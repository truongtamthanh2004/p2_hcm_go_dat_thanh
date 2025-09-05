package usecase

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"payment-service/internal/config"
	"payment-service/internal/model"
	"payment-service/internal/repository"
	"payment-service/utils"
	"strconv"
	"strings"
	"time"
)

type PaymentUsecase interface {
	CreatePaymentUrl(bookingID uint, clientIP string) (string, error)
	HandleVnpReturn(params url.Values) (string, error)
}

type paymentUsecaseImpl struct {
	txRepo            repository.TransactionRepository
	cfg               config.VnpayConfig
	httpClient        *http.Client
	bookingServiceURL string // e.g. http://booking-service:8080
}

func NewPaymentUsecase(repo repository.TransactionRepository, cfg config.VnpayConfig, bookingURL string) PaymentUsecase {
	return &paymentUsecaseImpl{
		txRepo:            repo,
		cfg:               cfg,
		httpClient:        &http.Client{Timeout: 5 * time.Second},
		bookingServiceURL: bookingURL,
	}
}

func (s *paymentUsecaseImpl) CreatePaymentUrl(bookingID uint, clientIP string) (string, error) {
	// Call booking-service to get booking
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/api/v1/bookings/%d", s.bookingServiceURL, bookingID))
	if err != nil {
		return "", fmt.Errorf("failed to call booking-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("booking not found")
	}

	var booking struct {
		ID         uint    `json:"id"`
		Status     string  `json:"status"`
		TotalPrice float64 `json:"TotalPrice"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		return "", err
	}

	if booking.Status != "PENDING" {
		return "", fmt.Errorf("booking already processed")
	}

	// Create transaction
	txnRef := strconv.FormatInt(time.Now().UnixNano(), 10)
	tx := &model.PaymentTransaction{TxnRef: txnRef, BookingID: bookingID, Status: "PENDING"}
	if err := s.txRepo.Create(tx); err != nil {
		return "", err
	}

	// Build VNPAY URL
	amount := int(math.Round(booking.TotalPrice * 100))
	params := url.Values{}
	params.Add("vnp_Version", "2.1.0")
	params.Add("vnp_Command", "pay")
	params.Add("vnp_TmnCode", s.cfg.TmnCode)
	params.Add("vnp_Amount", strconv.Itoa(amount))
	params.Add("vnp_CurrCode", "VND")
	params.Add("vnp_TxnRef", txnRef)
	params.Add("vnp_OrderInfo", fmt.Sprintf("Booking #%d", booking.ID))
	params.Add("vnp_OrderType", "other")
	params.Add("vnp_Locale", "vn")
	params.Add("vnp_ReturnUrl", s.cfg.ReturnURL)
	params.Add("vnp_IpAddr", clientIP)
	params.Add("vnp_CreateDate", time.Now().Format("20060102150405"))
	params.Add("vnp_ExpireDate", time.Now().Add(15*time.Minute).Format("20060102150405"))

	signedUrl := utils.BuildVnpUrl(params, s.cfg.HashSecret, s.cfg.PayURL, s.cfg.HashType)
	return signedUrl, nil
}

func (s *paymentUsecaseImpl) HandleVnpReturn(params url.Values) (string, error) {
	txnRef := params.Get("vnp_TxnRef")
	tx, err := s.txRepo.FindByTxnRef(txnRef)
	if err != nil {
		return "", err
	}

	if params.Get("vnp_ResponseCode") != "00" {
		tx.Status = "FAILED"
		_ = s.txRepo.Update(tx)
		return "", fmt.Errorf("payment failed")
	}

	tx.Status = "SUCCESS"
	_ = s.txRepo.Update(tx)

	// Call booking-service to update booking status
	bookingURL := fmt.Sprintf("%s/api/v1/bookings/%d/status", s.bookingServiceURL, tx.BookingID)
	body := strings.NewReader(`{"status":"CONFIRMED"}`)
	req, _ := http.NewRequest(http.MethodPut, bookingURL, body)
	req.Header.Set("Content-Type", "application/json")

	_, err = s.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	return s.cfg.ReturnURL, nil
}
