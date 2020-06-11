package database

import "github.com/gocraft/dbr/v2"

// NewNullInt64 creates dbr.NullInt64 from int64 pointer
func NewNullInt64(v *int64) dbr.NullInt64 {
	res := dbr.NullInt64{}
	if v != nil {
		res.Int64 = *v
		res.Valid = true
	}

	return res
}

// NewNullString creates dbr.NullString from string pointer
func NewNullString(v *string) dbr.NullString {
	res := dbr.NullString{}
	if v != nil {
		res.String = *v
		res.Valid = true
	}

	return res
}

// NullStringToPtr creates pointer from dbr.NullString
func NullStringToPtr(v dbr.NullString) *string {
	if !v.Valid {
		return nil
	}

	return &v.String
}

// NullInt64ToPtr creates pointer from dbr.NullInt64
func NullInt64ToPtr(v dbr.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}

	return &v.Int64
}
