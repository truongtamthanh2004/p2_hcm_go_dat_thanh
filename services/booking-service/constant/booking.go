package constant

const (
	BookingStatusPending   = "PENDING"
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusCancelled = "CANCELLED"
)

var AllowedBookingStatuses = map[string]bool{
	BookingStatusPending:   true,
	BookingStatusConfirmed: true,
	BookingStatusCancelled: true,
}
