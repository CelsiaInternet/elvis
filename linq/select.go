package linq

import (
	"reflect"

	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/strs"
)

func (c *Linq) Select(sel ...any) *Linq {
	var cols []*Column = []*Column{}
	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			cols = append(cols, &v)
		case *Column:
			cols = append(cols, v)
		case string:
			c := c.GetCol(v)
			if c != nil {
				cols = append(cols, c)
			}
		case []string:
			for _, n := range v {
				c := c.GetCol(n)
				if c != nil {
					cols = append(cols, c)
				}
			}
		case []*Column:
			cols = v
		default:
			console.ErrorF("Linq select type (%v) value:%v", reflect.TypeOf(v), v)
		}
	}

	c._select = cols

	return c
}

/**
*
**/
func (c *Linq) Find() (e.Items, error) {
	c.SqlSelect()

	c.sql = strs.Format(`%s;`, c.sql)

	items, err := c.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) All() (e.Items, error) {
	c.sql = c.SqlAll()

	items, err := c.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) First() (e.Item, error) {
	c.sql = c.SqlLimit(1)

	item, err := c.QueryOne()
	if err != nil {
		return e.Item{}, err
	}

	c.Details(&item.Result)

	return item, nil
}

func (c *Linq) Limit(limit int) (e.Items, error) {
	c.sql = c.SqlLimit(limit)

	items, err := c.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) Page(page, rows int) (e.Items, error) {
	offset := (page - 1) * rows
	c.sql = c.SqlOffset(rows, offset)

	items, err := c.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) Count() int {
	c.sql = c.SqlCount()

	return c.QueryCount()
}

func (c *Linq) List(page, rows int) (e.List, error) {
	all := c.Count()

	items, err := c.Page(page, rows)
	if err != nil {
		return e.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
