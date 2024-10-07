package tests_test

import (
	"context"
	"testing"

	"L0/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Save(ctx context.Context, orderUID string, orderInfo []byte) error {
	args := m.Called(ctx, orderUID, orderInfo)
	return args.Error(0)
}

func (m *MockDB) Get(ctx context.Context, orderUID string) ([]byte, error) {
	args := m.Called(ctx, orderUID)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDB) GetAll(ctx context.Context) ([]model.OrderInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.OrderInfo), args.Error(1)
}

func TestStorage_Save(t *testing.T) {
	mockDB := new(MockDB)
	orderUID := "12345"
	orderInfo := []byte(`{"order_uid": "12345", "track_number": "TRACK123"}`)

	mockDB.On("Save", mock.Anything, orderUID, orderInfo).Return(nil)

	err := mockDB.Save(context.Background(), orderUID, orderInfo)
	assert.NoError(t, err)

	mockDB.AssertCalled(t, "Save", mock.Anything, orderUID, orderInfo)
}
