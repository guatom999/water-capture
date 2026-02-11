package repositories

import (
	"context"
	"time"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockWaterLevelRepository is a mock implementation of WaterLevelRepositoryInterface
type MockWaterLevelRepository struct {
	mock.Mock
}

func (m *MockWaterLevelRepository) GetAll(ctx context.Context, limit int) ([]models.LocationWithWaterLevel, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LocationWithWaterLevel), args.Error(1)
}

func (m *MockWaterLevelRepository) GetLocationByID(ctx context.Context, stationID int) (*entities.Location, error) {
	args := m.Called(ctx, stationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Location), args.Error(1)
}

func (m *MockWaterLevelRepository) GetWaterLevelByID(ctx context.Context, locationID int) ([]*entities.WaterLevel, error) {
	args := m.Called(ctx, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WaterLevel), args.Error(1)
}

func (m *MockWaterLevelRepository) CreateProvince(ctx context.Context, req []*entities.Province) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockWaterLevelRepository) CreateStationLocation(ctx context.Context, req []*entities.Location) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockWaterLevelRepository) CreateWaterLevel(ctx context.Context, req []*entities.WaterLevel) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockWaterLevelRepository) MarkForDeletion(ctx context.Context, id int64, scheduledAt time.Time) error {
	args := m.Called(ctx, id, scheduledAt)
	return args.Error(0)
}

func (m *MockWaterLevelRepository) GetPendingDeletions(ctx context.Context) ([]*entities.WaterLevel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WaterLevel), args.Error(1)
}

func (m *MockWaterLevelRepository) HardDelete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
