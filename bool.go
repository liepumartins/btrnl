package btrnl

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	var nz sql.NullBool
	err := nz.Scan(value)
	if err != nil {
		return err
	}
	n.Bool = nz.Bool
	return nil
}

// Value implements the driver Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullBool) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(bytes, &n.Bool)
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Bool)
}
