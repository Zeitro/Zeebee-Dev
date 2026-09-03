package swagger

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testApiContent = `syntax = "v1"

info (
	title:   "test-api"
	desc:    "test api service"
	version: "v1"
)

@server (
	desc:       "user service"
	version:    v1
	group:      user
	prefix:     /api/v1
	middleware: AuthMiddleware
	tags:       user
)
service test-api {
	@doc "user login"
	@handler UserLogin
	post /user/login (EmptyRequest) returns (EmptyResponse)

	@doc "user info"
	@handler UserInfo
	get /user/info (EmptyRequest) returns (EmptyResponse)
}

@server (
	desc:    "loan service"
	version: v1
	group:   loan
	prefix:  /api/v1
	tags:    loan
)
service test-api {
	@doc "loan list"
	@handler LoanList
	get /loan/list (EmptyRequest) returns (EmptyResponse)
}

@server (
	desc:    "loan service"
	version: v1
	group:   loan
	prefix:  /api/v1
	tags:    loan
	middleware: AuthMiddleware
)
service test-api {
	@doc "delete"
	@handler LoanDelete
	delete /loan/:id (EmptyRequest) returns (EmptyResponse)
}

@server (
	desc:       "borrower property service"
	version:    v1
	group:      borrower_property
	prefix:     /api/v1
	middleware: AuthMiddleware
	tags:       property
)
service test-api {
	@doc "page"
	@handler PropertyPage
	post /property/page (PropertyPageRequest) returns (PropertyPageResponse)

	@doc "all"
	@handler PropertyAll
	get /property/all (EmptyRequest) returns (PropertyAllResponse)

	@doc "property Summary"
	@handler PropertyDetailSummary
	get /property/:id/summary (PropertyDetailRequest) returns (PropertyDetailSummaryResponse)

	@doc "detail"
	@handler PropertyDetail
	get /property/:id (PathIdRequest) returns (PropertyDetailResponse)

	@doc "create"
	@handler PropertyCreate
	post /property (PropertyCreateRequest) returns (PropertyCreateResponse)

	@doc "update"
	@handler PropertyUpdate
	put /property/:id (PropertyUpdateRequest) returns (EmptyResponse)

	@doc "delete"
	@handler PropertyDelete
	delete /property/:id (PathIdRequest) returns (EmptyResponse)
}
`

func TestParseApi_GroupsByTag(t *testing.T) {
	result, err := ParseApi(testApiContent)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result["user"], 2)
	assert.Len(t, result["loan"], 2)
	assert.Len(t, result["property"], 7)
}

func TestParseApi_ExtractsFields(t *testing.T) {
	result, err := ParseApi(testApiContent)
	require.NoError(t, err)

	login := result["user"][0]
	assert.Equal(t, "/api/v1/user/login", login.Path)
	assert.Equal(t, "POST", login.Method)
	assert.Equal(t, "user login", login.Name)
	assert.Equal(t, "user login", login.Description)
}

func TestParseApi_PropertyRoutes(t *testing.T) {
	result, err := ParseApi(testApiContent)
	require.NoError(t, err)

	property := result["property"]
	assert.Equal(t, "/api/v1/property/:id/summary", property[2].Path)
	assert.Equal(t, "property Summary", property[2].Name)
	assert.Equal(t, "property Summary", property[2].Description)

	assert.Equal(t, "/api/v1/property/:id", property[3].Path)
	assert.Equal(t, "DELETE", property[6].Method)
}

func TestParseApi_EmptyContent(t *testing.T) {
	result, err := ParseApi(`syntax = "v1"
info (
	title: "empty"
	version: "v1"
)
`)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseApi_MinimalContent(t *testing.T) {
	result, err := ParseApi("")
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseApi_FullExpectedData(t *testing.T) {
	result, err := ParseApi(testApiContent)
	require.NoError(t, err)

	expected := ApiGroup{
		"user": {
			{Path: "/api/v1/user/login", Method: "POST", Name: "user login", Description: "user login"},
			{Path: "/api/v1/user/info", Method: "GET", Name: "user info", Description: "user info"},
		},
		"loan": {
			{Path: "/api/v1/loan/list", Method: "GET", Name: "loan list", Description: "loan list"},
			{Path: "/api/v1/loan/:id", Method: "DELETE", Name: "delete", Description: "delete"},
		},
		"property": {
			{Path: "/api/v1/property/page", Method: "POST", Name: "page", Description: "page"},
			{Path: "/api/v1/property/all", Method: "GET", Name: "all", Description: "all"},
			{Path: "/api/v1/property/:id/summary", Method: "GET", Name: "property Summary", Description: "property Summary"},
			{Path: "/api/v1/property/:id", Method: "GET", Name: "detail", Description: "detail"},
			{Path: "/api/v1/property", Method: "POST", Name: "create", Description: "create"},
			{Path: "/api/v1/property/:id", Method: "PUT", Name: "update", Description: "update"},
			{Path: "/api/v1/property/:id", Method: "DELETE", Name: "delete", Description: "delete"},
		},
	}

	assert.Equal(t, expected, result)
}

func TestParseSwagger(t *testing.T) {
	result, err := ParseSwagger([]byte(`{
		"paths": {
			"/api/v1/project_management/project/{id}": {
				"get": {
					"tags": ["project_management_project"],
					"summary": "detail"
				}
			}
		}
	}`))
	require.NoError(t, err)
	assert.Equal(t, ApiGroup{
		"project_management_project": {
			{Path: "/api/v1/project_management/project/:id", Method: "GET", Name: "detail", Description: "detail"},
		},
	}, result)
}

func TestParseSwagger_BorrowerManifestIncludesImportedServers(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	content, err := os.ReadFile(filepath.Join(filepath.Dir(file), "../../app/borrower/api/swagger/borrower.json"))
	require.NoError(t, err)

	result, err := ParseSwagger(content)
	require.NoError(t, err)
	assert.Contains(t, result["project_management_project"], Api{
		Path:        "/api/v1/project_management/project/:id",
		Method:      "GET",
		Name:        "detail",
		Description: "detail",
	})
}
