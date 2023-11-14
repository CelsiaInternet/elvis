package linq

import (
	"reflect"

	"github.com/cgalvisleon/elvis/console"
	j "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
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
func (c *Linq) Find() (j.Items, error) {
	c.SqlSelect()

	c.sql = utility.Format(`%s;`, c.sql)

	items, err := c.Query()
	if err != nil {
		return j.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) All() (j.Items, error) {
	c.sql = c.SqlAll()

	items, err := c.Query()
	if err != nil {
		return j.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) First() (j.Item, error) {
	c.sql = c.SqlLimit(1)

	item, err := c.QueryOne()
	if err != nil {
		return j.Item{}, err
	}

	c.Details(&item.Result)

	return item, nil
}

func (c *Linq) Limit(limit int) (j.Items, error) {
	c.sql = c.SqlLimit(limit)

	items, err := c.Query()
	if err != nil {
		return j.Items{}, err
	}

	for _, data := range items.Result {
		c.Details(&data)
	}

	return items, nil
}

func (c *Linq) Page(page, rows int) (j.Items, error) {
	offset := (page - 1) * rows
	c.sql = c.SqlOffset(rows, offset)

	items, err := c.Query()
	if err != nil {
		return j.Items{}, err
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

func (c *Linq) List(page, rows int) (j.List, error) {
	all := c.Count()

	items, err := c.Page(page, rows)
	if err != nil {
		return j.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
