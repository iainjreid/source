package logger_test

import (
	"testing"

	"github.com/iainjreid/source/internal/logger"
)

var formattests = []struct {
	in string
}{
	{"text"},
	{"json"},
}

func TestFormatSetAndGet(t *testing.T) {
	var format logger.Format

	for _, tt := range formattests {
		t.Run(tt.in, func(t *testing.T) {
			if err := format.Set(tt.in); err != nil {
				t.Errorf("Format.Set(%v) returned an unexpected error = %v", tt.in, err)
			}

			if got := format.String(); got != tt.in {
				t.Errorf("Format(%s).String() = %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

var formatfailures = []struct {
	in string
}{
	{"invalid"},
	{""},
}

func TestInvalidFormat(t *testing.T) {
	var format logger.Format

	for _, tt := range formatfailures {
		t.Run(tt.in, func(t *testing.T) {
			if got, want := format.Set(tt.in), logger.InvalidFormatString(tt.in); got != want {
				t.Errorf("Format.Set(%v) = %v, want %s", tt.in, got, want)
			}
		})
	}
}
