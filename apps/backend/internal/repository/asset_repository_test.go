package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ark/internal/errs"
	"ark/internal/model"
	testingPkg "ark/internal/testing"
)

// ========== Constructor Tests ==========

// Test 1: TestNewAssetRepository_Success
func TestNewAssetRepository_Success(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	repo := NewAssetRepository(testDB.Pool)

	require.NotNil(t, repo, "Repository should not be nil")
	assert.NotNil(t, repo.db, "Database pool should be set")
}

// Test 2: TestAssetRepository_StructFields
func TestAssetRepository_StructFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	repo := NewAssetRepository(testDB.Pool)

	// Verify struct has db field
	assert.NotNil(t, repo.db, "Repository should have db field")
}

// ========== GetByID Tests ==========

// Test 3: TestAssetRepository_GetByID_Success
func TestAssetRepository_GetByID_Success(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Insert test asset with all fields
	assetID := uuid.New()
	userID := "test-user-1"
	metadata := json.RawMessage(`{"cpu": "8 cores", "ram": "32GB"}`)

	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO assets (id, user_id, name, type, hostname, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, assetID, userID, "Test Server", "server", "test.example.com", metadata)
	require.NoError(t, err, "Failed to insert test asset")

	// Action: GetByID
	asset, err := repo.GetByID(ctx, userID, assetID)

	// Assert
	require.NoError(t, err, "GetByID should succeed")
	assert.Equal(t, assetID, asset.ID, "ID should match")
	assert.Equal(t, userID, asset.UserID, "UserID should match")
	assert.Equal(t, "Test Server", asset.Name, "Name should match")
	assert.NotNil(t, asset.Type, "Type should not be nil")
	assert.Equal(t, "server", *asset.Type, "Type should match")
	assert.NotNil(t, asset.Hostname, "Hostname should not be nil")
	assert.Equal(t, "test.example.com", *asset.Hostname, "Hostname should match")
	assert.NotNil(t, asset.Metadata, "Metadata should not be nil")
	assert.JSONEq(t, string(metadata), string(asset.Metadata), "Metadata should match")
	assert.NotZero(t, asset.CreatedAt, "CreatedAt should be set")
	assert.NotZero(t, asset.UpdatedAt, "UpdatedAt should be set")
}

// Test 4: TestAssetRepository_GetByID_WithNullFields
func TestAssetRepository_GetByID_WithNullFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Insert asset with NULL type, hostname, metadata
	assetID := uuid.New()
	userID := "test-user-2"

	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO assets (id, user_id, name, type, hostname, metadata)
		VALUES ($1, $2, $3, NULL, NULL, NULL)
	`, assetID, userID, "Minimal Asset")
	require.NoError(t, err, "Failed to insert test asset")

	// Action: GetByID
	asset, err := repo.GetByID(ctx, userID, assetID)

	// Assert
	require.NoError(t, err, "GetByID should succeed")
	assert.Equal(t, "Minimal Asset", asset.Name, "Name should match")
	assert.Nil(t, asset.Type, "Type should be nil")
	assert.Nil(t, asset.Hostname, "Hostname should be nil")
	assert.Nil(t, asset.Metadata, "Metadata should be nil")
}

// Test 5: TestAssetRepository_GetByID_NotFound
func TestAssetRepository_GetByID_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Action: GetByID with random UUID (no asset exists)
	randomID := uuid.New()
	asset, err := repo.GetByID(ctx, "test-user-3", randomID)

	// Assert
	assert.Nil(t, asset, "Asset should be nil")
	assert.Error(t, err, "Should return error")

	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr, "Error should be HTTPError")
	assert.Equal(t, 404, httpErr.Status, "Status should be 404")
	assert.Contains(t, err.Error(), "asset not found", "Error message should indicate not found")
}

// Test 6: TestAssetRepository_GetByID_WrongUser
func TestAssetRepository_GetByID_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Insert asset for user "alice"
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO assets (id, user_id, name)
		VALUES ($1, $2, $3)
	`, assetID, "alice", "Alice's Server")
	require.NoError(t, err, "Failed to insert test asset")

	// Action: GetByID as user "bob" (security test)
	asset, err := repo.GetByID(ctx, "bob", assetID)

	// Assert
	assert.Nil(t, asset, "Asset should be nil")
	assert.Error(t, err, "Should return error")

	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr, "Error should be HTTPError")
	assert.Equal(t, 404, httpErr.Status, "Status should be 404 (not revealing existence)")
	assert.Contains(t, err.Error(), "asset not found", "Error message should be same as non-existent")
}

