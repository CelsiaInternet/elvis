package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utilities"
)

const TpColumn = 0
const TpAtrib = 1
const TpDetail = 2
const TpReference = 3
const TpFunction = 4
const TpClone = 5
const TpField = 6

/**
*
**/
type Col struct {
	from string
	name string
	cast string
	as   string
}

func (c *Col) Up() string {
	return utilities.Uppcase(c.name)
}

func (c *Col) Low() string {
	return utilities.Lowcase(c.name)
}

func (c *Col) Cast(cast string) *Col {
	c.cast = cast

	return c
}

func (c *Col) As() string {
	if len(c.as) == 0 {
		return c.name
	}

	return c.as
}

func (c *Col) AsUp() string {
	return utilities.Uppcase(c.As())
}

func (c *Col) AsLow() string {
	return utilities.Lowcase(c.As())
}

/**
*
**/
type Details func(col *Column, data *json.Json)

type Column struct {
	Model       *Model
	Tp          int
	Column      *Column
	name        string
	Title       string
	Description string
	Type        string
	Default     any
	Atribs      []*Column
	Reference   *Reference
	Definition  interface{}
	Function    string
	Details     Details
	Indexed     bool
	Unique      bool
	Required    bool
	PrimaryKey  bool
	ForeignKey  bool
	Hidden      bool
	from        string
	cast        string
}

func (c *Column) describe() json.Json {
	return json.Json{
		"name":        c.name,
		"description": c.Description,
		"type":        c.Type,
		"default":     c.Default,
		"tp":          c.Tp,
		"indexed":     c.Indexed,
		"unique":      c.Unique,
		"required":    c.Required,
		"primaryKey":  c.PrimaryKey,
		"foreignKey":  c.ForeignKey,
		"hidden":      c.Hidden,
	}
}

func (c *Column) Describe() json.Json {
	var atribs []json.Json = []json.Json{}
	for _, atrib := range c.Atribs {
		atribs = append(atribs, atrib.describe())
	}

	reference := json.Json{}
	if c.Reference != nil {
		reference = c.Reference.Describe()
	}

	return json.Json{
		"name":        c.name,
		"description": c.Description,
		"type":        c.Type,
		"default":     c.Default,
		"atribs":      atribs,
		"reference":   reference,
		"tp":          c.Tp,
		"indexed":     c.Indexed,
		"unique":      c.Unique,
		"required":    c.Required,
		"hidden":      c.Hidden,
	}
}

func NewColumn(model *Model, name, description, _type string, _default any) *Column {
	result := &Column{
		Model:       model,
		Tp:          TpColumn,
		name:        utilities.Uppcase(name),
		Description: description,
		Type:        _type,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
	}

	if !model.UseDateMake {
		model.UseDateMake = utilities.Uppcase(result.name) == utilities.Uppcase(model.DateMakeField)
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = utilities.Uppcase(result.name) == utilities.Uppcase(model.DateUpdateField)
	}

	if !model.UseState {
		model.UseState = utilities.Uppcase(result.name) == utilities.Uppcase(model.StateField)
	}

	if !model.UseRecycle {
		model.UseRecycle = utilities.Uppcase(result.name) == utilities.Uppcase(model.StateField)
	}

	if !model.UseProject {
		model.UseProject = utilities.Uppcase(result.name) == utilities.Uppcase(model.ProjectField)
	}

	if !model.UseIndex {
		model.UseIndex = utilities.Uppcase(result.name) == utilities.Uppcase(model.IndexField)
	}

	if !model.UseSource {
		model.UseSource = utilities.Uppcase(result.name) == utilities.Uppcase(model.SourceField)
	}

	model.Definition = append(model.Definition, result)
	return result
}

func NewVirtualAtrib(model *Model, name, description, _type string, _default any) *Column {
	result := &Column{
		Model:       model,
		Tp:          TpColumn,
		name:        utilities.Uppcase(name),
		Description: description,
		Type:        _type,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
	}

	return result
}

/**
* DDL
**/
func (c *Column) DDL() string {
	var result string

	if c.Model.integrityReference && c.ForeignKey {
		result = c.Reference.DDL()
	}

	_default := utilities.NewAny(c.Default)

	if _default.String() == "NOW()" {
		result = utilities.Append(`DEFAULT NOW()`, result, " ")
	} else {
		result = utilities.Append(utilities.Format(`DEFAULT %v`, json.Quoted(c.Default)), result, " ")
	}

	if len(c.Type) > 0 {
		result = utilities.Append(utilities.Uppcase(c.Type), result, " ")
	}
	if len(c.name) > 0 {
		result = utilities.Append(utilities.Uppcase(c.name), result, " ")
	}

	return result
}

