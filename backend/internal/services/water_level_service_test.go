package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper to create a service with a mock repo
func setupService(mockRepo *repositories.MockWaterLevelRepository) WaterLevelServiceInterface {
	cfg := &config.Config{
		App: config.App{
			BaseURL:   "http://localhost:8080",
			UploadDir: "/tmp/uploads",
		},
	}
	return NewWaterLevelService(mockRepo, cfg.App.BaseURL, cfg)
}

// ============================================================
// GetAllLocations
// ============================================================

func TestGetAllLocations_Success(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	now := time.Now()
	levelCm := 1.5
	danger := "SAFE"
	isFlooded := false
	note := "test note"
	waterLevelID := int64(1)

	mockData := []models.LocationWithWaterLevel{
		{
			StationID:           100,
			LocationName:        "Station A",
			LocationDescription: "Description A",
			ProvinceID:          13,
			Latitude:            13.75,
			Longitude:           100.50,
			IsActive:            true,
			BankLevel:           3.0,
			WaterLevelID:        &waterLevelID,
			LevelCm:             &levelCm,
			Image:               nil,
			Danger:              &danger,
			IsFlooded:           &isFlooded,
			MeasuredAt:          &now,
			Note:                &note,
		},
	}

	mockRepo.On("GetAll", mock.Anything, 10).Return(mockData, nil)

	result, err := svc.GetAllLocations(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(100), result[0].StationID)
	assert.Equal(t, "Station A", result[0].LocationName)
	assert.Equal(t, 13, result[0].ProvinceID)
	assert.Equal(t, &levelCm, result[0].LevelCm)
	assert.Nil(t, result[0].Image) // nil image should stay nil
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocations_WithImage(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	image := "water_photo.png"
	mockData := []models.LocationWithWaterLevel{
		{
			StationID:    200,
			LocationName: "Station B",
			IsActive:     true,
			Image:        &image,
		},
	}

	mockRepo.On("GetAll", mock.Anything, 5).Return(mockData, nil)

	result, err := svc.GetAllLocations(context.Background(), 5)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NotNil(t, result[0].Image)
	assert.Contains(t, *result[0].Image, "http://localhost:8080")
	assert.Contains(t, *result[0].Image, "water_photo.png")
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocations_EmptyResult(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("GetAll", mock.Anything, 10).Return([]models.LocationWithWaterLevel{}, nil)

	result, err := svc.GetAllLocations(context.Background(), 10)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocations_RepoError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("GetAll", mock.Anything, 10).Return(nil, errors.New("database error"))

	result, err := svc.GetAllLocations(context.Background(), 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============================================================
// GetByLocationID
// ============================================================

func TestGetByLocationID_Success(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	now := time.Now()
	location := &entities.Location{
		ID:        1,
		StationID: 100,
		Name:      "Station A",
		BankLevel: 3.0,
	}

	waterLevels := []*entities.WaterLevel{
		{
			ID:         1,
			StationID:  100,
			LevelCm:    1.5,
			Image:      "img.png",
			Danger:     "SAFE",
			IsFlooded:  false,
			Source:     sql.NullString{String: "CCTV", Valid: true},
			MeasuredAt: now,
			Note:       "normal",
		},
		{
			ID:         2,
			StationID:  100,
			LevelCm:    2.0,
			Image:      "img2.png",
			Danger:     "WATCH",
			IsFlooded:  false,
			Source:     sql.NullString{String: "CCTV", Valid: true},
			MeasuredAt: now.Add(-time.Hour),
			Note:       "rising",
		},
	}

	mockRepo.On("GetLocationByID", mock.Anything, 100).Return(location, nil)
	mockRepo.On("GetWaterLevelByID", mock.Anything, 100).Return(waterLevels, nil)

	result, err := svc.GetByLocationID(context.Background(), "100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.StationID)
	assert.Equal(t, float64(3.0), result.BankLevel)
	assert.Len(t, result.Detail, 2)
	assert.Equal(t, 1.5, result.Detail[0].LevelCm)
	assert.Equal(t, "SAFE", result.Detail[0].Danger)
	assert.Equal(t, 2.0, result.Detail[1].LevelCm)
	mockRepo.AssertExpectations(t)
}

func TestGetByLocationID_InvalidID(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	result, err := svc.GetByLocationID(context.Background(), "not-a-number")

	assert.Error(t, err)
	assert.Nil(t, result)
	// repo should never be called with invalid ID
	mockRepo.AssertNotCalled(t, "GetLocationByID")
	mockRepo.AssertNotCalled(t, "GetWaterLevelByID")
}

func TestGetByLocationID_LocationRepoError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("GetLocationByID", mock.Anything, 999).Return(nil, errors.New("not found"))

	result, err := svc.GetByLocationID(context.Background(), "999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetByLocationID_WaterLevelRepoError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	location := &entities.Location{
		ID:        1,
		StationID: 100,
		BankLevel: 3.0,
	}

	mockRepo.On("GetLocationByID", mock.Anything, 100).Return(location, nil)
	mockRepo.On("GetWaterLevelByID", mock.Anything, 100).Return(nil, errors.New("query failed"))

	result, err := svc.GetByLocationID(context.Background(), "100")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "query failed", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============================================================
// CreateWaterLevel
// ============================================================

func TestCreateWaterLevel_Success(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	req := &models.CreateWaterLevelReq{
		LocationID: 1,
		LevelCm:    1.5,
		Image:      "test.png",
		Danger:     "SAFE",
		IsFlooded:  false,
		MeasuredAt: time.Now(),
		Note:       "test",
	}

	// CreateWaterLevel is currently a no-op (returns nil without calling repo)
	err := svc.CreateWaterLevel(context.Background(), req)

	assert.NoError(t, err)
}

// ============================================================
// GetAllLocations — edge cases
// ============================================================

func TestGetAllLocations_EmptyImageString(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	emptyImage := ""
	mockData := []models.LocationWithWaterLevel{
		{
			StationID:    300,
			LocationName: "Station C",
			IsActive:     true,
			Image:        &emptyImage,
		},
	}

	mockRepo.On("GetAll", mock.Anything, 10).Return(mockData, nil)

	result, err := svc.GetAllLocations(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Nil(t, result[0].Image) // empty string image should become nil
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocations_NilOptionalFields(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockData := []models.LocationWithWaterLevel{
		{
			StationID:    400,
			LocationName: "Station D",
			IsActive:     true,
			WaterLevelID: nil,
			LevelCm:      nil,
			Image:        nil,
			Danger:       nil,
			IsFlooded:    nil,
			MeasuredAt:   nil,
			Note:         nil,
		},
	}

	mockRepo.On("GetAll", mock.Anything, 10).Return(mockData, nil)

	result, err := svc.GetAllLocations(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Nil(t, result[0].WaterLevelID)
	assert.Nil(t, result[0].LevelCm)
	assert.Nil(t, result[0].Image)
	assert.Nil(t, result[0].Danger)
	assert.Nil(t, result[0].IsFlooded)
	assert.Nil(t, result[0].Note)
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocations_MultipleResults(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	now := time.Now()
	level1 := 1.5
	level2 := 2.0
	danger1 := "SAFE"
	danger2 := "WATCH"
	flooded1 := false
	flooded2 := false
	wlID1 := int64(1)
	wlID2 := int64(2)
	img := "photo.png"

	mockData := []models.LocationWithWaterLevel{
		{
			StationID:    100,
			LocationName: "Station A",
			IsActive:     true,
			WaterLevelID: &wlID1,
			LevelCm:      &level1,
			Danger:       &danger1,
			IsFlooded:    &flooded1,
			MeasuredAt:   &now,
			Image:        &img,
		},
		{
			StationID:    200,
			LocationName: "Station B",
			IsActive:     true,
			WaterLevelID: &wlID2,
			LevelCm:      &level2,
			Danger:       &danger2,
			IsFlooded:    &flooded2,
			MeasuredAt:   &now,
			Image:        nil,
		},
	}

	mockRepo.On("GetAll", mock.Anything, 20).Return(mockData, nil)

	result, err := svc.GetAllLocations(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NotNil(t, result[0].Image)
	assert.Contains(t, *result[0].Image, "http://localhost:8080")
	assert.Nil(t, result[1].Image)
	mockRepo.AssertExpectations(t)
}

// ============================================================
// GetByLocationID — edge cases
// ============================================================

func TestGetByLocationID_WithImageURL(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	now := time.Now()
	location := &entities.Location{
		ID:        1,
		StationID: 100,
		Name:      "Station A",
		BankLevel: 3.0,
	}

	waterLevels := []*entities.WaterLevel{
		{
			ID:         1,
			StationID:  100,
			LevelCm:    1.5,
			Image:      "img.png",
			Danger:     "SAFE",
			IsFlooded:  false,
			Source:     sql.NullString{String: "", Valid: false},
			MeasuredAt: now,
			Note:       "test",
		},
	}

	mockRepo.On("GetLocationByID", mock.Anything, 100).Return(location, nil)
	mockRepo.On("GetWaterLevelByID", mock.Anything, 100).Return(waterLevels, nil)

	result, err := svc.GetByLocationID(context.Background(), "100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.StationID)
	assert.Equal(t, float64(3.0), result.BankLevel)
	assert.Len(t, result.Detail, 1)
	assert.Equal(t, "img.png", result.Detail[0].Image)
	assert.Equal(t, false, result.Detail[0].Source.Valid)
	mockRepo.AssertExpectations(t)
}

func TestGetByLocationID_MultipleWaterLevels(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	now := time.Now()
	location := &entities.Location{
		ID:        1,
		StationID: 100,
		Name:      "Station A",
		BankLevel: 3.0,
	}

	waterLevels := []*entities.WaterLevel{
		{ID: 1, StationID: 100, LevelCm: 1.0, Danger: "SAFE", IsFlooded: false, Source: sql.NullString{Valid: false}, MeasuredAt: now, Note: "a"},
		{ID: 2, StationID: 100, LevelCm: 2.0, Danger: "WATCH", IsFlooded: false, Source: sql.NullString{String: "CCTV", Valid: true}, MeasuredAt: now.Add(-time.Hour), Note: "b"},
		{ID: 3, StationID: 100, LevelCm: 4.0, Danger: "DANGER", IsFlooded: true, Source: sql.NullString{String: "CCTV", Valid: true}, MeasuredAt: now.Add(-2 * time.Hour), Note: "c"},
	}

	mockRepo.On("GetLocationByID", mock.Anything, 100).Return(location, nil)
	mockRepo.On("GetWaterLevelByID", mock.Anything, 100).Return(waterLevels, nil)

	result, err := svc.GetByLocationID(context.Background(), "100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Detail, 3)
	assert.Equal(t, "SAFE", result.Detail[0].Danger)
	assert.Equal(t, "WATCH", result.Detail[1].Danger)
	assert.Equal(t, "DANGER", result.Detail[2].Danger)
	assert.True(t, result.Detail[2].IsFlooded)
	mockRepo.AssertExpectations(t)
}

func TestGetByLocationID_ZeroID(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	location := &entities.Location{ID: 1, StationID: 0, BankLevel: 2.0}
	waterLevels := []*entities.WaterLevel{
		{ID: 1, StationID: 0, LevelCm: 1.0, Danger: "SAFE", MeasuredAt: time.Now()},
	}

	mockRepo.On("GetLocationByID", mock.Anything, 0).Return(location, nil)
	mockRepo.On("GetWaterLevelByID", mock.Anything, 0).Return(waterLevels, nil)

	result, err := svc.GetByLocationID(context.Background(), "0")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.StationID)
	mockRepo.AssertExpectations(t)
}

// ============================================================
// ScheduleDeleteWaterLevel
// ============================================================

func TestScheduleDeleteWaterLevel_Success(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("MarkForDeletion", mock.Anything, int64(28), mock.AnythingOfType("time.Time")).Return(nil)
	mockRepo.On("GetPendingDeletions", mock.Anything).Return([]*entities.WaterLevel{}, nil)

	err := svc.ScheduleDeleteWaterLevel(context.Background(), "test.png", 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestScheduleDeleteWaterLevel_MarkForDeletionError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("MarkForDeletion", mock.Anything, int64(28), mock.AnythingOfType("time.Time")).Return(errors.New("mark failed"))

	err := svc.ScheduleDeleteWaterLevel(context.Background(), "test.png", 1)

	assert.Error(t, err)
	assert.Equal(t, "mark failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestScheduleDeleteWaterLevel_GetPendingDeletionsError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	mockRepo.On("MarkForDeletion", mock.Anything, int64(28), mock.AnythingOfType("time.Time")).Return(nil)
	mockRepo.On("GetPendingDeletions", mock.Anything).Return(nil, errors.New("query failed"))

	err := svc.ScheduleDeleteWaterLevel(context.Background(), "test.png", 1)

	assert.Error(t, err)
	assert.Equal(t, "query failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestScheduleDeleteWaterLevel_WithPendingItems(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	pendingItems := []*entities.WaterLevel{
		{ID: 1, Image: "old1.png", Status: "PENDING_DELETION"},
		{ID: 2, Image: "old2.png", Status: "PENDING_DELETION"},
	}

	mockRepo.On("MarkForDeletion", mock.Anything, int64(28), mock.AnythingOfType("time.Time")).Return(nil)
	mockRepo.On("GetPendingDeletions", mock.Anything).Return(pendingItems, nil)
	// HardDelete will be called for each pending item
	mockRepo.On("HardDelete", mock.Anything, int64(1)).Return(nil)
	mockRepo.On("HardDelete", mock.Anything, int64(2)).Return(nil)

	err := svc.ScheduleDeleteWaterLevel(context.Background(), "test.png", 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestScheduleDeleteWaterLevel_HardDeleteError(t *testing.T) {
	mockRepo := new(repositories.MockWaterLevelRepository)
	svc := setupService(mockRepo)

	pendingItems := []*entities.WaterLevel{
		{ID: 1, Image: "old1.png", Status: "PENDING_DELETION"},
	}

	mockRepo.On("MarkForDeletion", mock.Anything, int64(28), mock.AnythingOfType("time.Time")).Return(nil)
	mockRepo.On("GetPendingDeletions", mock.Anything).Return(pendingItems, nil)
	// HardDelete fails but ScheduleDeleteWaterLevel continues (logs error, doesn't return it)
	mockRepo.On("HardDelete", mock.Anything, int64(1)).Return(errors.New("delete failed"))

	err := svc.ScheduleDeleteWaterLevel(context.Background(), "test.png", 1)

	// ScheduleDeleteWaterLevel doesn't propagate HardDelete errors
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