// Test 7: TestAssetRepository_GetByID_ContextCanceled
func TestAssetRepository_GetByID_ContextCanceled(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	repo := NewAssetRepository(testDB.Pool)

	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Action: GetByID with canceled context
	asset, err := repo.GetByID(ctx, "test-user-4", uuid.New())

	// Assert
	assert.Nil(t, asset, "Asset should be nil")
	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "context canceled", "Error should indicate context cancellation")
}

// Test 8: TestAssetRepository_GetByID_ComplexMetadata
func TestAssetRepository_GetByID_ComplexMetadata(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Insert asset with complex nested metadata
	assetID := uuid.New()
	userID := "test-user-5"
	complexMetadata := json.RawMessage(`{
		"hardware": {
			"cpu": {"cores": 16, "model": "AMD EPYC"},
			"ram": {"size": "128GB", "type": "DDR4"},
			"storage": [
				{"type": "NVMe", "size": "1TB"},
				{"type": "SSD", "size": "4TB"}
			]
		},
		"network": {
			"interfaces": ["eth0", "eth1"],
			"ip_addresses": ["192.168.1.100", "10.0.0.50"]
		}
	}`)

	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO assets (id, user_id, name, metadata)
		VALUES ($1, $2, $3, $4)
	`, assetID, userID, "Complex Server", complexMetadata)
	require.NoError(t, err, "Failed to insert test asset")

	// Action: GetByID
	asset, err := repo.GetByID(ctx, userID, assetID)

	// Assert
	require.NoError(t, err, "GetByID should succeed")
	assert.NotNil(t, asset.Metadata, "Metadata should not be nil")

	// Verify metadata unmarshals correctly
	var retrievedMeta map[string]interface{}
	err = json.Unmarshal(asset.Metadata, &retrievedMeta)
	require.NoError(t, err, "Metadata should unmarshal")
	assert.Contains(t, retrievedMeta, "hardware", "Metadata should contain hardware")
	assert.Contains(t, retrievedMeta, "network", "Metadata should contain network")
}

// ========== List Tests ==========

// Test 9: TestAssetRepository_List_NoFilters
func TestAssetRepository_List_NoFilters(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "alice"
	// Insert 3 assets for alice, 2 for bob
	for i := 1; i <= 3; i++ {
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, userID, fmt.Sprintf("Asset %d", i))
		require.NoError(t, err)
	}
	for i := 1; i <= 2; i++ {
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, "bob", fmt.Sprintf("Bob Asset %d", i))
		require.NoError(t, err)
	}

	params := &model.AssetQueryParams{
		Limit:     10,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	assets, err := repo.List(ctx, userID, params)

	require.NoError(t, err)
	assert.Len(t, assets, 3, "Should return only alice's assets")
	for _, asset := range assets {
		assert.Equal(t, userID, asset.UserID, "All assets should belong to alice")
	}
}

// Test 10: TestAssetRepository_List_TypeFilter
func TestAssetRepository_List_TypeFilter(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	// Insert assets with different types
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, "Server 1", "server")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, "Server 2", "server")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, "NAS", "nas")
	require.NoError(t, err)

	params := &model.AssetQueryParams{
		Type:      testingPkg.Ptr("server"),
		Limit:     10,
		Offset:    0,
		SortBy:    "name",
		SortOrder: "asc",
	}

	assets, err := repo.List(ctx, userID, params)

	require.NoError(t, err)
	assert.Len(t, assets, 2, "Should return only server assets")
	for _, asset := range assets {
		assert.Equal(t, "server", *asset.Type)
	}
}

// Test 11: TestAssetRepository_List_SearchByName
func TestAssetRepository_List_SearchByName(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, userID, "prod-server")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, userID, "dev-server")
	require.NoError(t, err)
	_, err = testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, userID, "nas-backup")
	require.NoError(t, err)

	params := &model.AssetQueryParams{
		Search:    testingPkg.Ptr("server"),
		Limit:     10,
		Offset:    0,
		SortBy:    "name",
		SortOrder: "asc",
	}

	assets, err := repo.List(ctx, userID, params)

	require.NoError(t, err)
	assert.Len(t, assets, 2, "Should return assets with 'server' in name")
}

// Test 12: TestAssetRepository_List_InvalidSortBy
func TestAssetRepository_List_InvalidSortBy(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	params := &model.AssetQueryParams{
		Limit:     10,
		Offset:    0,
		SortBy:    "DROP TABLE assets",
		SortOrder: "asc",
	}

	assets, err := repo.List(ctx, "test-user", params)

	assert.Nil(t, assets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_by")
}

// ========== Count Tests ==========

// Test 13: TestAssetRepository_Count_NoFilters
func TestAssetRepository_Count_NoFilters(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	for i := 1; i <= 7; i++ {
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name) VALUES ($1, $2)`, userID, fmt.Sprintf("Asset %d", i))
		require.NoError(t, err)
	}

	params := &model.AssetQueryParams{}
	count, err := repo.Count(ctx, userID, params)

	require.NoError(t, err)
	assert.Equal(t, int64(7), count)
}

