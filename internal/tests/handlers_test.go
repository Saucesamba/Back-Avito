package handlers

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/handlers"
	"Backend_trainee_assigment_2025/internal/schemas"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockAvitoDB - это mock для AvitoDB
type MockAvitoDB struct {
	NewDBFn         func(cfg *config.DBConfig) (db.Database, error)
	GetUserFn       func(ctx context.Context, user schemas.UserLogin) (*schemas.User, error)
	CreateUserFn    func(ctx context.Context, user *schemas.UserReg) (*schemas.User, error)
	OpenPVZFn       func(ctx context.Context, city string) (*schemas.PVZ, error)
	GetPVZFn        func(ctx context.Context, startTime, endTime string, offset, limit int) ([]schemas.PVZWithReceptionsAndProducts, error)
	OpenRecFn       func(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error)
	CloseLastRecFn  func(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error)
	CreateProductFn func(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error)
	DeleteProductFn func(ctx context.Context, pvzId uuid.UUID) error
	GetProductFn    func(id uuid.UUID) ([]schemas.Product, error)
}

// Implement all the methods from DB.Database interface
func (m *MockAvitoDB) NewDB(cfg *config.DBConfig) (db.Database, error) {
	if m.NewDBFn != nil {
		return m.NewDBFn(cfg)
	}
	return nil, nil
}

func (m *MockAvitoDB) GetUser(ctx context.Context, user schemas.UserLogin) (*schemas.User, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, user)
	}
	return nil, nil
}

func (m *MockAvitoDB) CreateUser(ctx context.Context, user *schemas.UserReg) (*schemas.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return nil, nil
}

func (m *MockAvitoDB) OpenPVZ(ctx context.Context, city string) (*schemas.PVZ, error) {
	if m.OpenPVZFn != nil {
		return m.OpenPVZFn(ctx, city)
	}
	return nil, nil
}

func (m *MockAvitoDB) GetPVZ(ctx context.Context, startTime, endTime string, offset, limit int) ([]schemas.PVZWithReceptionsAndProducts, error) {
	if m.GetPVZFn != nil {
		return m.GetPVZFn(ctx, startTime, endTime, offset, limit)
	}
	return nil, nil
}

func (m *MockAvitoDB) OpenRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
	if m.OpenRecFn != nil {
		return m.OpenRecFn(ctx, pvzId)
	}
	return nil, nil
}

func (m *MockAvitoDB) CloseLastRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
	if m.CloseLastRecFn != nil {
		return m.CloseLastRecFn(ctx, pvzId)
	}
	return nil, nil
}
func (m *MockAvitoDB) AddProduct(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error) {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, typ, pvzId)
	}

	return nil, errors.New("Not Implemented Method") // Implement
}

func (m *MockAvitoDB) CreateProduct(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error) {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, typ, pvzId)
	}

	return nil, errors.New("Not Implemented Method") // Implement
}

func (m *MockAvitoDB) DeleteProduct(ctx context.Context, pvzId uuid.UUID) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, pvzId)
	}

	return errors.New("Not Implemented Method") // Implement
}

func (m *MockAvitoDB) GetProduct(id uuid.UUID) ([]schemas.Product, error) {
	return nil, nil
}

func TestAddProductHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	//Setup expected product for createProduct
	expectedProduct := &schemas.Product{
		Type: "electronics",
	}
	mockDB.CreateProductFn = func(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error) {
		if typ == expectedProduct.Type {
			expectedProduct.ReceptionId = pvzId

			return expectedProduct, nil
		}

		return nil, errors.New("product error")
	}
	handler := handlers.NewProductHandler(mockDB, cfg)

	pvzIdTest := uuid.New()
	body := map[string]string{
		"type":  expectedProduct.Type,
		"pvzId": pvzIdTest.String(),
	}

	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/products", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler.AddProductHandler(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code, "It's supposed to return  Status Created")

	var response schemas.Product
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, expectedProduct.Type, response.Type)

}

func TestDeleteLastProductHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	mockDB.DeleteProductFn = func(ctx context.Context, pvzId uuid.UUID) error {

		//Add more conditions you have
		return nil
	}

	handler := handlers.NewProductHandler(mockDB, cfg) // Initialize handler

	pvzIdTest := uuid.New() //Create a test UUID

	req, err := http.NewRequest("POST", "/pvz/{pvzId}/delete_last_product", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a ResponseRecorder to record the response.
	recorder := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"pvzId": pvzIdTest.String()}) //Added to MUX Urls
	// Act
	handler.DeleteLastProductHandler(recorder, req)

	// Assert - I use testing library to confirm
	assert.Equal(t, http.StatusOK, recorder.Code, "it should return status OK")
}

func TestCreatePVZHandler(t *testing.T) {
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)
	expectedPVZ := &schemas.PVZ{
		City: "Москва",
	}

	// Setup mock to return successful value
	mockDB.OpenPVZFn = func(ctx context.Context, city string) (*schemas.PVZ, error) {
		if city == expectedPVZ.City {
			expectedPVZ.Id = uuid.New()

			return expectedPVZ, nil
		}
		return nil, errors.New("pvz error")
	}

	handler := handlers.NewPVZHandler(mockDB, cfg)
	body := map[string]string{
		"city": expectedPVZ.City,
	}

	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/pvz", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler.CreatePVZHandler(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code, "It's supposed to return  Status Created")

	var response schemas.PVZ
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, expectedPVZ.City, response.City)
}

func TestGetPVZsHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)
	pvzIdTest := uuid.New()

	// Setup the mock database to return specific data
	expectedPVZs := []schemas.PVZWithReceptionsAndProducts{{
		PVZ: schemas.PVZ{
			Id:               pvzIdTest,
			RegistrationDate: time.Now(),
			City:             "Москва",
		},
	}}
	mockDB.GetPVZFn = func(ctx context.Context, startTime, endTime string, offset, limit int) ([]schemas.PVZWithReceptionsAndProducts, error) {
		return expectedPVZs, nil
	}

	// Create Request
	handler := handlers.NewPVZHandler(mockDB, cfg)
	reqBody := map[string]string{
		"start_time": "2024-04-20T18:22:03.710748Z",
		"end_time":   "2024-04-20T22:22:03.235135Z",
		"page":       "1",
		"limit":      "10",
	}

	//Format all the information
	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("GET", "/pvzs", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// We create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler.GetPVZsHandler(rr, req)

	// Assert
	var actual []schemas.PVZWithReceptionsAndProducts
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("Error unmarshaling reponse: %v", err)
	}
	// Basic check to see if the values line up, for now
	assert.Equal(t, expectedPVZs[0].PVZ.City, actual[0].PVZ.City, "Check")
}

func TestCreateReceptionHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)
	pvzIdTest := uuid.New()

	expectedReception := &schemas.Reception{
		PVZId: pvzIdTest,
	}

	mockDB.OpenRecFn = func(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
		if pvzId == expectedReception.PVZId {
			expectedReception.Id = uuid.New() // Assign new ID
			expectedReception.DateTime = time.Now()
			expectedReception.Status = "in_progress" // Set status

			return expectedReception, nil
		}

		return nil, errors.New("reception error")
	}

	handler := handlers.NewReceptionHandler(mockDB, cfg)

	body := map[string]string{
		"pvzId": expectedReception.PVZId.String(),
	}

	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/receptions", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler.CreateReceptionHandler(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code, "It's supposed to return Status Created")

	var response schemas.Reception
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, expectedReception.PVZId.String(), response.PVZId.String(), "Check if response has the PVZ set")

}

func TestCloseLastReceptionHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	//Setup
	pvzIdTest := uuid.New()

	// Setup the mock database to return specific data
	expectedReceptions := schemas.Reception{
		PVZId:  pvzIdTest,
		Status: "close",
	}
	mockDB.CloseLastRecFn = func(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
		return &expectedReceptions, nil
	}
	handlers.NewProductHandler(mockDB, cfg)
	// Create Request

	handler := handlers.NewReceptionHandler(mockDB, cfg)
	req, err := http.NewRequest("POST", "/pvz/{pvzId}/close_last_reception", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Setup url vars since this will happen during runtime with the mux
	vars := map[string]string{
		"pvzId": pvzIdTest.String(),
	}
	req = mux.SetURLVars(req, vars)

	// We create a ResponseRecorder to record the response.
	recorder := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler.CloseLastReceptionHandler)
	// Act

	testHandler.ServeHTTP(recorder, req)
	// Assert

	var actual schemas.Reception
	err = json.Unmarshal(recorder.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("Error unmarshaling reponse: %v", err)
	}

	assert.Equal(t, http.StatusOK, recorder.Code, "it should return status OK")

	assert.Equal(t, expectedReceptions.Status, actual.Status)
}

func TestDummyLoginHandler(t *testing.T) {

	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	handler := handlers.NewUserHandler(mockDB, cfg)

	body := map[string]string{
		"role": "moderator",
	}

	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/dummyLogin", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler.DummyLoginHandler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "It's supposed to return  Status Created")

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, response["token"], response["token"])

}

func TestRegisterHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	expectedUser := &schemas.User{
		Email: "test@example.com",
		Role:  "employee",
	}

	mockDB.CreateUserFn = func(ctx context.Context, user *schemas.UserReg) (*schemas.User, error) {
		if user.Email == expectedUser.Email && user.Role == expectedUser.Role {
			//This returns a object
			return &schemas.User{
				Email: user.Email,
				Role:  user.Role,
			}, nil // This user contains a value

		}
		return nil, errors.New("user error")
	}

	handler := handlers.NewUserHandler(mockDB, cfg)

	body := map[string]string{
		"email":    expectedUser.Email,
		"password": "password",
		"role":     expectedUser.Role,
	}

	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler.RegisterHandler(recorder, req)

	//Assert
	assert.Equal(t, http.StatusCreated, recorder.Code, "It's supposed to return  Status Created")

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
		return
	}
}

func TestLoginHandler(t *testing.T) {
	// Arrange
	mockDB := new(MockAvitoDB)
	cfg := new(config.AppConfig)

	emailTest := "test@example.com"

	mockDB.GetUserFn = func(ctx context.Context, user schemas.UserLogin) (*schemas.User, error) {
		if user.Email == emailTest {
			//Setup fake ID that it will always be equals
			return &schemas.User{
				Email: user.Email,
			}, nil

		}
		return nil, errors.New("not found")
	}

	handler := handlers.NewUserHandler(mockDB, cfg) //Pass App.Config

	// Create a request body
	reqBody := map[string]string{
		"email":    emailTest,
		"password": "test",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Could not marshal request body: %v", err)
	}
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	//Create a response recorder
	recorder := httptest.NewRecorder()
	handler.LoginHandler(recorder, req)

}
