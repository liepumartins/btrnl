package btrnl

import (
	"database/sql/driver"
	"fmt"

	"github.com/robaho/fixed"
)

// NullFixed represents a nullable decimal with compatibility for
// scanning null values from the database.
type NullFixed struct {
	Fixed fixed.Fixed
	Valid bool
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullFixed) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	d.Valid = true
	switch v := value.(type) {
	case float32:
		d.Fixed = fixed.NewF(float64(v))
		return nil

	case float64:
		// numeric in sqlite3 sends us float64
		d.Fixed = fixed.NewF(v)
		return nil

	case int64:
		// at least in sqlite3 when the value is 0 in db, the data is sent
		// to us as an int64 instead of a float64 ...
		d.Fixed = fixed.NewI(v, 0)
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := unquoteIfQuoted(v)
		if err != nil {
			d.Valid = false
			return err
		}
		d.Fixed, err = fixed.NewSErr(str)
		if err != nil {
			d.Valid = false
		}
		return err
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (d NullFixed) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Fixed.String(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *NullFixed) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == "null" {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Fixed.UnmarshalJSON(decimalBytes)
}

// MarshalJSON implements the json.Marshaler interface.
func (d NullFixed) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return d.Fixed.MarshalJSON()
}

func unquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("Could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}
