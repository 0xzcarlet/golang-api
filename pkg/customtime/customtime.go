package customtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date represents a date without time (YYYY-MM-DD format)
type Date struct {
	time.Time
}

// DateTime represents a datetime without timezone (YYYY-MM-DD HH:MM:SS format)
type DateTime struct {
	time.Time
}

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// Date methods

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Format(DateFormat))), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}

	// Remove quotes
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed, err := time.Parse(DateFormat, str)
	if err != nil {
		// Try parsing full datetime format and extract date
		parsed, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}
	d.Time = parsed
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
	return nil
}

// DateTime methods

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, dt.Format(DateTimeFormat))), nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}

	// Remove quotes
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Try multiple formats
	formats := []string{
		DateTimeFormat,
		time.RFC3339,
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
	}

	var parsed time.Time
	var err error
	for _, format := range formats {
		parsed, err = time.Parse(format, str)
		if err == nil {
			dt.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("cannot parse datetime: %s", str)
}

func (dt DateTime) Value() (driver.Value, error) {
	if dt.IsZero() {
		return nil, nil
	}
	return dt.Time, nil
}

func (dt *DateTime) Scan(value interface{}) error {
	if value == nil {
		dt.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		dt.Time = v
	default:
		return fmt.Errorf("cannot scan type %T into DateTime", value)
	}
	return nil
}

// Helper functions to create pointers

func NewDate(t time.Time) *Date {
	if t.IsZero() {
		return nil
	}
	return &Date{Time: t}
}

func NewDateTime(t time.Time) *DateTime {
	if t.IsZero() {
		return nil
	}
	return &DateTime{Time: t}
}