func (c *Column) DDLIndex() string {
	return jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1($3);`, utilities.Lowcase(c.Model.Name), utilities.Uppcase(c.Model.Table), utilities.Uppcase(c.name))
}

/**
*
**/
func (c *Column) Up() string {
	return utilities.Uppcase(c.name)
}

func (c *Column) Low() string {
	return utilities.Lowcase(c.name)
}

/**
* This function not use in Select
**/
func (c *Column) As(linq *Linq) string {
	switch c.Tp {
	case TpColumn:
		from := linq.GetFrom(c)
		return utilities.Append(from.As(), c.Up(), ".")
	case TpAtrib:
		from := linq.GetFrom(c)
		col := utilities.Append(from.As(), c.Column.Up(), ".")
		return utilities.Format(`%s#>>'{%s}'`, col, c.Low())
	case TpClone:
		return utilities.Append(utilities.Uppcase(c.from), c.Up(), ".")
	case TpReference:
		from := linq.GetFrom(c)
		as := linq.GetAs()
		as = utilities.Format(`A%s`, as)
		fn := utilities.Format(`%s.%s`, as, c.Reference.Reference.Up())
		fm := utilities.Format(`%s AS %s`, c.Reference.Reference.Model.Name, as)
		key := utilities.Format(`%s.%s`, as, c.Reference.Key)
		Fkey := utilities.Append(from.As(), c.Reference.Fkey, ".")
		return utilities.Format(`(SELECT %s FROM %s WHERE %s=%v LIMIT 1)`, fn, fm, key, Fkey)
	case TpDetail:
		return utilities.Format(`%v`, json.Quoted(c.Default))
	case TpFunction:
		def := FunctionDef(linq, c)
		return utilities.Append(def, c.Up(), " AS ")
	case TpField:
		as := linq.As(c)
		if len(as) > 0 {
			as = utilities.Format(`%s.`, as)
		}
		def := utilities.Format(`(%v)`, c.Definition)
		return utilities.ReplaceAll(def, []string{"{AS}.", "{as}.", "{AS}", "{as}"}, as)
	default:
		return utilities.Format(`%s`, c.Up())
	}
}

/**
* This function use in Select
**/
func (c *Column) Def(linq *Linq) string {
	from := linq.GetFrom(c)

	if linq.Tp == TpData {
		switch c.Tp {
		case TpColumn:
			def := c.As(linq)
			return utilities.Format(`'%s', %s`, c.Low(), def)
		case TpAtrib:
			def := c.As(linq)
			def = utilities.Format(`COALESCE(%s, %v)`, def, json.Quoted(c.Default))
			return utilities.Format(`'%s', %s`, c.Low(), def)
		case TpClone:
			return c.As(linq)
		case TpReference:
			Fkey := utilities.Append(from.As(), c.Reference.Fkey, ".")
			def := c.As(linq)
			def = utilities.Format(`jsonb_build_object('_id', %s, 'name', %s)`, Fkey, def)
			return utilities.Format(`'%s', %s`, c.Title, def)
		case TpDetail:
			def := json.Quoted(c.Default)
			return utilities.Format(`'%s', %s`, c.Low(), def)
		case TpFunction:
			def := FunctionDef(linq, c)
			return utilities.Format(`'%s', %s`, c.Low(), def)
		case TpField:
			def := c.As(linq)
			return utilities.Format(`'%s', %s`, c.Low(), def)
		default:
			def := json.Quoted(c.Default)
			return utilities.Format(`'%s', %s`, c.Low(), def)
		}
	}

	switch c.Tp {
	case TpColumn:
		return utilities.Append(from.As(), c.Up(), ".")
	case TpAtrib:
		col := utilities.Append(from.As(), c.Column.Up(), ".")
		def := utilities.Format(`%s#>>'{%s}'`, col, c.Low())
		def = utilities.Format(`COALESCE(%s, %v)`, def, json.Quoted(c.Default))
		return utilities.Format(`%s AS %s`, def, c.Up())
	case TpClone:
		return utilities.Append(utilities.Uppcase(c.from), c.Up(), ".")
	case TpReference:
		def := c.As(linq)
		return utilities.Format(`%s AS %s`, def, utilities.Uppcase(c.Title))
	case TpDetail:
		def := json.Quoted(c.Default)
		return utilities.Format(`%v AS %s`, def, c.Up())
	case TpFunction:
		def := FunctionDef(linq, c)
		return utilities.Format(`%s AS %s`, def, c.Up())
	case TpField:
		def := c.As(linq)
		return utilities.Format(`%s AS %s`, def, c.Up())
	default:
		return utilities.Format(`%s`, c.Up())
	}
}

/**
* This function use in Select
**/
func (c *Column) From(from string) *Column {
	result := &Column{
		Model:   c.Model,
		Tp:      TpClone,
		name:    c.name,
		Type:    c.Type,
		Default: c.Default,
		Atribs:  c.Atribs,
		from:    from,
	}

	return result
}

func (c *Column) Name(name string) *Column {
	c.name = name

	return c
}

func (c *Column) Cast(cast string) *Column {
	c.cast = cast

	return c
}

/**
*
**/
func (c *Column) Eq(val any) *Where {
	return NewWhere(c, "=", val)
}

func (c *Column) Neg(val any) *Where {
	return NewWhere(c, "!=", val)
}

func (c *Column) In(vals ...any) *Where {
	return NewWhere(c, "IN", vals)
}

func (c *Column) Like(val any) *Where {
	if c.Tp == TpFunction {
		c.Cast("TEXT")
	}

	if val == "%"+"%" {
		val = "%"
	}

	return NewWhere(c, "ILIKE", val)
}

func (c *Column) More(val any) *Where {
	return NewWhere(c, ">", val)
}

func (c *Column) Less(val any) *Where {
	return NewWhere(c, "<", val)
}

func (c *Column) MoreEq(val any) *Where {
	return NewWhere(c, ">=", val)
}

func (c *Column) LessEq(val any) *Where {
	return NewWhere(c, "<=", val)
}
