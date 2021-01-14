package btrnl

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var s NullString
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  if s.Valid {
//     // use s.String
//  } else {
//     // NULL value
//  }
//
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	var nz sql.NullString
	err := nz.Scan(value)
	if err != nil {
		return err
	}
	ns.String = nz.String
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ns *NullString) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(bytes, &ns.String)
}

// MarshalJSON implements the json.Marshaler interface.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
