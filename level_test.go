package slfg

import "testing"

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{TRACE, "TRACE"},
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{PANIC, "PANIC"},
		{FATAL, "FATAL"},
		{Level(99), "LEVEL(99)"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", int(tt.level), got, tt.want)
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    Level
		wantErr bool
	}{
		{"TRACE", TRACE, false},
		{"debug", DEBUG, false},
		{"Info", INFO, false},
		{"WARN", WARN, false},
		{"error", ERROR, false},
		{"PANIC", PANIC, false},
		{"fatal", FATAL, false},
		{"unknown", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		got, err := ParseLevel(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
