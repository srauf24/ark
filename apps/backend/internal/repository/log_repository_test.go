package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ark/internal/errs"
	testingPkg "ark/internal/testing"
)

// ========== Constructor Tests ==========

// Test 1: TestNewLogRepository_Success
func TestNewLogRepository_Success(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	repo := NewLogRepository(testDB.Pool)

	require.NotNil(t, repo, "Repository should not be nil")
	assert.NotNil(t, repo.db, "Database pool should be set")
}

// Test 2: TestLogRepository_StructFields
func TestLogRepository_StructFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	repo := NewLogRepository(testDB.Pool)

	// Verify struct has db field
	assert.NotNil(t, repo.db, "Repository should have db field")
}

// ========== GetByID Tests ==========

// Test 3: TestLogRepository_GetByID_WithTags
func TestLogRepository_GetByID_WithTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset first
	assetID := uuid.New()
	userID := "test-user-1"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert log with tags
	logID := uuid.New()
	tags := []string{"nginx", "ssl", "fix"}
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO asset_logs (id, asset_id, user_id, content, tags)
		VALUES ($1, $2, $3, $4, $5)
	`, logID, assetID, userID, "Fixed SSL certificate issue", tags)
	require.NoError(t, err)

	// Action: GetByID
	log, err := repo.GetByID(ctx, userID, logID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, logID, log.ID)
	assert.Equal(t, assetID, log.AssetID)
	assert.Equal(t, userID, log.UserID)
	assert.Equal(t, "Fixed SSL certificate issue", log.Content)
	assert.Equal(t, tags, log.Tags, "Tags should match")
	assert.NotZero(t, log.CreatedAt)
	assert.NotZero(t, log.UpdatedAt)
}

// Test 4: TestLogRepository_GetByID_NullTags
func TestLogRepository_GetByID_NullTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-2"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert log with NULL tags
	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO asset_logs (id, asset_id, user_id, content, tags)
		VALUES ($1, $2, $3, $4, NULL)
	`, logID, assetID, userID, "Log without tags")
	require.NoError(t, err)

	// Action: GetByID
	log, err := repo.GetByID(ctx, userID, logID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Log without tags", log.Content)
	assert.Nil(t, log.Tags, "NULL tags should be nil slice")
}

// Test 5: TestLogRepository_GetByID_EmptyTags
func TestLogRepository_GetByID_EmptyTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-3"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert log with empty tags array
	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO asset_logs (id, asset_id, user_id, content, tags)
		VALUES ($1, $2, $3, $4, $5)
	`, logID, assetID, userID, "Log with empty tags", []string{})
	require.NoError(t, err)

	// Action: GetByID
	log, err := repo.GetByID(ctx, userID, logID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Log with empty tags", log.Content)
	assert.NotNil(t, log.Tags, "Empty array should not be nil")
	assert.Len(t, log.Tags, 0, "Tags should be empty slice")
}

// Test 6: TestLogRepository_GetByID_NotFound
func TestLogRepository_GetByID_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Action: GetByID with random UUID (no log exists)
	randomID := uuid.New()
	log, err := repo.GetByID(ctx, "test-user-4", randomID)

	// Assert
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
	assert.Contains(t, err.Error(), "log not found")
}

// Test 7: TestLogRepository_GetByID_WrongUser
func TestLogRepository_GetByID_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	// Insert log for alice
	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO asset_logs (id, asset_id, user_id, content)
		VALUES ($1, $2, $3, $4)
	`, logID, assetID, "alice", "Alice's log")
	require.NoError(t, err)

	// Action: GetByID as bob (security test)
	log, err := repo.GetByID(ctx, "bob", logID)

	// Assert
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status, "Should be 404 (not revealing existence)")
	assert.Contains(t, err.Error(), "log not found")
}

