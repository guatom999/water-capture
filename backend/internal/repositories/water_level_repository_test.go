package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// helper to create a sqlmock-backed repository
func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, WaterLevelRepositoryInterface) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWaterLevelRepository(sqlxDB)
	return sqlxDB, mock, repo
}

// ============================================================
// GetAll
// ============================================================

func TestGetAll_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"station_id", "location_name", "location_description", "province_id",
		"latitude", "longitude", "is_active", "bank_level",
		"water_level_id", "level_cm", "image", "danger", "is_flooded", "measured_at", "note",
	}).AddRow(100, "Station A", "Desc A", 13, 13.75, 100.50, true, 3.0, 1, 1.5, nil, "SAFE", false, now, "ok")

	mock.ExpectQuery("SELECT DISTINCT ON").WillReturnRows(rows)

	result, err := repo.GetAll(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(100), result[0].StationID)
	assert.Equal(t, "Station A", result[0].LocationName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_QueryError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT DISTINCT ON").WillReturnError(sql.ErrConnDone)

	result, err := repo.GetAll(context.Background(), 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetLocationByID
// ============================================================

func TestGetLocationByID_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "station_id", "name", "description", "province_id",
		"bank_level", "latitude", "longitude", "is_active", "created_at", "updated_at",
	}).AddRow(1, 100, "Station A", "Desc A", 13, 3.0, 13.75, 100.50, true, now, now)

	mock.ExpectQuery("SELECT \\* FROM locations WHERE station_id").
		WithArgs(100).
		WillReturnRows(rows)

	result, err := repo.GetLocationByID(context.Background(), 100)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.StationID)
	assert.Equal(t, "Station A", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocationByID_NotFound(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT \\* FROM locations WHERE station_id").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetLocationByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetWaterLevelByID
// ============================================================

func TestGetWaterLevelByID_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "station_id", "level_cm", "image", "danger", "is_flooded",
		"source", "measured_at", "note", "status", "deleted_at", "scheduled_delete_at",
	}).
		AddRow(1, 100, 1.5, "img.png", "SAFE", false, "CCTV", now, "ok", "ACTIVE", nil, nil).
		AddRow(2, 100, 2.0, "img2.png", "WATCH", false, "CCTV", now, "rising", "ACTIVE", nil, nil)

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE station_id").
		WithArgs(100).
		WillReturnRows(rows)

	result, err := repo.GetWaterLevelByID(context.Background(), 100)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, 1.5, result[0].LevelCm)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWaterLevelByID_Empty(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	rows := sqlmock.NewRows([]string{
		"id", "station_id", "level_cm", "image", "danger", "is_flooded",
		"source", "measured_at", "note", "status", "deleted_at", "scheduled_delete_at",
	})

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE station_id").
		WithArgs(999).
		WillReturnRows(rows)

	result, err := repo.GetWaterLevelByID(context.Background(), 999)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateProvince
// ============================================================

func TestCreateProvince_EmptyInput(t *testing.T) {
	sqlxDB, _, repo := setupMockDB(t)
	defer sqlxDB.Close()

	err := repo.CreateProvince(context.Background(), nil)
	assert.NoError(t, err)
}

func TestCreateProvince_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	provinces := []*entities.Province{
		{Name: "Pathum Thani", Code: "13"},
		{Name: "Bangkok", Code: "10"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO provinces").
		WithArgs("Pathum Thani", "13").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO provinces").
		WithArgs("Bangkok", "10").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err := repo.CreateProvince(context.Background(), provinces)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateStationLocation
// ============================================================

func TestCreateStationLocation_EmptyInput(t *testing.T) {
	sqlxDB, _, repo := setupMockDB(t)
	defer sqlxDB.Close()

	err := repo.CreateStationLocation(context.Background(), nil)
	assert.NoError(t, err)
}

func TestCreateStationLocation_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	locations := []*entities.Location{
		{
			StationID:   100,
			Name:        "Station A",
			Description: "Desc A",
			Latitude:    13.75,
			Longitude:   100.50,
			IsActive:    true,
			BankLevel:   3.0,
			ProvinceID:  13,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO locations").
		WithArgs(int64(100), "Station A", "Desc A", 13.75, 100.50, true, 3.0, 13).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateStationLocation(context.Background(), locations)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateWaterLevel
// ============================================================

func TestCreateWaterLevel_EmptyInput(t *testing.T) {
	sqlxDB, _, repo := setupMockDB(t)
	defer sqlxDB.Close()

	err := repo.CreateWaterLevel(context.Background(), nil)
	assert.NoError(t, err)
}

func TestCreateWaterLevel_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	waterLevels := []*entities.WaterLevel{
		{
			StationID:  100,
			LevelCm:    1.5,
			Image:      "img.png",
			Danger:     "SAFE",
			IsFlooded:  false,
			Source:     sql.NullString{String: "CCTV", Valid: true},
			MeasuredAt: now,
			Note:       "test",
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO water_levels").
		WithArgs(int64(100), 1.5, "img.png", "SAFE", false, sql.NullString{String: "CCTV", Valid: true}, now, "test", "ACTIVE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateWaterLevel(context.Background(), waterLevels)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// MarkForDeletion
// ============================================================

func TestMarkForDeletion_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	scheduledAt := time.Now()

	mock.ExpectExec("UPDATE water_levels").
		WithArgs(scheduledAt, int64(28)).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := repo.MarkForDeletion(context.Background(), 28, scheduledAt)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetPendingDeletions
// ============================================================

func TestGetPendingDeletions_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "station_id", "level_cm", "image", "danger", "is_flooded",
		"source", "measured_at", "note", "status", "deleted_at", "scheduled_delete_at",
	}).AddRow(1, 100, 1.5, "img.png", "SAFE", false, "CCTV", now, "ok", "PENDING_DELETION", nil, now)

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE status").WillReturnRows(rows)

	result, err := repo.GetPendingDeletions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "PENDING_DELETION", result[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// HardDelete
// ============================================================

func TestHardDelete_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectExec("DELETE FROM water_levels WHERE id").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.HardDelete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHardDelete_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectExec("DELETE FROM water_levels WHERE id").
		WithArgs(int64(1)).
		WillReturnError(sql.ErrConnDone)

	err := repo.HardDelete(context.Background(), 1)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetWaterLevelByID — error path
// ============================================================

func TestGetWaterLevelByID_QueryError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE station_id").
		WithArgs(100).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.GetWaterLevelByID(context.Background(), 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateProvince — error paths
// ============================================================

func TestCreateProvince_BeginTxError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	provinces := []*entities.Province{{Name: "Test", Code: "99"}}
	err := repo.CreateProvince(context.Background(), provinces)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProvince_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	provinces := []*entities.Province{
		{Name: "Pathum Thani", Code: "13"},
		{Name: "Bangkok", Code: "10"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO provinces").
		WithArgs("Pathum Thani", "13").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO provinces").
		WithArgs("Bangkok", "10").
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := repo.CreateProvince(context.Background(), provinces)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProvince_CommitError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	provinces := []*entities.Province{{Name: "Test", Code: "99"}}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO provinces").
		WithArgs("Test", "99").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

	err := repo.CreateProvince(context.Background(), provinces)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateStationLocation — error paths
// ============================================================

func TestCreateStationLocation_BeginTxError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	locations := []*entities.Location{{StationID: 100, Name: "Test"}}
	err := repo.CreateStationLocation(context.Background(), locations)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStationLocation_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	locations := []*entities.Location{
		{StationID: 100, Name: "Station A", Description: "Desc A", Latitude: 13.75, Longitude: 100.50, IsActive: true, BankLevel: 3.0, ProvinceID: 13},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO locations").
		WithArgs(int64(100), "Station A", "Desc A", 13.75, 100.50, true, 3.0, 13).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := repo.CreateStationLocation(context.Background(), locations)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStationLocation_CommitError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	locations := []*entities.Location{
		{StationID: 100, Name: "Station A", Description: "Desc A", Latitude: 13.75, Longitude: 100.50, IsActive: true, BankLevel: 3.0, ProvinceID: 13},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO locations").
		WithArgs(int64(100), "Station A", "Desc A", 13.75, 100.50, true, 3.0, 13).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

	err := repo.CreateStationLocation(context.Background(), locations)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CreateWaterLevel — error paths
// ============================================================

func TestCreateWaterLevel_BeginTxError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	now := time.Now()
	waterLevels := []*entities.WaterLevel{
		{StationID: 100, LevelCm: 1.5, Image: "img.png", Danger: "SAFE", IsFlooded: false, Source: sql.NullString{String: "CCTV", Valid: true}, MeasuredAt: now, Note: "test"},
	}

	err := repo.CreateWaterLevel(context.Background(), waterLevels)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWaterLevel_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	waterLevels := []*entities.WaterLevel{
		{StationID: 100, LevelCm: 1.5, Image: "img.png", Danger: "SAFE", IsFlooded: false, Source: sql.NullString{String: "CCTV", Valid: true}, MeasuredAt: now, Note: "test"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO water_levels").
		WithArgs(int64(100), 1.5, "img.png", "SAFE", false, sql.NullString{String: "CCTV", Valid: true}, now, "test", "ACTIVE").
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := repo.CreateWaterLevel(context.Background(), waterLevels)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWaterLevel_CommitError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	waterLevels := []*entities.WaterLevel{
		{StationID: 100, LevelCm: 1.5, Image: "img.png", Danger: "SAFE", IsFlooded: false, Source: sql.NullString{String: "CCTV", Valid: true}, MeasuredAt: now, Note: "test"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO water_levels").
		WithArgs(int64(100), 1.5, "img.png", "SAFE", false, sql.NullString{String: "CCTV", Valid: true}, now, "test", "ACTIVE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

	err := repo.CreateWaterLevel(context.Background(), waterLevels)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// MarkForDeletion — error path
// ============================================================

func TestMarkForDeletion_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	scheduledAt := time.Now()

	mock.ExpectExec("UPDATE water_levels").
		WithArgs(scheduledAt, int64(28)).
		WillReturnError(sql.ErrConnDone)

	err := repo.MarkForDeletion(context.Background(), 28, scheduledAt)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetPendingDeletions — error path
// ============================================================

func TestGetPendingDeletions_QueryError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE status").
		WillReturnError(sql.ErrConnDone)

	result, err := repo.GetPendingDeletions(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// CancelDeletion
// ============================================================

func TestCancelDeletion_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	// CancelDeletion is on the concrete type, access via interface assertion
	concreteRepo := repo.(*waterLevelRepository)

	mock.ExpectExec("UPDATE water_levels SET status").
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := concreteRepo.CancelDeletion(context.Background(), 5)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCancelDeletion_ExecError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	concreteRepo := repo.(*waterLevelRepository)

	mock.ExpectExec("UPDATE water_levels SET status").
		WithArgs(int64(5)).
		WillReturnError(sql.ErrConnDone)

	err := concreteRepo.CancelDeletion(context.Background(), 5)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetAll — empty result
// ============================================================

func TestGetAll_EmptyResult(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	rows := sqlmock.NewRows([]string{
		"station_id", "location_name", "location_description", "province_id",
		"latitude", "longitude", "is_active", "bank_level",
		"water_level_id", "level_cm", "image", "danger", "is_flooded", "measured_at", "note",
	})

	mock.ExpectQuery("SELECT DISTINCT ON").WillReturnRows(rows)

	result, err := repo.GetAll(context.Background(), 10)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetLocationByID — connection error
// ============================================================

func TestGetLocationByID_QueryError(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT \\* FROM locations WHERE station_id").
		WithArgs(100).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.GetLocationByID(context.Background(), 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// GetLastWaterLevelWithLimit
// ============================================================

func TestGetLastWaterLevelWithLimit_Success(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	concreteRepo := repo.(*waterLevelRepository)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "station_id", "level_cm", "image", "danger", "is_flooded",
		"source", "measured_at", "note", "status", "deleted_at", "scheduled_delete_at",
	}).AddRow(1, 100, 1.5, "img.png", "SAFE", false, "CCTV", now, "ok", "ACTIVE", nil, nil)

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE location_id").
		WithArgs(100).
		WillReturnRows(rows)

	result, err := concreteRepo.GetLastWaterLevelWithLimit(context.Background(), 100, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, 1.5, result.LevelCm)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastWaterLevelWithLimit_NotFound(t *testing.T) {
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	concreteRepo := repo.(*waterLevelRepository)

	mock.ExpectQuery("SELECT \\* FROM water_levels WHERE location_id").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	result, err := concreteRepo.GetLastWaterLevelWithLimit(context.Background(), 999, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// Notification & Auth repository additional tests
// ============================================================

func TestAuth_CreateUser_Success(t *testing.T) {
	sqlxDB, mock, _ := setupMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(5, now, now))

	repo := NewAuthRepository(sqlxDB)
	u := &entities.User{Email: "a@a.com", Password: "pass", Name: "Name", Role: "user", IsActive: true}
	res, err := repo.CreateUser(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), res.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuth_GetUserByEmail_NotFound(t *testing.T) {
	sqlxDB, mock, _ := setupMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT \\* FROM users WHERE email").WithArgs("x@x.com").WillReturnError(sql.ErrNoRows)

	repo := NewAuthRepository(sqlxDB)
	user, err := repo.GetUserByEmail(context.Background(), "x@x.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuth_CreateRefreshToken_ExecError(t *testing.T) {
	sqlxDB, mock, _ := setupMockDB(t)
	defer sqlxDB.Close()

	token := &entities.RefreshToken{UserID: 1, Token: "t", ExpiresAt: time.Now(), IsRevoked: false}
	mock.ExpectExec("INSERT INTO refresh_tokens").WillReturnError(sql.ErrConnDone)

	repo := NewAuthRepository(sqlxDB)
	err := repo.CreateRefreshToken(context.Background(), token)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