// Test 14: TestAssetRepository_Count_TypeFilter
func TestAssetRepository_Count_TypeFilter(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	for i := 1; i <= 3; i++ {
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, fmt.Sprintf("Server %d", i), "server")
		require.NoError(t, err)
	}
	for i := 1; i <= 2; i++ {
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, fmt.Sprintf("NAS %d", i), "nas")
		require.NoError(t, err)
	}

	params := &model.AssetQueryParams{Type: testingPkg.Ptr("server")}
	count, err := repo.Count(ctx, userID, params)

	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// Test 15: TestAssetRepository_Count_MatchesList
func TestAssetRepository_Count_MatchesList(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	for i := 1; i <= 10; i++ {
		assetType := "server"
		if i%2 == 0 {
			assetType = "nas"
		}
		_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (user_id, name, type) VALUES ($1, $2, $3)`, userID, fmt.Sprintf("Asset %d", i), assetType)
		require.NoError(t, err)
	}

	params := &model.AssetQueryParams{
		Type:      testingPkg.Ptr("server"),
		Limit:     100,
		Offset:    0,
		SortBy:    "name",
		SortOrder: "asc",
	}

	assets, err := repo.List(ctx, userID, params)
	require.NoError(t, err)

	count, err := repo.Count(ctx, userID, params)
	require.NoError(t, err)

	assert.Equal(t, int64(len(assets)), count, "Count should match List results")
}

// ========== Create Tests ==========

// Test 16: TestAssetRepository_Create_AllFields
func TestAssetRepository_Create_AllFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	req := &model.CreateAssetRequest{
		Name:     "Production Server",
		Type:     testingPkg.Ptr("server"),
		Hostname: testingPkg.Ptr("prod.example.com"),
		Metadata: json.RawMessage(`{"cpu": "16 cores"}`),
	}

	asset, err := repo.Create(ctx, userID, req)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, asset.ID, "ID should be generated")
	assert.Equal(t, userID, asset.UserID)
	assert.Equal(t, "Production Server", asset.Name)
	assert.Equal(t, "server", *asset.Type)
	assert.Equal(t, "prod.example.com", *asset.Hostname)
	assert.NotNil(t, asset.Metadata)
	assert.NotZero(t, asset.CreatedAt)
	assert.NotZero(t, asset.UpdatedAt)
}

// Test 17: TestAssetRepository_Create_MinimalRequired
func TestAssetRepository_Create_MinimalRequired(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	req := &model.CreateAssetRequest{
		Name: "Minimal Asset",
	}

	asset, err := repo.Create(ctx, "test-user", req)

	require.NoError(t, err)
	assert.Equal(t, "Minimal Asset", asset.Name)
	assert.Nil(t, asset.Type)
	assert.Nil(t, asset.Hostname)
	assert.Nil(t, asset.Metadata)
}

// Test 18: TestAssetRepository_Create_GeneratesUniqueIDs
func TestAssetRepository_Create_GeneratesUniqueIDs(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	req := &model.CreateAssetRequest{Name: "Test"}

	asset1, err := repo.Create(ctx, "test-user", req)
	require.NoError(t, err)

	asset2, err := repo.Create(ctx, "test-user", req)
	require.NoError(t, err)

	assert.NotEqual(t, asset1.ID, asset2.ID, "IDs should be unique")
}

// ========== Update Tests ==========

// Test 19: TestAssetRepository_Update_SingleField_Name
func TestAssetRepository_Update_SingleField_Name(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	// Create asset
	userID := "test-user"
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name, type) VALUES ($1, $2, $3, $4)`, assetID, userID, "Old Name", "server")
	require.NoError(t, err)

	// Update only name
	updateReq := &model.UpdateAssetRequest{
		Name: testingPkg.Ptr("New Name"),
	}

	asset, err := repo.Update(ctx, userID, assetID, updateReq)

	require.NoError(t, err)
	assert.Equal(t, "New Name", asset.Name)
	assert.Equal(t, "server", *asset.Type, "Type should be unchanged")
}