// Test 8: TestLogRepository_GetByID_MultipleTags
func TestLogRepository_GetByID_MultipleTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-5"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert log with many tags
	logID := uuid.New()
	tags := []string{"production", "database", "postgresql", "backup", "automated"}
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO asset_logs (id, asset_id, user_id, content, tags)
		VALUES ($1, $2, $3, $4, $5)
	`, logID, assetID, userID, "Automated database backup completed", tags)
	require.NoError(t, err)

	// Action: GetByID
	log, err := repo.GetByID(ctx, userID, logID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, log.Tags, 5, "Should have 5 tags")
	assert.ElementsMatch(t, tags, log.Tags, "Tags should match (order may differ)")
}

// ========== ListByAsset Tests ==========

// Test 9: TestLogRepository_ListByAsset_NoFilters
func TestLogRepository_ListByAsset_NoFilters(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-9"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert 3 logs
	log1ID := uuid.New()
	log2ID := uuid.New()
	log3ID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`, log1ID, assetID, userID, "Log 1")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`, log2ID, assetID, userID, "Log 2")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`, log3ID, assetID, userID, "Log 3")
	require.NoError(t, err)

	// Action: List with no filters
	params := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 3, "Should return 3 logs")
	assert.Equal(t, "Log 3", logs[0].Content, "Should be ordered by created_at desc")
}

// Test 10: TestLogRepository_ListByAsset_SingleTag
func TestLogRepository_ListByAsset_SingleTag(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-10"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different tags
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Nginx fix", []string{"nginx", "fix"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Docker issue", []string{"docker", "issue"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Nginx config", []string{"nginx", "config"})
	require.NoError(t, err)

	// Action: Filter by single tag "nginx"
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		Tags:      []string{"nginx"},
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 2, "Should return 2 logs with 'nginx' tag")
	assert.Contains(t, logs[0].Content, "Nginx", "Should contain nginx logs")
	assert.Contains(t, logs[1].Content, "Nginx", "Should contain nginx logs")
}

// Test 11: TestLogRepository_ListByAsset_MultipleTags_AND
func TestLogRepository_ListByAsset_MultipleTags_AND(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-11"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different tag combinations
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Log with nginx and fix", []string{"nginx", "fix", "ssl"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Log with only nginx", []string{"nginx"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Log with nginx and fix only", []string{"nginx", "fix"})
	require.NoError(t, err)

	// Action: Filter by multiple tags ["nginx", "fix"] - AND logic
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		Tags:      []string{"nginx", "fix"},
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert: Should return logs that have BOTH nginx AND fix
	require.NoError(t, err)
	assert.Len(t, logs, 2, "Should return 2 logs with BOTH 'nginx' AND 'fix' tags")
	for _, log := range logs {
		assert.Contains(t, log.Tags, "nginx", "Should contain nginx tag")
		assert.Contains(t, log.Tags, "fix", "Should contain fix tag")
	}
}

// Test 12: TestLogRepository_ListByAsset_Search
func TestLogRepository_ListByAsset_Search(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-12"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different content
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "Fixed SSL certificate issue")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "Updated nginx configuration")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "SSL handshake errors resolved")
	require.NoError(t, err)

	// Action: Search for "ssl" (case-insensitive)
	search := "ssl"
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		Search:    &search,
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 2, "Should return 2 logs containing 'ssl'")
	assert.Contains(t, strings.ToLower(logs[0].Content), "ssl")
	assert.Contains(t, strings.ToLower(logs[1].Content), "ssl")
}

// Test 13: TestLogRepository_ListByAsset_DateRange_StartDate
func TestLogRepository_ListByAsset_DateRange_StartDate(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-13"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with specific timestamps
	oldDate := time.Now().Add(-72 * time.Hour) // 3 days ago
	newDate := time.Now().Add(-24 * time.Hour) // 1 day ago

	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Old log", oldDate)
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Recent log", newDate)
	require.NoError(t, err)

	// Action: Filter by start_date (last 2 days)
	startDate := time.Now().Add(-48 * time.Hour)
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		StartDate: &startDate,
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 1, "Should return 1 log after start_date")
	assert.Equal(t, "Recent log", logs[0].Content)
}

