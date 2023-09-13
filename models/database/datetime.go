package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"manuel71sj/go-api-template/constants"
	"strings"
	"time"
)

// DateTime custom time types
// Used to format time into a human-readable string
type DateTime sql.NullTime

// Scan Scanner 인터페이스 구현
func (t *DateTime) Scan(value interface{}) error {
	return (*sql.NullTime)(t).Scan(value)
}

// Value driver Valuer 인터페이스 구현
func (t *DateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return t.Time.Format(constants.TimeFormat), nil
}

func (t *DateTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(constants.TimeFormat))), nil
	}

	return json.Marshal(nil)
}

func (t *DateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		t.Valid = false
		t.Time = time.Time{}

		return nil
	}

	cst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return fmt.Errorf("time.LoadLocation error: %s", err.Error())
	}

	t.Time, err = time.ParseInLocation(constants.TimeFormat, s, cst)
	if err != nil {
		// When time cannot be resolved using the default format, try RFC3339Nano
		if t.Time, err = time.ParseInLocation(time.RFC3339Nano, s, cst); err == nil {
			t.Time = t.Time.In(cst)
		}
	}

	t.Valid = true

	return err
}
