package stringx

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple string",
			format:   "Hello, World!",
			args:     nil,
			expected: "Hello, World!",
		},
		{
			name:     "string with placeholder",
			format:   "Hello, %s!",
			args:     []interface{}{"Go"},
			expected: "Hello, Go!",
		},
		{
			name:     "multiple placeholders",
			format:   "%s %d %v",
			args:     []interface{}{"test", 42, true},
			expected: "test 42 true",
		},
		{
			name:     "integer placeholder",
			format:   "Number: %d",
			args:     []interface{}{123},
			expected: "Number: 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args...)
			if result != tt.expected {
				t.Errorf("Format() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseFullName(t *testing.T) {
	tests := []struct {
		name      string
		fullName  string
		firstName string
		midName   string
		lastName  string
	}{
		{
			name:      "empty string",
			fullName:  "",
			firstName: "",
			midName:   "",
			lastName:  "",
		},
		{
			name:      "only first name",
			fullName:  "John",
			firstName: "John",
			midName:   "",
			lastName:  "",
		},
		{
			name:      "first and last name",
			fullName:  "John Doe",
			firstName: "John",
			midName:   "",
			lastName:  "Doe",
		},
		{
			name:      "first, middle, and last name",
			fullName:  "John Michael Doe",
			firstName: "John",
			midName:   "Michael",
			lastName:  "Doe",
		},
		{
			name:      "multiple middle names",
			fullName:  "John Michael William Doe",
			firstName: "John",
			midName:   "Michael William",
			lastName:  "Doe",
		},
		{
			name:      "with extra spaces",
			fullName:  "  John Doe  ",
			firstName: "John",
			midName:   "",
			lastName:  "Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstName, midName, lastName := ParseFullName(tt.fullName)
			if firstName != tt.firstName || midName != tt.midName || lastName != tt.lastName {
				t.Errorf("ParseFullName(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.fullName, firstName, midName, lastName,
					tt.firstName, tt.midName, tt.lastName)
			}
		})
	}
}

func TestFormatFullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		midName   string
		lastName  string
		expected  string
	}{
		{
			name:      "all empty",
			firstName: "",
			midName:   "",
			lastName:  "",
			expected:  "",
		},
		{
			name:      "only first name",
			firstName: "John",
			midName:   "",
			lastName:  "",
			expected:  "John",
		},
		{
			name:      "first and last name",
			firstName: "John",
			midName:   "",
			lastName:  "Doe",
			expected:  "John Doe",
		},
		{
			name:      "first, middle, and last name",
			firstName: "John",
			midName:   "Michael",
			lastName:  "Doe",
			expected:  "John Michael Doe",
		},
		{
			name:      "only middle name",
			firstName: "",
			midName:   "Michael",
			lastName:  "",
			expected:  "Michael",
		},
		{
			name:      "first and middle name",
			firstName: "John",
			midName:   "Michael",
			lastName:  "",
			expected:  "John Michael",
		},
		{
			name:      "middle and last name",
			firstName: "",
			midName:   "Michael",
			lastName:  "Doe",
			expected:  "Michael Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFullName(tt.firstName, tt.midName, tt.lastName)
			if result != tt.expected {
				t.Errorf("FormatFullName(%q, %q, %q) = %q, want %q",
					tt.firstName, tt.midName, tt.lastName, result, tt.expected)
			}
		})
	}
}

func TestTransNumberToString(t *testing.T) {
	tests := []struct {
		name        string
		num         interface{}
		zeroToEmpty bool
		expected    string
	}{
		{
			name:        "int with zero and zeroToEmpty=true",
			num:         int(0),
			zeroToEmpty: true,
			expected:    "",
		},
		{
			name:        "int with zero and zeroToEmpty=false",
			num:         int(0),
			zeroToEmpty: false,
			expected:    "0",
		},
		{
			name:        "int with positive value",
			num:         int(42),
			zeroToEmpty: true,
			expected:    "42",
		},
		{
			name:        "int64 with value",
			num:         int64(1234567890),
			zeroToEmpty: true,
			expected:    "1234567890",
		},
		{
			name:        "float64 with zero and zeroToEmpty=true",
			num:         float64(0.0),
			zeroToEmpty: true,
			expected:    "",
		},
		{
			name:        "float64 with value",
			num:         float64(123.456),
			zeroToEmpty: true,
			expected:    "123.456",
		},
		{
			name:        "uint with value",
			num:         uint(100),
			zeroToEmpty: true,
			expected:    "100",
		},
		{
			name:        "negative int",
			num:         int(-50),
			zeroToEmpty: true,
			expected:    "-50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.num.(type) {
			case int:
				result := TransNumberToString(v, tt.zeroToEmpty)
				if result != tt.expected {
					t.Errorf("TransNumberToString(%d, %v) = %q, want %q", v, tt.zeroToEmpty, result, tt.expected)
				}
			case int64:
				result := TransNumberToString(v, tt.zeroToEmpty)
				if result != tt.expected {
					t.Errorf("TransNumberToString(%d, %v) = %q, want %q", v, tt.zeroToEmpty, result, tt.expected)
				}
			case float64:
				result := TransNumberToString(v, tt.zeroToEmpty)
				if result != tt.expected {
					t.Errorf("TransNumberToString(%f, %v) = %q, want %q", v, tt.zeroToEmpty, result, tt.expected)
				}
			case uint:
				result := TransNumberToString(v, tt.zeroToEmpty)
				if result != tt.expected {
					t.Errorf("TransNumberToString(%d, %v) = %q, want %q", v, tt.zeroToEmpty, result, tt.expected)
				}
			}
		})
	}
}

func TestFormatUsPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid 10-digit number",
			input:    "1234567890",
			expected: "(123) 456-7890",
		},
		{
			name:     "valid number with spaces",
			input:    " 1234567890 ",
			expected: "(123) 456-7890",
		},
		{
			name:     "less than 10 digits",
			input:    "123456",
			expected: "123456",
		},
		{
			name:     "more than 10 digits",
			input:    "1234567890123",
			expected: "1234567890123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "non-numeric characters",
			input:    "abcdefghij",
			expected: "(abc) def-ghij",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUsPhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("FormatUsPhoneNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatUsSsn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid 9-digit SSN",
			input:    "123456789",
			expected: "123-45-6789",
		},
		{
			name:     "valid SSN with spaces",
			input:    " 123456789 ",
			expected: "123-45-6789",
		},
		{
			name:     "less than 9 digits",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "more than 9 digits",
			input:    "1234567890123",
			expected: "1234567890123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "non-numeric characters",
			input:    "abcdefghi",
			expected: "abc-de-fghi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUsSsn(tt.input)
			if result != tt.expected {
				t.Errorf("FormatUsSsn(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		amount   float64
		expected string
	}{
		{
			name:     "zero amount",
			symbol:   "$",
			amount:   0,
			expected: "",
		},
		{
			name:     "positive amount with dollar",
			symbol:   "$",
			amount:   123.456,
			expected: "$123.46",
		},
		{
			name:     "positive amount with euro",
			symbol:   "€",
			amount:   1000.00,
			expected: "€1000.00",
		},
		{
			name:     "small amount",
			symbol:   "$",
			amount:   0.01,
			expected: "$0.01",
		},
		{
			name:     "large amount",
			symbol:   "$",
			amount:   999999.99,
			expected: "$999999.99",
		},
		{
			name:     "negative amount",
			symbol:   "$",
			amount:   -50.25,
			expected: "$-50.25",
		},
		{
			name:     "empty symbol",
			symbol:   "",
			amount:   100.00,
			expected: "100.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCurrency(tt.symbol, tt.amount)
			if result != tt.expected {
				t.Errorf("FormatCurrency(%q, %f) = %q, want %q", tt.symbol, tt.amount, result, tt.expected)
			}
		})
	}
}

func TestFormatUsDollarCurrency(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{
			name:     "zero amount",
			amount:   0,
			expected: "",
		},
		{
			name:     "positive amount",
			amount:   123.45,
			expected: "$123.45",
		},
		{
			name:     "round number",
			amount:   100.00,
			expected: "$100",
		},
		{
			name:     "small amount",
			amount:   0.01,
			expected: "$0.01",
		},
		{
			name:     "large amount",
			amount:   999999.99,
			expected: "$999,999.99",
		},
		{
			name:     "negative amount",
			amount:   -50.25,
			expected: "-$50.25",
		},
		{
			name:     "amount with cents",
			amount:   1234.56,
			expected: "$1,234.56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUsDollarCurrency(tt.amount)
			if result != tt.expected {
				t.Errorf("FormatUsDollarCurrency(%f) = %q, want %q", tt.amount, result, tt.expected)
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "all uppercase",
			input:    "WORLD",
			expected: "World",
		},
		{
			name:     "sentence",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "mixed case",
			input:    "hElLo WoRLd",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "multiple words",
			input:    "the quick brown fox",
			expected: "The Quick Brown Fox",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToTitle(tt.input)
			if result != tt.expected {
				t.Errorf("ToTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkParseFullName(b *testing.B) {
	fullName := "John Michael William Doe"
	for i := 0; i < b.N; i++ {
		ParseFullName(fullName)
	}
}

func BenchmarkFormatFullName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatFullName("John", "Michael", "Doe")
	}
}

func BenchmarkFormatUsPhoneNumber(b *testing.B) {
	phone := "1234567890"
	for i := 0; i < b.N; i++ {
		FormatUsPhoneNumber(phone)
	}
}

func BenchmarkToTitle(b *testing.B) {
	str := "hello world"
	for i := 0; i < b.N; i++ {
		ToTitle(str)
	}
}