// Test 14: TestLogRepository_ListByAsset_DateRange_EndDate
func TestLogRepository_ListByAsset_DateRange_EndDate(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-14"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with specific timestamps
	oldDate := time.Now().Add(-72 * time.Hour)
	newDate := time.Now().Add(-24 * time.Hour)

	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Old log", oldDate)
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Recent log", newDate)
	require.NoError(t, err)

	// Action: Filter by end_date (before 2 days ago)
	endDate := time.Now().Add(-48 * time.Hour)
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		EndDate:   &endDate,
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 1, "Should return 1 log before end_date")
	assert.Equal(t, "Old log", logs[0].Content)
}

// Test 15: TestLogRepository_ListByAsset_DateRange_Both
func TestLogRepository_ListByAsset_DateRange_Both(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-15"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs across 5 days
	dates := []time.Time{
		time.Now().Add(-96 * time.Hour), // 4 days ago
		time.Now().Add(-72 * time.Hour), // 3 days ago
		time.Now().Add(-48 * time.Hour), // 2 days ago
		time.Now().Add(-24 * time.Hour), // 1 day ago
		time.Now(),                       // now
	}

	for i, date := range dates {
		_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(), assetID, userID, fmt.Sprintf("Log %d", i+1), date)
		require.NoError(t, err)
	}

	// Action: Filter by start_date AND end_date (days 2-4)
	startDate := time.Now().Add(-80 * time.Hour)
	endDate := time.Now().Add(-36 * time.Hour)
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		StartDate: &startDate,
		EndDate:   &endDate,
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 2, "Should return 2 logs within date range")
	assert.Equal(t, "Log 2", logs[0].Content)
	assert.Equal(t, "Log 3", logs[1].Content)
}

// Test 16: TestLogRepository_ListByAsset_CombinedFilters
func TestLogRepository_ListByAsset_CombinedFilters(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-16"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with various attributes
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), assetID, userID, "Fixed SSL nginx issue", []string{"nginx", "ssl", "fix"}, time.Now().Add(-24*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), assetID, userID, "SSL certificate renewed", []string{"ssl", "cert"}, time.Now().Add(-12*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), assetID, userID, "Nginx SSL configuration", []string{"nginx", "ssl"}, time.Now().Add(-6*time.Hour))
	require.NoError(t, err)

	// Action: Filter by tags=["nginx","ssl"] AND search="fixed" AND start_date
	search := "fixed"
	startDate := time.Now().Add(-48 * time.Hour)
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		Tags:      []string{"nginx", "ssl"},
		Search:    &search,
		StartDate: &startDate,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 1, "Should return 1 log matching all filters")
	assert.Equal(t, "Fixed SSL nginx issue", logs[0].Content)
	assert.Contains(t, logs[0].Tags, "nginx")
	assert.Contains(t, logs[0].Tags, "ssl")
}

// Test 17: TestLogRepository_ListByAsset_Pagination
func TestLogRepository_ListByAsset_Pagination(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-17"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert 5 logs
	for i := 1; i <= 5; i++ {
		_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
			uuid.New(), assetID, userID, fmt.Sprintf("Log %d", i))
		require.NoError(t, err)
	}

	// Action: Get page 1 (limit=2, offset=0)
	params1 := &model.LogQueryParams{Limit: 2, Offset: 0, SortBy: "created_at", SortOrder: "asc"}
	logs1, err := repo.ListByAsset(ctx, userID, assetID, params1)
	require.NoError(t, err)
	assert.Len(t, logs1, 2, "Page 1 should have 2 logs")
	assert.Equal(t, "Log 1", logs1[0].Content)

	// Action: Get page 2 (limit=2, offset=2)
	params2 := &model.LogQueryParams{Limit: 2, Offset: 2, SortBy: "created_at", SortOrder: "asc"}
	logs2, err := repo.ListByAsset(ctx, userID, assetID, params2)
	require.NoError(t, err)
	assert.Len(t, logs2, 2, "Page 2 should have 2 logs")
	assert.Equal(t, "Log 3", logs2[0].Content)

	// Action: Get page 3 (limit=2, offset=4)
	params3 := &model.LogQueryParams{Limit: 2, Offset: 4, SortBy: "created_at", SortOrder: "asc"}
	logs3, err := repo.ListByAsset(ctx, userID, assetID, params3)
	require.NoError(t, err)
	assert.Len(t, logs3, 1, "Page 3 should have 1 log")
	assert.Equal(t, "Log 5", logs3[0].Content)
}

