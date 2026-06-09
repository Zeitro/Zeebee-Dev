package copierx

import (
	"database/sql"
	"testing"
	"time"
)

// Test models
type sourceModel struct {
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *sql.NullTime
}

type destModel struct {
	Name      string
	Age       int
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type nestedSource struct {
	FirstName string
	LastName  string
}

type nestedDest struct {
	FirstName string
	LastName  string
}

// TestTrans_NilSource tests Trans with nil source
func TestTrans_NilSource(t *testing.T) {
	result, err := Trans[nestedSource, nestedDest](nil)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

// TestTrans_BasicConversion tests basic struct conversion
func TestTrans_BasicConversion(t *testing.T) {
	now := time.Now()
	src := &sourceModel{
		Name:      "John Doe",
		Age:       30,
		CreatedAt: now,
	}

	result, err := Trans[sourceModel, destModel](src)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", result.Name)
	}
	if result.Age != 30 {
		t.Errorf("Expected Age 30, got %d", result.Age)
	}
}

// TestTrans_WithTimeConverter tests Trans with time to string converter
func TestTrans_WithTimeConverter(t *testing.T) {
	now := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	src := &sourceModel{
		Name:      "Jane Doe",
		Age:       25,
		CreatedAt: now,
	}

	result, err := Trans[sourceModel, destModel](src, TimeToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	expectedTime := "2024-12-29T14:30:00Z"
	if result.CreatedAt != expectedTime {
		t.Errorf("Expected CreatedAt '%s', got '%s'", expectedTime, result.CreatedAt)
	}
}

// TestTrans_WithTimePtrConverter tests Trans with time pointer to string converter
func TestTrans_WithTimePtrConverter(t *testing.T) {
	src := &sourceModel{
		Name:      "Bob Smith",
		Age:       35,
		UpdatedAt: new(time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)),
	}

	result, err := Trans[sourceModel, destModel](src, TimePtrToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	expectedTime := "2024-12-29T14:30:00Z"
	if result.UpdatedAt != expectedTime {
		t.Errorf("Expected UpdatedAt '%s', got '%s'", expectedTime, result.UpdatedAt)
	}
}

// TestTrans_WithNilTimePtr tests Trans with nil time pointer
func TestTrans_WithNilTimePtr(t *testing.T) {
	src := &sourceModel{
		Name:      "Alice",
		Age:       28,
		UpdatedAt: nil,
	}

	result, err := Trans[sourceModel, destModel](src, TimePtrToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.UpdatedAt != "" {
		t.Errorf("Expected empty UpdatedAt for nil time ptr, got '%s'", result.UpdatedAt)
	}
}

// TestTrans_WithSqlNullTimeConverter tests Trans with sql.NullTime to string converter
func TestTrans_WithSqlNullTimeConverter(t *testing.T) {
	now := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	nullTime := &sql.NullTime{Time: now, Valid: true}
	src := &sourceModel{
		Name:      "Charlie",
		Age:       40,
		DeletedAt: nullTime,
	}

	result, err := Trans[sourceModel, destModel](src, SqlTimePtrToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	expectedTime := "2024-12-29T14:30:00Z"
	if result.DeletedAt != expectedTime {
		t.Errorf("Expected DeletedAt '%s', got '%s'", expectedTime, result.DeletedAt)
	}
}

// TestTrans_WithInvalidSqlNullTime tests Trans with invalid sql.NullTime
func TestTrans_WithInvalidSqlNullTime(t *testing.T) {
	nullTime := &sql.NullTime{Valid: false}
	src := &sourceModel{
		Name:      "David",
		Age:       45,
		DeletedAt: nullTime,
	}

	result, err := Trans[sourceModel, destModel](src, SqlTimePtrToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.DeletedAt != "" {
		t.Errorf("Expected empty DeletedAt for invalid NullTime, got '%s'", result.DeletedAt)
	}
}

// TestTrans_WithMultipleConverters tests Trans with multiple type converters
func TestTrans_WithMultipleConverters(t *testing.T) {
	now := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	nullTime := &sql.NullTime{Time: now, Valid: true}
	src := &sourceModel{
		Name:      "Eve",
		Age:       32,
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: nullTime,
	}

	result, err := Trans[sourceModel, destModel](src,
		TimeToRFC3339Str(),
		TimePtrToRFC3339Str(),
		SqlTimePtrToRFC3339Str(),
	)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	expectedTime := "2024-12-29T14:30:00Z"
	if result.CreatedAt != expectedTime {
		t.Errorf("Expected CreatedAt '%s', got '%s'", expectedTime, result.CreatedAt)
	}
	if result.UpdatedAt != expectedTime {
		t.Errorf("Expected UpdatedAt '%s', got '%s'", expectedTime, result.UpdatedAt)
	}
	if result.DeletedAt != expectedTime {
		t.Errorf("Expected DeletedAt '%s', got '%s'", expectedTime, result.DeletedAt)
	}
}

// TestTrans_ZeroTime tests Trans with zero time value
func TestTrans_ZeroTime(t *testing.T) {
	zeroTime := time.Time{}
	src := &sourceModel{
		Name:      "Frank",
		Age:       50,
		CreatedAt: zeroTime,
	}

	result, err := Trans[sourceModel, destModel](src, TimeToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.CreatedAt != "" {
		t.Errorf("Expected empty CreatedAt for zero time, got '%s'", result.CreatedAt)
	}
}

// TestTrans_SliceConversion tests slice conversion
func TestTrans_SliceConversion(t *testing.T) {
	now := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	srcList := []*sourceModel{
		{Name: "Alice", Age: 25, CreatedAt: now},
		{Name: "Bob", Age: 30, CreatedAt: now},
		{Name: "Charlie", Age: 35, CreatedAt: now},
	}

	result, err := Trans[[]*sourceModel, []*destModel](&srcList, TimeToRFC3339Str())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(*result) != 3 {
		t.Errorf("Expected 3 items, got %d", len(*result))
	}
	if (*result)[0].Name != "Alice" {
		t.Errorf("Expected first name 'Alice', got '%s'", (*result)[0].Name)
	}
	if (*result)[1].Name != "Bob" {
		t.Errorf("Expected second name 'Bob', got '%s'", (*result)[1].Name)
	}
	if (*result)[2].Name != "Charlie" {
		t.Errorf("Expected third name 'Charlie', got '%s'", (*result)[2].Name)
	}
}

// TestTrans_EmptySlice tests empty slice conversion
func TestTrans_EmptySlice(t *testing.T) {
	srcList := []*sourceModel{}

	result, err := Trans[[]*sourceModel, []*destModel](&srcList)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(*result) != 0 {
		t.Errorf("Expected 0 items, got %d", len(*result))
	}
}

// TestTrans_NestedSliceConversion tests nested slice conversion with nil elements
func TestTrans_NestedSliceConversion(t *testing.T) {
	srcList := []*sourceModel{
		{Name: "Alice", Age: 25},
		nil,
		{Name: "Bob", Age: 30},
	}

	result, err := Trans[[]*sourceModel, []*destModel](&srcList)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	// copier should handle nil elements in slice
	if len(*result) != 3 {
		t.Errorf("Expected 3 items, got %d", len(*result))
	}
}
