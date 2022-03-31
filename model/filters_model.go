package model

import (
	"database/sql/driver"
	"fmt"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"strings"
	"time"
)

var logger logging.Logger

const Layout = "20060102"

type MyTime struct {
	time.Time
}

func (c *MyTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" || string(data) == "" {
		logger.Error("date  is not specified")
		return fmt.Errorf("date is not specified")
	} else {
		s := strings.Trim(string(data), "\"")
		// Fractional seconds are handled implicitly by Parse.
		tt, err := time.Parse(Layout, s)
		*c = MyTime{tt}
		return err
	}
}

func (c MyTime) Value() (driver.Value, error) {
	return driver.Value(c.Time), nil
}

func (c *MyTime) Scan(src interface{}) error {
	switch t := src.(type) {
	case time.Time:
		c.Time = t
		return nil
	default:
		return fmt.Errorf("column type not supported")
	}
}
func (c MyTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(Layout))), nil
}

type RequestFilters struct {
	ShowDeleted bool   `form:"show_deleted,omitempty"`
	FilterData  bool   `form:"filter_data,omitempty"`
	StartTime   MyTime `form:"start_time,omitempty"`
	EndTime     MyTime `form:"end_time,omitempty"`
	Role        string `form:"role,omitempty"`
}