// Test 18: TestLogRepository_ListByAsset_WrongUser
func TestLogRepository_ListByAsset_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset for alice
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	// Insert log for alice
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, "alice", "Alice's log")
	require.NoError(t, err)

	// Action: Try to list as bob (security test)
	params := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logs, err := repo.ListByAsset(ctx, "bob", assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 0, "Bob should see 0 logs (user isolation)")
}

// Test 19: TestLogRepository_ListByAsset_WrongAsset
func TestLogRepository_ListByAsset_WrongAsset(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create 2 assets for same user
	asset1ID := uuid.New()
	asset2ID := uuid.New()
	userID := "test-user-19"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, asset1ID, userID, "Asset 1")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, asset2ID, userID, "Asset 2")
	require.NoError(t, err)

	// Insert log for asset1
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), asset1ID, userID, "Log for asset 1")
	require.NoError(t, err)

	// Action: Try to list logs for asset2 (should be empty)
	params := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logs, err := repo.ListByAsset(ctx, userID, asset2ID, params)

	// Assert
	require.NoError(t, err)
	assert.Len(t, logs, 0, "Asset 2 should have 0 logs")
}

// Test 20: TestLogRepository_ListByAsset_Sorting
func TestLogRepository_ListByAsset_Sorting(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-20"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different timestamps
	log1ID := uuid.New()
	log2ID := uuid.New()
	log3ID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		log1ID, assetID, userID, "Log 1", time.Now().Add(-48*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		log2ID, assetID, userID, "Log 2", time.Now().Add(-24*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		log3ID, assetID, userID, "Log 3", time.Now())
	require.NoError(t, err)

	// Test sorting: created_at ASC
	paramsAsc := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "asc"}
	logsAsc, err := repo.ListByAsset(ctx, userID, assetID, paramsAsc)
	require.NoError(t, err)
	assert.Equal(t, "Log 1", logsAsc[0].Content, "ASC: oldest first")
	assert.Equal(t, "Log 3", logsAsc[2].Content, "ASC: newest last")

	// Test sorting: created_at DESC
	paramsDesc := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logsDesc, err := repo.ListByAsset(ctx, userID, assetID, paramsDesc)
	require.NoError(t, err)
	assert.Equal(t, "Log 3", logsDesc[0].Content, "DESC: newest first")
	assert.Equal(t, "Log 1", logsDesc[2].Content, "DESC: oldest last")
}

// Test 21: TestLogRepository_ListByAsset_InvalidSortBy
func TestLogRepository_ListByAsset_InvalidSortBy(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-21"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Try to sort by invalid column (SQL injection attempt)
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		SortBy:    "id; DROP TABLE asset_logs;",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	assert.Nil(t, logs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_by", "Should reject invalid sort_by")
}

// Test 22: TestLogRepository_ListByAsset_InvalidSortOrder
func TestLogRepository_ListByAsset_InvalidSortOrder(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-22"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Try to use invalid sort order
	params := &model.LogQueryParams{
		Limit:     10,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "INVALID",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	assert.Nil(t, logs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_order", "Should reject invalid sort_order")
}

// Test 23: TestLogRepository_ListByAsset_EmptyResult
func TestLogRepository_ListByAsset_EmptyResult(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset with no logs
	assetID := uuid.New()
	userID := "test-user-23"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: List logs for empty asset
	params := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, logs, "Should return non-nil slice")
	assert.Len(t, logs, 0, "Should return empty slice")
}

// ========== CountByAsset Tests ==========

// Test 24: TestLogRepository_CountByAsset_NoFilters
func TestLogRepository_CountByAsset_NoFilters(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-24"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert 5 logs
	for i := 1; i <= 5; i++ {
		_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
			uuid.New(), assetID, userID, fmt.Sprintf("Log %d", i))
		require.NoError(t, err)
	}

	// Action: Count with no filters
	params := &model.LogQueryParams{}
	count, err := repo.CountByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(5), count, "Should count 5 logs")
}

