package btrnl

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	var nz sql.NullFloat64
	err := nz.Scan(value)
	if err != nil {
		return err
	}
	n.Float64 = nz.Float64
	return nil
}

// Value implements the driver Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Float64, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullFloat64) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(bytes, &n.Float64)
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}