// Test 20: TestAssetRepository_Update_AllFields
func TestAssetRepository_Update_AllFields(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "Old")
	require.NoError(t, err)

	updateReq := &model.UpdateAssetRequest{
		Name:     testingPkg.Ptr("Updated"),
		Type:     testingPkg.Ptr("vm"),
		Hostname: testingPkg.Ptr("new.host.com"),
		Metadata: json.RawMessage(`{"updated": true}`),
	}

	asset, err := repo.Update(ctx, userID, assetID, updateReq)

	require.NoError(t, err)
	assert.Equal(t, "Updated", asset.Name)
	assert.Equal(t, "vm", *asset.Type)
	assert.Equal(t, "new.host.com", *asset.Hostname)
	assert.NotNil(t, asset.Metadata)
}

// Test 21: TestAssetRepository_Update_AllNil
func TestAssetRepository_Update_AllNil(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name, type) VALUES ($1, $2, $3, $4)`, assetID, userID, "Original", "server")
	require.NoError(t, err)

	// Get original
	original, err := repo.GetByID(ctx, userID, assetID)
	require.NoError(t, err)

	// Update with all nil (should only update timestamp)
	updateReq := &model.UpdateAssetRequest{}
	asset, err := repo.Update(ctx, userID, assetID, updateReq)

	require.NoError(t, err)
	assert.Equal(t, "Original", asset.Name, "Name should be unchanged")
	assert.Equal(t, "server", *asset.Type, "Type should be unchanged")
	assert.True(t, asset.UpdatedAt.After(original.UpdatedAt) || asset.UpdatedAt.Equal(original.UpdatedAt), "UpdatedAt should change or stay same")
}

// Test 22: TestAssetRepository_Update_NotFound
func TestAssetRepository_Update_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	updateReq := &model.UpdateAssetRequest{Name: testingPkg.Ptr("New")}
	asset, err := repo.Update(ctx, "test-user", uuid.New(), updateReq)

	assert.Nil(t, asset)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
}

// Test 23: TestAssetRepository_Update_WrongUser
func TestAssetRepository_Update_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	updateReq := &model.UpdateAssetRequest{Name: testingPkg.Ptr("Hacked")}
	asset, err := repo.Update(ctx, "bob", assetID, updateReq)

	assert.Nil(t, asset)
	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
}

// ========== Delete Tests ==========

// Test 24: TestAssetRepository_Delete_Success
func TestAssetRepository_Delete_Success(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	userID := "test-user"
	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, userID, "To Delete")
	require.NoError(t, err)

	err = repo.Delete(ctx, userID, assetID)
	require.NoError(t, err)

	// Verify asset is gone
	asset, err := repo.GetByID(ctx, userID, assetID)
	assert.Nil(t, asset)
	assert.Error(t, err)
}

// Test 25: TestAssetRepository_Delete_NotFound
func TestAssetRepository_Delete_NotFound(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	err := repo.Delete(ctx, "test-user", uuid.New())

	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)
}

// Test 26: TestAssetRepository_Delete_WrongUser
func TestAssetRepository_Delete_WrongUser(t *testing.T) {
	testDB, _, cleanup := testingPkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewAssetRepository(testDB.Pool)

	assetID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `INSERT INTO assets (id, user_id, name) VALUES ($1, $2, $3)`, assetID, "alice", "Alice's Asset")
	require.NoError(t, err)

	err = repo.Delete(ctx, "bob", assetID)

	assert.Error(t, err)
	var httpErr *errs.HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 404, httpErr.Status)

	// Verify asset still exists for alice
	asset, err := repo.GetByID(ctx, "alice", assetID)
	require.NoError(t, err)
	assert.NotNil(t, asset)
}