// Test 25: TestLogRepository_CountByAsset_WithTagFilter
func TestLogRepository_CountByAsset_WithTagFilter(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-25"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different tags
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Nginx fix", []string{"nginx", "fix"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Docker issue", []string{"docker", "issue"})
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Nginx config", []string{"nginx", "config"})
	require.NoError(t, err)

	// Action: Count logs with "nginx" tag
	params := &model.LogQueryParams{Tags: []string{"nginx"}}
	count, err := repo.CountByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "Should count 2 logs with 'nginx' tag")
}

// Test 26: TestLogRepository_CountByAsset_WithSearch
func TestLogRepository_CountByAsset_WithSearch(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-26"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with different content
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "Fixed SSL certificate issue")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "Updated nginx configuration")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		uuid.New(), assetID, userID, "SSL handshake errors resolved")
	require.NoError(t, err)

	// Action: Count logs containing "ssl"
	search := "ssl"
	params := &model.LogQueryParams{Search: &search}
	count, err := repo.CountByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "Should count 2 logs containing 'ssl'")
}

// Test 27: TestLogRepository_CountByAsset_WithDateRange
func TestLogRepository_CountByAsset_WithDateRange(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-27"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs across different dates
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Old log", time.Now().Add(-96*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Recent log 1", time.Now().Add(-24*time.Hour))
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), assetID, userID, "Recent log 2", time.Now().Add(-12*time.Hour))
	require.NoError(t, err)

	// Action: Count logs within last 48 hours
	startDate := time.Now().Add(-48 * time.Hour)
	params := &model.LogQueryParams{StartDate: &startDate}
	count, err := repo.CountByAsset(ctx, userID, assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "Should count 2 logs within last 48 hours")
}

// Test 28: TestLogRepository_CountByAsset_MatchesList
func TestLogRepository_CountByAsset_MatchesList(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-28"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Insert logs with various attributes
	for i := 1; i <= 10; i++ {
		tags := []string{"test"}
		if i%2 == 0 {
			tags = append(tags, "nginx")
		}
		_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(), assetID, userID, fmt.Sprintf("Log %d", i), tags)
		require.NoError(t, err)
	}

	// Action: List and Count with same filter (nginx tag)
	params := &model.LogQueryParams{
		Limit:     100,
		Offset:    0,
		Tags:      []string{"nginx"},
		SortBy:    "created_at",
		SortOrder: "asc",
	}
	logs, err := repo.ListByAsset(ctx, userID, assetID, params)
	require.NoError(t, err)

	countParams := &model.LogQueryParams{Tags: []string{"nginx"}}
	count, err := repo.CountByAsset(ctx, userID, assetID, countParams)
	require.NoError(t, err)

	// Assert: Count should match List length
	assert.Equal(t, int64(len(logs)), count, "Count should match List length")
	assert.Equal(t, int64(5), count, "Should have 5 logs with 'nginx' tag")
}

// Test 29: TestLogRepository_CountByAsset_WrongUser
func TestLogRepository_CountByAsset_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset for alice
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	// Insert logs for alice
	for i := 1; i <= 3; i++ {
		_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
			uuid.New(), assetID, "alice", fmt.Sprintf("Alice's log %d", i))
		require.NoError(t, err)
	}

	// Action: Try to count as bob (security test)
	params := &model.LogQueryParams{}
	count, err := repo.CountByAsset(ctx, "bob", assetID, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Bob should see 0 logs (user isolation)")
}

// ========== Create Tests ==========

