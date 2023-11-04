package linq

import (
	"strings"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

const TpSelect = 1
const TpData = 2

const ActSelect = 3
const ActInsert = 4
const ActDelete = 5
const ActUpdate = 6
const ActUpsert = 7

/**
*
**/
type SQL struct {
	linq *Linq
	val  string
}

/**
*
**/
type FRom struct {
	model *Model
	as    string
}

func (c *FRom) As() string {
	return c.as
}

func (c *FRom) NameAs() string {
	return Append(c.model.Name, c.as, " AS ")
}

func (c *FRom) Col(name string, cast ...string) *Col {
	_cast := ""
	if len(cast) > 0 {
		_cast = cast[0]
	}

	return &Col{
		from: c.as,
		name: name,
		cast: _cast,
	}
}

func (c *FRom) Column(name string, cast ...string) *Col {
	return c.Col(name, cast...)
}

/**
*
**/
type Join struct {
	kind  string
	from  *FRom
	join  *FRom
	where *Where
}

/**
*
**/
type OrderBy struct {
	colum  *Column
	sorted bool
}

/**
*
**/
type ReferenceValue struct {
	Schema string
	Table  string
	Key    any
	Op     int
}

type Rerences func(references []*ReferenceValue)

/**
*
**/
type Linq struct {
	Tp         int
	Act        int
	db         int
	_select    []*Column
	from       []*FRom
	where      []*Where
	_join      []*Join
	orderBy    []*OrderBy
	groupBy    []*Column
	_return    []*Column
	concat     string
	fromAs     []*FRom
	as         int
	details    []*Column
	data       Json
	new        *Json
	change     bool
	references []*ReferenceValue
	debug      bool
	sql        string
}

func GetAs(n int) string {
	limit := 18251
	base := 26
	as := ""
	a := n % base
	b := n / base
	c := b / base

	if n >= limit {
		n = n - limit + 702
		a = n % base
		b = n / base
		c = b / base
		b = b / base
		a = 65 + a
		b = 65 + b - 1
		c = 65 + c - 1
		as = Format(`A%c%c%c`, rune(c), rune(b), rune(a))
	} else if b > base {
		b = b / base
		a = 65 + a
		b = 65 + b - 1
		c = 65 + c - 1
		as = Format(`%c%c%c`, rune(c), rune(b), rune(a))
	} else if b > 0 {
		a = 65 + a
		b = 65 + b - 1
		as = Format(`%c%c`, rune(b), rune(a))
	} else {
		a = 65 + a
		as = Format(`%c`, rune(a))
	}

	return as
}

/**
*
**/
func NewLinq(tp int, act int, model *Model, as ...string) *Linq {
	if len(as) == 0 && act == ActSelect {
		as = []string{GetAs(0)}
	} else if len(as) == 0 {
		as = []string{""}
	}
	from := &FRom{model: model, as: Uppcase(as[0])}
	return &Linq{
		Tp:      tp,
		Act:     act,
		db:      model.Db,
		from:    []*FRom{from},
		fromAs:  []*FRom{from},
		where:   []*Where{},
		_join:   []*Join{},
		orderBy: []*OrderBy{},
		groupBy: []*Column{},
		details: []*Column{},
		as:      1,
	}
}

/**
*
**/
func (c *Linq) GetAs() string {
	result := GetAs(c.as)
	c.as++

	return result
}

func (c *Linq) Details(data *Json) *Json {
	for _, col := range c.details {
		col.Details(col, data)
	}

	return data
}

func (c *Linq) GetFrom(col *Column) *FRom {
	model := col.Model
	for _, item := range c.fromAs {
		if item.model.Up() == model.Up() {
			return item
		}
	}

	as := c.GetAs()
	result := &FRom{
		model: model,
		as:    as,
	}
	c.fromAs = append(c.fromAs, result)

	return result
}

func (c *Linq) SetFromAs(from *FRom) *FRom {
	model := from.model
	for _, item := range c.fromAs {
		if item.model.Up() == model.Up() {
			item.as = from.as
			return item
		}
	}

	as := c.GetAs()
	result := &FRom{
		model: model,
		as:    as,
	}
	c.fromAs = append(c.fromAs, result)

	return result
}

func (c *Linq) SetAs(model *Model, as string) string {
	for _, item := range c.fromAs {
		if item.model.Up() == model.Up() {
			item.as = as
			return as
		}
	}

	result := &FRom{
		model: model,
		as:    as,
	}
	c.fromAs = append(c.fromAs, result)

	return as
}

func (c *Linq) As(val any) string {
	switch v := val.(type) {
	case Column:
		col := &v
		return c.GetFrom(col).as
	case *Column:
		col := v
		return c.GetFrom(col).as
	case Col:
		col := &v
		return col.from
	case *Col:
		col := v
		return col.from
	default:
		return c.GetAs()
	}
}

func (c *Linq) GetCol(name string) *Column {
	if f := ReplaceAll(name, []string{" "}, ""); len(f) == 0 {
		return nil
	}

	pars := strings.Split(name, ".")
	if len(pars) == 2 {
		modelN := pars[0]
		colN := pars[0]
		for _, item := range c.fromAs {
			if item.model.Up() == Uppcase(modelN) {
				return item.model.Column(colN)
			}
		}
	} else if len(pars) == 1 {
		colN := pars[0]
		if len(c.fromAs) == 0 {
			return nil
		}

		return c.fromAs[0].model.Column(colN)
	}

	return nil
}

func (c *Linq) ToCols(sel ...any) []*Column {
	var cols []*Column
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
		}
	}

	return cols
}

/**
* Query
**/
func (c *Linq) Query() (Items, error) {	
	if c.Tp == TpData {
		result, err := DBQueryData(c.db, c.sql)
		if err != nil {
		return Items{}, err
	}

	if c.debug {
		console.Log(c.sql)
	}

	return result, nil
	} 
	
	result, err := DBQuery(c.db, c.sql)
	if err != nil {
		return Items{}, err
	}

	if c.debug {
		console.Log(c.sql)
	}

	return result, nil
}

func (c *Linq) QueryOne() (Item, error) {
	if c.Tp == TpData {
		result, err := DBQueryDataOne(c.db, c.sql)
		if err != nil {
			return Item{}, err
		}

		if c.debug {
			console.Log(c.sql)
		}

		return result, nil
	}

	result, err := DBQueryOne(c.db, c.sql)
	if err != nil {
		return Item{}, err
	}

	if c.debug {
		console.Log(c.sql)
	}

	return result, nil
}

func (c *Linq) QueryCount() int {
	result := DBQueryCount(c.db, c.sql)

	if c.debug {
		console.Log(c.sql)
	}

	return result
}
