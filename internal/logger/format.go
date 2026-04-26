package logger

import (
	"fmt"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	default:
		return fmt.Sprintf("%d", int(f))
	}
}

func (f *Format) Set(str string) error {
	switch str {
	case "text":
		*f = FormatText
	case "json":
		*f = FormatJSON
	default:
		return InvalidFormatString(str)
	}

	return nil
}

type InvalidFormatString string

func (e InvalidFormatString) Error() string {
	return fmt.Sprintf("logger: invalid format string '%s'", string(e))
}
