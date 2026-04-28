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

	if c.tx != nil {
		if c.Tp == TpData {
			return c.tx.CommandSource(SourceField.Upp(), c.sql)
		}
		return c.tx.Command(c.sql)
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

/**
* queryWithCount executes a query that includes COUNT(*) OVER() AS _total
* and returns items plus the total count in a single round-trip.
* @return et.Items, int, error
**/
func (c *Linq) queryWithCount() (et.Items, int, error) {
	if c.debug {
		logs.Debug(c.sql)
	}

	if c.Tp == TpData {
		return c.db.SourceWithTotal("_total", SourceField.Upp(), c.sql)
	}

	return c.db.QueryWithTotal("_total", c.sql)
}
