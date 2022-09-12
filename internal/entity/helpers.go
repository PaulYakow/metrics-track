package entity

import (
	"database/sql/driver"
	"fmt"
)

type NullString string

func (s *NullString) Scan(value any) error {
	if value == nil {
		*s = ""
		return nil
	}

	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}

	*s = NullString(strVal)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	return string(s), nil
}