// Test 30: TestLogRepository_Create_AllFields
func TestLogRepository_Create_AllFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-30"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Create log with all fields
	req := &model.CreateLogRequest{
		Content: "Fixed nginx SSL certificate issue",
		Tags:    []string{"nginx", "ssl", "fix"},
	}
	log, err := repo.Create(ctx, userID, assetID, req)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, log.ID, "Should generate UUID")
	assert.Equal(t, assetID, log.AssetID)
	assert.Equal(t, userID, log.UserID)
	assert.Equal(t, "Fixed nginx SSL certificate issue", log.Content)
	assert.Equal(t, []string{"nginx", "ssl", "fix"}, log.Tags)
	assert.NotZero(t, log.CreatedAt)
	assert.NotZero(t, log.UpdatedAt)
}

// Test 31: TestLogRepository_Create_MinimalFields
func TestLogRepository_Create_MinimalFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-31"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Create log with only content (no tags)
	req := &model.CreateLogRequest{
		Content: "Updated configuration",
	}
	log, err := repo.Create(ctx, userID, assetID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Updated configuration", log.Content)
	assert.Nil(t, log.Tags, "Nil tags should remain nil")
}

// Test 32: TestLogRepository_Create_EmptyTags
func TestLogRepository_Create_EmptyTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-32"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Create log with empty tags array
	req := &model.CreateLogRequest{
		Content: "Log with empty tags",
		Tags:    []string{},
	}
	log, err := repo.Create(ctx, userID, assetID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Log with empty tags", log.Content)
	assert.NotNil(t, log.Tags, "Empty array should not be nil")
	assert.Len(t, log.Tags, 0, "Should be empty slice")
}

// Test 33: TestLogRepository_Create_AssetNotFound
func TestLogRepository_Create_AssetNotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Action: Try to create log for non-existent asset
	randomAssetID := uuid.New()
	req := &model.CreateLogRequest{
		Content: "Log for non-existent asset",
	}
	log, err := repo.Create(ctx, "test-user-33", randomAssetID, req)

	// Assert
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
	assert.Contains(t, err.Error(), "asset not found")
}

// Test 34: TestLogRepository_Create_WrongUser
func TestLogRepository_Create_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset for alice
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	// Action: Try to create log as bob (security test)
	req := &model.CreateLogRequest{
		Content: "Bob trying to create log for Alice's asset",
	}
	log, err := repo.Create(ctx, "bob", assetID, req)

	// Assert: Should succeed because user_id is what we provide, not validated by FK
	// The asset_id FK only checks that asset exists, not that user owns it
	// This means bob can create a log with bob's user_id for alice's asset
	// However, bob won't be able to see it via ListByAsset because of dual scoping
	require.NoError(t, err)
	assert.Equal(t, "bob", log.UserID)
	assert.Equal(t, assetID, log.AssetID)

	// Verify bob cannot see this log via ListByAsset (security check)
	params := &model.LogQueryParams{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "desc"}
	logs, err := repo.ListByAsset(ctx, "bob", assetID, params)
	require.NoError(t, err)
	assert.Len(t, logs, 0, "Bob should not see logs (wrong asset owner)")

	// Verify alice cannot see bob's log either (dual scoping)
	logsAlice, err := repo.ListByAsset(ctx, "alice", assetID, params)
	require.NoError(t, err)
	assert.Len(t, logsAlice, 0, "Alice should not see bob's log")
}

// Test 35: TestLogRepository_Create_UniqueIDs
func TestLogRepository_Create_UniqueIDs(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-35"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Create multiple logs
	ids := make([]uuid.UUID, 0, 3)
	for i := 1; i <= 3; i++ {
		req := &model.CreateLogRequest{
			Content: fmt.Sprintf("Log %d", i),
		}
		log, err := repo.Create(ctx, userID, assetID, req)
		require.NoError(t, err)
		ids = append(ids, log.ID)
	}

	// Assert: All IDs should be unique
	uniqueIDs := make(map[uuid.UUID]bool)
	for _, id := range ids {
		assert.False(t, uniqueIDs[id], "IDs should be unique")
		uniqueIDs[id] = true
	}
	assert.Len(t, uniqueIDs, 3)
}

