package custom

import (
	"database/sql/driver"
	"encoding/json"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}
