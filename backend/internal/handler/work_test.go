package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB) {
	// Create test database
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(
		&model.User{},
		&model.WorkSession{},
		&model.Transaction{},
	)
	require.NoError(t, err)

	// Set test database globally (work handler uses database.GetDB())
	database.SetDB(db)

	// Create Fiber app
	app := fiber.New()

	return app, db
}

func createTestUser(t *testing.T, db *gorm.DB, balance float64) *model.User {
	timestamp := time.Now().UnixNano()
	googleID := fmt.Sprintf("test-google-id-%d", timestamp)
	user := &model.User{
		GoogleID: &googleID,
		Email:    fmt.Sprintf("test%d@example.com", timestamp),
		Username: fmt.Sprintf("testuser%d", timestamp),
		Name:     fmt.Sprintf("Test User %d", timestamp),
		Balance:  balance,
	}

	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

func TestWorkHandlerStartWork(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	// Set up route with middleware to inject userID
	app.Post("/work/start", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.StartWork(c)
	})

	// Make request
	req := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response, "data")
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["started_at"])
}

func TestWorkHandlerStartWorkUnauthorized(t *testing.T) {
	app, _ := setupTestApp(t)
	handler := NewWorkHandler()

	// Set up route without userID in locals (unauthorized)
	app.Post("/work/start", handler.StartWork)

	req := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkHandlerStartWorkAlreadyInProgress(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	app.Post("/work/start", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.StartWork(c)
	})

	// Start work once
	req1 := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Try to start again (should fail)
	req2 := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusConflict, resp2.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["error"].(bool))
	assert.Equal(t, "work session already in progress", response["message"])
}

func TestWorkHandlerGetStatus(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	app.Get("/work/status", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.GetStatus(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/work/status", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.False(t, data["is_working"].(bool))
}

func TestWorkHandlerGetStatusWorking(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	// Set up routes
	app.Post("/work/start", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.StartWork(c)
	})

	app.Get("/work/status", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.GetStatus(c)
	})

	// Start work
	startReq := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	_, err := app.Test(startReq)
	require.NoError(t, err)

	// Get status (should be working)
	statusReq := httptest.NewRequest(http.MethodGet, "/work/status", nil)
	resp, err := app.Test(statusReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.True(t, data["is_working"].(bool))
	assert.NotNil(t, data["started_at"])
	assert.NotNil(t, data["remaining_seconds"])
}

func TestWorkHandlerCompleteWorkTooEarly(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	// Set up routes
	app.Post("/work/start", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.StartWork(c)
	})

	app.Post("/work/complete", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.CompleteWork(c)
	})

	// Start work
	startReq := httptest.NewRequest(http.MethodPost, "/work/start", nil)
	_, err := app.Test(startReq)
	require.NoError(t, err)

	// Try to complete immediately (should fail - not enough time passed)
	completeReq := httptest.NewRequest(http.MethodPost, "/work/complete", nil)
	resp, err := app.Test(completeReq)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["error"].(bool))
	assert.Contains(t, response["message"], "work not completed yet")
}

func TestWorkHandlerCompleteWorkNoActiveSession(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	handler := NewWorkHandler()

	app.Post("/work/complete", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.CompleteWork(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/work/complete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["error"].(bool))
	assert.Equal(t, "no active work session", response["message"])
}

func TestWorkHandlerGetHistory(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	// Create some completed work sessions directly in DB (simulating past completed work)
	for i := 0; i < 3; i++ {
		session := &model.WorkSession{
			UserID:          user.ID,
			DurationSeconds: 180,
			Earned:          500.0,
			CompletedAt:     time.Now().Add(-time.Duration(i+1) * time.Hour),
		}
		err := db.Create(session).Error
		require.NoError(t, err)
	}

	handler := NewWorkHandler()

	app.Get("/work/history", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.GetHistory(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/work/history", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	sessions := data["sessions"].([]interface{})
	assert.Len(t, sessions, 3)
}

func TestWorkHandlerGetHistoryWithPagination(t *testing.T) {
	app, db := setupTestApp(t)
	user := createTestUser(t, db, 100.0)

	// Create 10 completed work sessions
	for i := 0; i < 10; i++ {
		session := &model.WorkSession{
			UserID:          user.ID,
			DurationSeconds: 180,
			Earned:          500.0,
			CompletedAt:     time.Now().Add(-time.Duration(i+1) * time.Hour),
		}
		err := db.Create(session).Error
		require.NoError(t, err)
	}

	handler := NewWorkHandler()

	app.Get("/work/history", func(c *fiber.Ctx) error {
		c.Locals("userID", user.ID)
		return handler.GetHistory(c)
	})

	// Request with limit=5, offset=3
	req := httptest.NewRequest(http.MethodGet, "/work/history?limit=5&offset=3", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	sessions := data["sessions"].([]interface{})
	assert.Len(t, sessions, 5)
}
