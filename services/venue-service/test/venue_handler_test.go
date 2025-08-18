package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"venue-service/internal/handler"
	"venue-service/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockVenueUsecase struct {
	mock.Mock
}

func (m *MockVenueUsecase) CreateVenue(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}
func (m *MockVenueUsecase) GetVenue(id uint) (*model.Venue, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Venue), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockVenueUsecase) UpdateVenue(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}
func (m *MockVenueUsecase) DeleteVenue(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockVenueUsecase) SearchVenues(city, name string) ([]model.Venue, error) {
	args := m.Called(city, name)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Venue), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupRouter(mockUC *MockVenueUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handler.NewVenueHandler(mockUC)
	r.POST("/venues", h.CreateVenue)
	r.GET("/venues/:id", h.GetVenue)
	r.PUT("/venues/:id", h.UpdateVenue)
	r.DELETE("/venues/:id", h.DeleteVenue)
	r.GET("/venues", h.SearchVenues)
	return r
}

func TestCreateVenue_Success(t *testing.T) {
	mockUC := new(MockVenueUsecase)
	mockUC.On("CreateVenue", mock.Anything).Return(nil)

	router := setupRouter(mockUC)

	body := `{"name":"Stadium","address":"123 Street","city":"HCM","description":"Nice"}`
	req, _ := http.NewRequest("POST", "/venues", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "venue.created", resp["message"])
}

func TestGetVenue_NotFound_WithInvalidID(t *testing.T) {
	mockUC := new(MockVenueUsecase)
	mockUC.On("GetVenue", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	router := setupRouter(mockUC)
	req, _ := http.NewRequest("GET", "/venues/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateVenue_Success(t *testing.T) {
	mockUC := new(MockVenueUsecase)
	existing := &model.Venue{Name: "Old", Address: "A"}
	existing.ID = uint(1)
	mockUC.On("GetVenue", uint(1)).Return(existing, nil)
	mockUC.On("UpdateVenue", mock.Anything).Return(nil)

	router := setupRouter(mockUC)
	body := `{"name":"New","address":"B","city":"HCM","status":"open","description":"desc"}`
	req, _ := http.NewRequest("PUT", "/venues/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteVenue_Success(t *testing.T) {
	mockUC := new(MockVenueUsecase)
	mockUC.On("DeleteVenue", uint(1)).Return(nil)

	router := setupRouter(mockUC)
	req, _ := http.NewRequest("DELETE", "/venues/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestSearchVenues_Success(t *testing.T) {
	mockUC := new(MockVenueUsecase)
	mockUC.On("SearchVenues", "HCM", "").Return([]model.Venue{{Name: "Stadium"}}, nil)

	router := setupRouter(mockUC)
	req, _ := http.NewRequest("GET", "/venues?city=HCM", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "venues.found", resp["message"])
}
