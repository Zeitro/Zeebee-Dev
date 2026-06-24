package timex

import (
	"testing"
	"time"
)

func TestStringToDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"seconds", "10s", 10 * time.Second, false},
		{"minutes", "5m", 5 * time.Minute, false},
		{"hours", "2h", 2 * time.Hour, false},
		{"composite", "1h30m45s", 1*time.Hour + 30*time.Minute + 45*time.Second, false},
		{"composite 2", "1h45s", 1*time.Hour + 45*time.Second, false},
		{"milliseconds", "500ms", 500 * time.Millisecond, false},
		{"with space", " 500ms", 500 * time.Millisecond, false},
		{"with space 2", " 500ms ", 500 * time.Millisecond, false},
		{"milliseconds trimmed", "500ms", 500 * time.Millisecond, false},
		{"zero", "0s", 0, false},
		{"empty", "", 0, true},
		{"invalid", "abc", 0, true},
		{"missing unit", "10", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToDuration(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("StringToDuration(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDurationToString(t *testing.T) {
	tests := []struct {
		name          string
		d             time.Duration
		humanReadable bool
		want          string
	}{
		{"seconds", 10 * time.Second, false, "10s"},
		{"minutes", 5 * time.Minute, false, "5m"},
		{"whole tens of minutes", 10 * time.Minute, false, "10m"},
		{"whole tens of minutes 2", 20 * time.Minute, false, "20m"},
		{"composite", 1*time.Hour + 30*time.Minute + 45*time.Second, false, "1h30m45s"},
		{"hour without minutes", 2 * time.Hour, false, "2h"},
		{"hour with zero minutes and seconds", 2*time.Hour + 30*time.Second, false, "2h30s"},
		{"zero", 0, false, "0s"},
		{"negative", -5 * time.Minute, false, "-5m"},
		{"milliseconds", 500 * time.Millisecond, false, "500ms"},

		{"readable seconds", 10 * time.Second, true, "10 seconds"},
		{"readable minutes", 5 * time.Minute, true, "5 minutes"},
		{"singular minute", 1 * time.Minute, true, "1 minute"},
		{"readable composite", 1*time.Hour + 30*time.Minute + 45*time.Second, true, "1 hour 30 minutes 45 seconds"},
		{"composite 2", 1*time.Hour + 45*time.Second, true, "1 hour 45 seconds"},
		{"hours only", 2 * time.Hour, true, "2 hours"},
		{"readable zero", 0, true, "0 second"},
		{"readable negative", -5 * time.Minute, true, "0 second"},
		{"readable milliseconds", 500 * time.Millisecond, true, "500 milliseconds"},
		{"singular millisecond", 1 * time.Millisecond, true, "1 millisecond"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationToString(tt.d, tt.humanReadable); got != tt.want {
				t.Errorf("DurationToString(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
