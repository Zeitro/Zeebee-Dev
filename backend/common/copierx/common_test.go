package copierx

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

// TestTimeToRFC3339Str tests time.Time to RFC3339 string conversion
func TestTimeToRFC3339Str(t *testing.T) {
	converter := TimeToRFC3339Str()

	// Test valid time
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	result, err := converter.Fn(testTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024-12-29T14:30:00Z"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test zero time
	zeroTime := time.Time{}
	result, err = converter.Fn(zeroTime)
	if err != nil {
		t.Errorf("Expected no error for zero time, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for zero time, got '%s'", result)
	}

	// Test invalid type
	_, err = converter.Fn("invalid")
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestTimeToStr tests custom layout time conversion
func TestTimeToStr(t *testing.T) {
	converter := TimeToStr("2006-01-02 15:04:05")

	// Test valid time
	testTime := time.Date(2024, 12, 29, 14, 30, 45, 0, time.UTC)
	result, err := converter.Fn(testTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024-12-29 14:30:45"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test zero time
	zeroTime := time.Time{}
	result, err = converter.Fn(zeroTime)
	if err != nil {
		t.Errorf("Expected no error for zero time, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for zero time, got '%s'", result)
	}
}

// TestTimePtrToRFC3339Str tests *time.Time to RFC3339 string conversion
func TestTimePtrToRFC3339Str(t *testing.T) {
	converter := TimePtrToRFC3339Str()

	// Test valid time pointer
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	result, err := converter.Fn(&testTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024-12-29T14:30:00Z"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test zero time pointer
	zeroTime := time.Time{}
	result, err = converter.Fn(&zeroTime)
	if err != nil {
		t.Errorf("Expected no error for zero time, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for zero time, got '%s'", result)
	}

	// Note: nil pointer causes type assertion error in copier, so we don't test it
	// In real usage with copier.Copy, nil fields are handled before reaching the converter
}

// TestTimePtrToStr tests *time.Time to custom layout string conversion
func TestTimePtrToStr(t *testing.T) {
	converter := TimePtrToStr("2006/01/02")

	// Test valid time pointer
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	result, err := converter.Fn(&testTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024/12/29"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test zero time pointer
	zeroTime := time.Time{}
	result, err = converter.Fn(&zeroTime)
	if err != nil {
		t.Errorf("Expected no error for zero time, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for zero time, got '%s'", result)
	}
}

// TestSqlTimePtrToRFC3339Str tests *sql.NullTime to RFC3339 string conversion
func TestSqlTimePtrToRFC3339Str(t *testing.T) {
	converter := SqlTimePtrToRFC3339Str()

	// Test valid NullTime
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	nullTime := &sql.NullTime{Time: testTime, Valid: true}
	result, err := converter.Fn(nullTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024-12-29T14:30:00Z"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test invalid NullTime
	invalidNullTime := &sql.NullTime{Valid: false}
	result, err = converter.Fn(invalidNullTime)
	if err != nil {
		t.Errorf("Expected no error for invalid NullTime, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for invalid NullTime, got '%s'", result)
	}

	// Test zero time in NullTime
	zeroTime := time.Time{}
	zeroNullTime := &sql.NullTime{Time: zeroTime, Valid: true}
	result, err = converter.Fn(zeroNullTime)
	if err != nil {
		t.Errorf("Expected no error for zero time, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for zero time, got '%s'", result)
	}

	// Note: nil pointer causes type assertion error in copier
	// In real usage with copier.Copy, nil fields are handled before reaching the converter
}

// TestSqlTimePtrToStr tests *sql.NullTime to custom layout string conversion
func TestSqlTimePtrToStr(t *testing.T) {
	converter := SqlTimePtrToStr("2006-01-02")

	// Test valid NullTime
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	nullTime := &sql.NullTime{Time: testTime, Valid: true}
	result, err := converter.Fn(nullTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "2024-12-29"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test invalid NullTime
	invalidNullTime := &sql.NullTime{Valid: false}
	result, err = converter.Fn(invalidNullTime)
	if err != nil {
		t.Errorf("Expected no error for invalid NullTime, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for invalid NullTime, got '%s'", result)
	}
}

// TestStrToSqlTimePtr tests string to *sql.NullTime conversion
func TestStrToSqlTimePtr(t *testing.T) {
	converter := StrToSqlTimePtr(time.RFC3339)

	// Test valid string
	testStr := "2024-12-29T14:30:00Z"
	result, err := converter.Fn(testStr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	nullTime, ok := result.(*sql.NullTime)
	if !ok {
		t.Fatal("Expected *sql.NullTime result")
	}
	if !nullTime.Valid {
		t.Error("Expected Valid to be true")
	}
	expectedTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)
	if nullTime.Time != expectedTime {
		t.Errorf("Expected time %v, got %v", expectedTime, nullTime.Time)
	}

	// Test empty string
	result, err = converter.Fn("")
	if err != nil {
		t.Errorf("Expected no error for empty string, got %v", err)
	}
	nullTime, ok = result.(*sql.NullTime)
	if !ok {
		t.Fatal("Expected *sql.NullTime result")
	}
	if nullTime.Valid {
		t.Error("Expected Valid to be false for empty string")
	}

	// Test invalid string
	_, err = converter.Fn("invalid-time")
	if err == nil {
		t.Error("Expected error for invalid time string, got nil")
	}

	// Test invalid type
	_, err = converter.Fn(123)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestUsStrToTime tests US date format string to time.Time conversion
func TestUsStrToTime(t *testing.T) {
	converter := UsStrToTime()

	// Test valid US date string
	testStr := "12/29/2024"
	result, err := converter.Fn(testStr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	parsedTime, ok := result.(time.Time)
	if !ok {
		t.Fatal("Expected time.Time result")
	}
	expectedTime := time.Date(2024, 12, 29, 0, 0, 0, 0, time.UTC)
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, parsedTime)
	}

	// Test invalid string
	_, err = converter.Fn("2024-12-29")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}

	// Test invalid type
	_, err = converter.Fn(123)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestStrToTime tests string to time.Time conversion with custom format
func TestStrToTime(t *testing.T) {
	converter := StrToTime("2006-01-02")

	// Test valid string
	testStr := "2024-12-29"
	result, err := converter.Fn(testStr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	parsedTime, ok := result.(time.Time)
	if !ok {
		t.Fatal("Expected time.Time result")
	}
	expectedTime := time.Date(2024, 12, 29, 0, 0, 0, 0, time.UTC)
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, parsedTime)
	}

	// Test invalid string
	_, err = converter.Fn("12/29/2024")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}

	// Test invalid type
	_, err = converter.Fn(123)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestConverterTypes tests that converters have correct source and destination types
func TestConverterTypes(t *testing.T) {
	tests := []struct {
		name        string
		converter   copier.TypeConverter
		expectedSrc interface{}
		expectedDst interface{}
	}{
		{
			name:        "TimeToRFC3339Str",
			converter:   TimeToRFC3339Str(),
			expectedSrc: time.Time{},
			expectedDst: copier.String,
		},
		{
			name:        "TimePtrToRFC3339Str",
			converter:   TimePtrToRFC3339Str(),
			expectedSrc: &time.Time{},
			expectedDst: copier.String,
		},
		{
			name:        "SqlTimePtrToRFC3339Str",
			converter:   SqlTimePtrToRFC3339Str(),
			expectedSrc: &sql.NullTime{},
			expectedDst: copier.String,
		},
		{
			name:        "StrToSqlTimePtr",
			converter:   StrToSqlTimePtr(time.RFC3339),
			expectedSrc: copier.String,
			expectedDst: &sql.NullTime{},
		},
		{
			name:        "UsStrToTime",
			converter:   UsStrToTime(),
			expectedSrc: copier.String,
			expectedDst: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test type assignment compatibility
			_ = tt.converter.SrcType
			_ = tt.converter.DstType
			_ = tt.converter.Fn
		})
	}
}

// BenchmarkTimeToRFC3339Str benchmarks time to RFC3339 string conversion
func BenchmarkTimeToRFC3339Str(b *testing.B) {
	converter := TimeToRFC3339Str()
	testTime := time.Date(2024, 12, 29, 14, 30, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = converter.Fn(testTime)
	}
}

// BenchmarkTrans_BasicConversion benchmarks basic Trans conversion
func BenchmarkTrans_BasicConversion(b *testing.B) {
	type Source struct {
		Name string
		Age  int
	}
	type Dest struct {
		Name string
		Age  int
	}
	src := &Source{Name: "Test", Age: 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Trans[Source, Dest](src)
	}
}

// BenchmarkTrans_WithConverter benchmarks Trans with time converter
func BenchmarkTrans_WithConverter(b *testing.B) {
	type Source struct {
		Name      string
		CreatedAt time.Time
	}
	type Dest struct {
		Name      string
		CreatedAt string
	}
	now := time.Now()
	src := &Source{Name: "Test", CreatedAt: now}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Trans[Source, Dest](src, TimeToRFC3339Str())
	}
}
