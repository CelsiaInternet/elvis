package linq

import (
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
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
	return Uppcase(c.name)
}

func (c *Col) Low() string {
	return Lowcase(c.name)
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
	return Uppcase(c.As())
}

func (c *Col) AsLow() string {
	return Lowcase(c.As())
}

/**
*
**/
type Details func(col *Column, data *Json)

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

func (c *Column) describe() Json {
	return Json{
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

func (c *Column) Describe() Json {
	var atribs []Json = []Json{}
	for _, atrib := range c.Atribs {
		atribs = append(atribs, atrib.describe())
	}

	reference := Json{}
	if c.Reference != nil {
		reference = c.Reference.Describe()
	}

	return Json{
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
		name:        Uppcase(name),
		Description: description,
		Type:        _type,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
	}

	if !model.UseDateMake {
		model.UseDateMake = Uppcase(result.name) == Uppcase(model.DateMakeField)
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = Uppcase(result.name) == Uppcase(model.DateUpdateField)
	}

	if !model.UseState {
		model.UseState = Uppcase(result.name) == Uppcase(model.StateField)
	}

	if !model.UseRecycle {
		model.UseRecycle = Uppcase(result.name) == Uppcase(model.StateField)
	}

	if !model.UseProject {
		model.UseProject = Uppcase(result.name) == Uppcase(model.ProjectField)
	}

	if !model.UseIndex {
		model.UseIndex = Uppcase(result.name) == Uppcase(model.IndexField)
	}

	if !model.UseSource {
		model.UseSource = Uppcase(result.name) == Uppcase(model.SourceField)
	}

	model.Definition = append(model.Definition, result)
	return result
}

func NewVirtualAtrib(model *Model, name, description, _type string, _default any) *Column {
	result := &Column{
		Model:       model,
		Tp:          TpColumn,
		name:        Uppcase(name),
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

	_default := NewAny(c.Default)

	if _default.String() == "NOW()" {
		result = Append(`DEFAULT NOW()`, result, " ")
	} else {
		result = Append(Format(`DEFAULT %v`, Quoted(c.Default)), result, " ")
	}

	if len(c.Type) > 0 {
		result = Append(Uppcase(c.Type), result, " ")
	}
	if len(c.name) > 0 {
		result = Append(Uppcase(c.name), result, " ")
	}

	return result
}

func (c *Column) DDLIndex() string {
	return SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1($3);`, Lowcase(c.Model.Name), Uppcase(c.Model.Table), Uppcase(c.name))
}

/**
*
**/
func (c *Column) Up() string {
	return Uppcase(c.name)
}

func (c *Column) Low() string {
	return Lowcase(c.name)
}

/**
* This function not use in Select
**/
func (c *Column) As(linq *Linq) string {
	switch c.Tp {
	case TpColumn:
		from := linq.GetFrom(c)
		return Append(from.As(), c.Up(), ".")
	case TpAtrib:
		from := linq.GetFrom(c)
		col := Append(from.As(), c.Column.Up(), ".")
		return Format(`%s#>>'{%s}'`, col, c.Low())
	case TpClone:
		return Append(Uppcase(c.from), c.Up(), ".")
	case TpReference:
		from := linq.GetFrom(c)
		as := linq.GetAs()
		as = Format(`A%s`, as)
		fn := Format(`%s.%s`, as, c.Reference.Reference.Up())
		fm := Format(`%s AS %s`, c.Reference.Reference.Model.Name, as)
		key := Format(`%s.%s`, as, c.Reference.Key)
		Fkey := Append(from.As(), c.Reference.Fkey, ".")
		return Format(`(SELECT %s FROM %s WHERE %s=%v LIMIT 1)`, fn, fm, key, Fkey)
	case TpDetail:
		return Format(`%v`, Quoted(c.Default))
	case TpFunction:
		def := FunctionDef(linq, c)
		return Append(def, c.Up(), " AS ")
	case TpField:
		as := linq.As(c)
		if len(as) > 0 {
			as = Format(`%s.`, as)
		}
		def := Format(`(%v)`, c.Definition)
		return ReplaceAll(def, []string{"{AS}.", "{as}.", "{AS}", "{as}"}, as)
	default:
		return Format(`%s`, c.Up())
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
			return Format(`'%s', %s`, c.Low(), def)
		case TpAtrib:
			def := c.As(linq)
			def = Format(`COALESCE(%s, %v)`, def, Quoted(c.Default))
			return Format(`'%s', %s`, c.Low(), def)
		case TpClone:
			return c.As(linq)
		case TpReference:
			Fkey := Append(from.As(), c.Reference.Fkey, ".")
			def := c.As(linq)
			def = Format(`jsonb_build_object('_id', %s, 'name', %s)`, Fkey, def)
			return Format(`'%s', %s`, c.Title, def)
		case TpDetail:
			def := Quoted(c.Default)
			return Format(`'%s', %s`, c.Low(), def)
		case TpFunction:
			def := FunctionDef(linq, c)
			return Format(`'%s', %s`, c.Low(), def)
		case TpField:
			def := c.As(linq)
			return Format(`'%s', %s`, c.Low(), def)
		default:
			def := Quoted(c.Default)
			return Format(`'%s', %s`, c.Low(), def)
		}
	}

	switch c.Tp {
	case TpColumn:
		return Append(from.As(), c.Up(), ".")
	case TpAtrib:
		col := Append(from.As(), c.Column.Up(), ".")
		def := Format(`%s#>>'{%s}'`, col, c.Low())
		def = Format(`COALESCE(%s, %v)`, def, Quoted(c.Default))
		return Format(`%s AS %s`, def, c.Up())
	case TpClone:
		return Append(Uppcase(c.from), c.Up(), ".")
	case TpReference:
		def := c.As(linq)
		return Format(`%s AS %s`, def, Uppcase(c.Title))
	case TpDetail:
		def := Quoted(c.Default)
		return Format(`%v AS %s`, def, c.Up())
	case TpFunction:
		def := FunctionDef(linq, c)
		return Format(`%s AS %s`, def, c.Up())
	case TpField:
		def := c.As(linq)
		return Format(`%s AS %s`, def, c.Up())
	default:
		return Format(`%s`, c.Up())
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