// Test 36: TestLogRepository_Create_TimestampsAutoPopulated
func TestLogRepository_Create_TimestampsAutoPopulated(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset
	assetID := uuid.New()
	userID := "test-user-36"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	// Action: Create log
	before := time.Now().Add(-1 * time.Second)
	req := &model.CreateLogRequest{
		Content: "Test log",
	}
	log, err := repo.Create(ctx, userID, assetID, req)
	after := time.Now().Add(1 * time.Second)

	// Assert
	require.NoError(t, err)
	assert.True(t, log.CreatedAt.After(before) && log.CreatedAt.Before(after), "CreatedAt should be auto-populated")
	assert.True(t, log.UpdatedAt.After(before) && log.UpdatedAt.Before(after), "UpdatedAt should be auto-populated")
	assert.Equal(t, log.CreatedAt, log.UpdatedAt, "CreatedAt and UpdatedAt should be equal on creation")
}

// ========== Update Tests ==========

// Test 37: TestLogRepository_Update_ContentOnly
func TestLogRepository_Update_ContentOnly(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-37"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Original content", []string{"original", "tags"})
	require.NoError(t, err)

	// Action: Update content only
	newContent := "Updated content"
	req := &model.UpdateLogRequest{
		Content: &newContent,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Updated content", log.Content, "Content should be updated")
	assert.Equal(t, []string{"original", "tags"}, log.Tags, "Tags should remain unchanged")
}

// Test 38: TestLogRepository_Update_TagsOnly
func TestLogRepository_Update_TagsOnly(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-38"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Original content", []string{"original", "tags"})
	require.NoError(t, err)

	// Action: Update tags only
	newTags := []string{"nginx", "ssl", "fix"}
	req := &model.UpdateLogRequest{
		Tags: &newTags,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Original content", log.Content, "Content should remain unchanged")
	assert.Equal(t, []string{"nginx", "ssl", "fix"}, log.Tags, "Tags should be updated")
}

// Test 39: TestLogRepository_Update_BothFields
func TestLogRepository_Update_BothFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-39"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Original content", []string{"original"})
	require.NoError(t, err)

	// Action: Update both content and tags
	newContent := "Updated content"
	newTags := []string{"docker", "container"}
	req := &model.UpdateLogRequest{
		Content: &newContent,
		Tags:    &newTags,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Updated content", log.Content)
	assert.Equal(t, []string{"docker", "container"}, log.Tags)
}

// Test 40: TestLogRepository_Update_ClearTags
func TestLogRepository_Update_ClearTags(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log with tags
	assetID := uuid.New()
	userID := "test-user-40"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Log with tags", []string{"nginx", "ssl"})
	require.NoError(t, err)

	// Action: Clear tags by setting empty array
	emptyTags := []string{}
	req := &model.UpdateLogRequest{
		Tags: &emptyTags,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, log.Tags, "Empty array should not be nil")
	assert.Len(t, log.Tags, 0, "Tags should be empty")
}

// Test 41: TestLogRepository_Update_AllNil
func TestLogRepository_Update_AllNil(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-41"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Original content", []string{"original"})
	require.NoError(t, err)

	// Get original updated_at
	var originalUpdatedAt time.Time
	err = testDB.Pool.QueryRow(ctx, `SELECT updated_at FROM asset_logs WHERE id = $1`, logID).Scan(&originalUpdatedAt)
	require.NoError(t, err)

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Action: Update with all nil fields (only updated_at should change)
	req := &model.UpdateLogRequest{}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Original content", log.Content, "Content should remain unchanged")
	assert.Equal(t, []string{"original"}, log.Tags, "Tags should remain unchanged")
	assert.True(t, log.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
}

// Test 42: TestLogRepository_Update_NotFound
func TestLogRepository_Update_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Action: Try to update non-existent log
	randomID := uuid.New()
	newContent := "Updated content"
	req := &model.UpdateLogRequest{
		Content: &newContent,
	}
	log, err := repo.Update(ctx, "test-user-42", randomID, req)

	// Assert
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
	assert.Contains(t, err.Error(), "log not found")
}

