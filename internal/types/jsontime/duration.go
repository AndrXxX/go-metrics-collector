package jsontime

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		dVal, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration{dVal}
		return nil
	}
	return fmt.Errorf("cannot unmarshal JSON value")
}
