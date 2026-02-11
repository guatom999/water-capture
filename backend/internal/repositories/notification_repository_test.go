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

func setupNotificationMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, NotificationRepositoryInterface) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewNotificationRepository(sqlxDB)
	return sqlxDB, mock, repo
}

func TestNotification_GetSubscriptionsByLocationID_Success(t *testing.T) {
	sqlxDB, mock, _ := setupNotificationMockDB(t)
	defer sqlxDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "location_id", "channel", "target", "threshold_level", "is_active", "created_at", "updated_at"}).AddRow(1, 10, 5, "line", "@user", 3.0, true, now, now)

	mock.ExpectQuery("SELECT id, user_id").WithArgs(5).WillReturnRows(rows)

	repo := NewNotificationRepository(sqlxDB)
	result, err := repo.GetSubscriptionsByLocationID(context.Background(), 5)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	if result[0].UserID != nil {
		assert.Equal(t, 10, *result[0].UserID)
	} else {
		t.Fatalf("expected UserID to be non-nil")
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotification_GetSubscriptionsByLocationID_Error(t *testing.T) {
	sqlxDB, mock, _ := setupNotificationMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("SELECT id, user_id").WillReturnError(sql.ErrConnDone)

	repo := NewNotificationRepository(sqlxDB)
	result, err := repo.GetSubscriptionsByLocationID(context.Background(), 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotification_CreateSubscription_Success(t *testing.T) {
	sqlxDB, mock, _ := setupNotificationMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("INSERT INTO notification_subscriptions").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	repo := NewNotificationRepository(sqlxDB)
	one := 1
	sub := &entities.NotificationSubscription{UserID: &one, LocationID: nil, Channel: "line", Target: "@user", ThresholdLevel: 2.5, IsActive: true}
	err := repo.CreateSubscription(context.Background(), sub)
	assert.NoError(t, err)
	assert.Equal(t, 7, sub.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotification_LogNotification_Success(t *testing.T) {
	sqlxDB, mock, _ := setupNotificationMockDB(t)
	defer sqlxDB.Close()

	mock.ExpectQuery("INSERT INTO notification_logs").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	repo := NewNotificationRepository(sqlxDB)
	oneInt := 1
	lg := &entities.NotificationLog{SubscriptionID: &oneInt, LocationID: 5, WaterLevel: 1.5, Message: "msg", Channel: "line", Status: "SENT", ErrorMessage: nil}
	err := repo.LogNotification(context.Background(), lg)
	assert.NoError(t, err)
	assert.Equal(t, 3, lg.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
