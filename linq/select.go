package linq

import (
	"reflect"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
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
func (c *Linq) Find() (Items, error) {
	c.SqlSelect()

	c.sql = Format(`%s;`, c.sql)

	items, err := c.Query()
	if err != nil {
		return Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	if c.debug {
		console.Log(c.sql)
	}

	return items, nil
}

func (c *Linq) First() (Item, error) {
	c.sql = c.SqlLimit(1)

	item, err := c.QueryOne()
	if err != nil {
		return Item{}, err
	}

	c.Details(&item.Result)

	return item, nil
}

func (c *Linq) Limit(limit int) (Items, error) {
	c.sql = c.SqlLimit(limit)

	items, err := c.Query()
	if err != nil {
		return Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) Page(page, rows int) (Items, error) {
	offset := (page - 1) * rows
	c.sql = c.SqlOffset(rows, offset)

	items, err := c.Query()
	if err != nil {
		return Items{}, err
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

func (c *Linq) List(page, rows int) (List, error) {
	all := c.Count()

	items, err := c.Page(page, rows)
	if err != nil {
		return List{}, err
	}

	return items.ToList(all, page, rows), nil
}
