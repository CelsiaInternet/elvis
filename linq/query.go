package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
)

/**
* command
* @return et.Items, error
**/
func (c *Linq) command() (et.Items, error) {
	if c.debug {
		logs.Debug(c.sql)
	}

	if c.Tp == TpData {
		result, err := c.db.CommandSource(SourceField.Upp(), c.sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := c.db.Command(c.sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* query
* @return et.Items, error
**/
func (c *Linq) query() (et.Items, error) {
	if c.debug {
		logs.Debug(c.sql)
	}

	if c.Tp == TpData {
		result, err := c.db.Source(SourceField.Upp(), c.sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := c.db.Query(c.sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* queryCount
* @return int
**/
func (c *Linq) queryCount() int {
	if c.debug {
		logs.Debug(c.sql)
	}

	items, err := c.db.Query(c.sql)
	if err != nil {
		return 0
	}

	item := items.First()
	if !item.Ok {
		return 0
	}

	return item.Int("count")
}
