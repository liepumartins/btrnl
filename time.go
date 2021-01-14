package btrnl

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (nt *NullTime) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		nt.Valid = false
		return nil
	}
	nt.Valid = true
	return json.Unmarshal(bytes, &nt.Time)
}

// MarshalJSON implements the json.Marshaler interface.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}
