package btrnl

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	var nz sql.NullInt64
	err := nz.Scan(value)
	if err != nil {
		return err
	}
	n.Int64 = nz.Int64
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullInt64) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(bytes, &n.Int64)
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}