// Test 43: TestLogRepository_Update_WrongUser
func TestLogRepository_Update_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log for alice
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		logID, assetID, "alice", "Alice's log")
	require.NoError(t, err)

	// Action: Try to update as bob (security test)
	newContent := "Bob trying to update Alice's log"
	req := &model.UpdateLogRequest{
		Content: &newContent,
	}
	log, err := repo.Update(ctx, "bob", logID, req)

	// Assert
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status, "Should be 404 (not revealing existence)")
	assert.Contains(t, err.Error(), "log not found")
}

// Test 44: TestLogRepository_Update_UpdatedAtChanges
func TestLogRepository_Update_UpdatedAtChanges(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-44"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Original content", time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	// Get original log
	originalLog, err := repo.GetByID(ctx, userID, logID)
	require.NoError(t, err)

	// Wait to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Action: Update content
	newContent := "Updated content"
	req := &model.UpdateLogRequest{
		Content: &newContent,
	}
	updatedLog, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, originalLog.CreatedAt, updatedLog.CreatedAt, "CreatedAt should not change")
	assert.True(t, updatedLog.UpdatedAt.After(originalLog.UpdatedAt), "UpdatedAt should be newer")
}

// Test 45: TestLogRepository_Update_TagsNullToArray
func TestLogRepository_Update_TagsNullToArray(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log with NULL tags
	assetID := uuid.New()
	userID := "test-user-45"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, NULL)`,
		logID, assetID, userID, "Log without tags")
	require.NoError(t, err)

	// Action: Update tags from NULL to array
	newTags := []string{"nginx", "fix"}
	req := &model.UpdateLogRequest{
		Tags: &newTags,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, log.Tags, "Tags should not be nil")
	assert.Equal(t, []string{"nginx", "fix"}, log.Tags)
}

// Test 46: TestLogRepository_Update_TagsArrayToEmpty
func TestLogRepository_Update_TagsArrayToEmpty(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log with tags
	assetID := uuid.New()
	userID := "test-user-46"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content, tags) VALUES ($1, $2, $3, $4, $5)`,
		logID, assetID, userID, "Log with tags", []string{"nginx", "ssl", "fix"})
	require.NoError(t, err)

	// Action: Update tags to empty array
	emptyTags := []string{}
	req := &model.UpdateLogRequest{
		Tags: &emptyTags,
	}
	log, err := repo.Update(ctx, userID, logID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, log.Tags, "Empty array should not be nil")
	assert.Len(t, log.Tags, 0, "Tags should be empty")
}

// ========== Delete Tests ==========

// Test 47: TestLogRepository_Delete_Success
func TestLogRepository_Delete_Success(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log
	assetID := uuid.New()
	userID := "test-user-47"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Test Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		logID, assetID, userID, "Log to delete")
	require.NoError(t, err)

	// Action: Delete log
	err = repo.Delete(ctx, userID, logID)

	// Assert
	require.NoError(t, err)

	// Verify log is actually deleted
	log, err := repo.GetByID(ctx, userID, logID)
	assert.Nil(t, log)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
}

// Test 48: TestLogRepository_Delete_NotFound
func TestLogRepository_Delete_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Action: Try to delete non-existent log
	randomID := uuid.New()
	err := repo.Delete(ctx, "test-user-48", randomID)

	// Assert
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
	assert.Contains(t, err.Error(), "log not found")
}

// Test 49: TestLogRepository_Delete_WrongUser
func TestLogRepository_Delete_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewLogRepository(testDB.Pool)

	// Create asset and log for alice
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	logID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO asset_logs (id, asset_id, user_id, content) VALUES ($1, $2, $3, $4)`,
		logID, assetID, "alice", "Alice's log")
	require.NoError(t, err)

	// Action: Try to delete as bob (security test)
	err = repo.Delete(ctx, "bob", logID)

	// Assert
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status, "Should be 404 (not revealing existence)")
	assert.Contains(t, err.Error(), "log not found")

	// Verify log still exists for alice
	log, err := repo.GetByID(ctx, "alice", logID)
	require.NoError(t, err)
	assert.Equal(t, "Alice's log", log.Content, "Log should still exist")
}
