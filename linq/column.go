package linq

import (
	"github.com/cgalvisleon/elvis/generic"
	"github.com/cgalvisleon/elvis/jdb"
	j "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
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
	return utility.Uppcase(c.name)
}

func (c *Col) Low() string {
	return utility.Lowcase(c.name)
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
	return utility.Uppcase(c.As())
}

func (c *Col) AsLow() string {
	return utility.Lowcase(c.As())
}

/**
*
**/
type Details func(col *Column, data *j.Json)

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

func (c *Column) describe() j.Json {
	return j.Json{
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

func (c *Column) Describe() j.Json {
	var atribs []j.Json = []j.Json{}
	for _, atrib := range c.Atribs {
		atribs = append(atribs, atrib.describe())
	}

	reference := j.Json{}
	if c.Reference != nil {
		reference = c.Reference.Describe()
	}

	return j.Json{
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
		name:        utility.Uppcase(name),
		Description: description,
		Type:        _type,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
	}

	if !model.UseDateMake {
		model.UseDateMake = utility.Uppcase(result.name) == utility.Uppcase(model.DateMakeField)
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = utility.Uppcase(result.name) == utility.Uppcase(model.DateUpdateField)
	}

	if !model.UseState {
		model.UseState = utility.Uppcase(result.name) == utility.Uppcase(model.StateField)
	}

	if !model.UseRecycle {
		model.UseRecycle = utility.Uppcase(result.name) == utility.Uppcase(model.StateField)
	}

	if !model.UseProject {
		model.UseProject = utility.Uppcase(result.name) == utility.Uppcase(model.ProjectField)
	}

	if !model.UseIndex {
		model.UseIndex = utility.Uppcase(result.name) == utility.Uppcase(model.IndexField)
	}

	if !model.UseSource {
		model.UseSource = utility.Uppcase(result.name) == utility.Uppcase(model.SourceField)
	}

	model.Definition = append(model.Definition, result)
	return result
}

func NewVirtualAtrib(model *Model, name, description, _type string, _default any) *Column {
	result := &Column{
		Model:       model,
		Tp:          TpColumn,
		name:        utility.Uppcase(name),
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

	_default := generic.New(c.Default)

	if _default.Str() == "NOW()" {
		result = utility.Append(`DEFAULT NOW()`, result, " ")
	} else {
		result = utility.Append(utility.Format(`DEFAULT %v`, j.Quoted(c.Default)), result, " ")
	}

	if len(c.Type) > 0 {
		result = utility.Append(utility.Uppcase(c.Type), result, " ")
	}
	if len(c.name) > 0 {
		result = utility.Append(utility.Uppcase(c.name), result, " ")
	}

	return result
}

func (c *Column) DDLIndex() string {
	return jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1($3);`, utility.Lowcase(c.Model.Name), utility.Uppcase(c.Model.Table), utility.Uppcase(c.name))
}

/**
*
**/
func (c *Column) Up() string {
	return utility.Uppcase(c.name)
}

func (c *Column) Low() string {
	return utility.Lowcase(c.name)
}

/**
* This function not use in Select
**/
func (c *Column) As(linq *Linq) string {
	switch c.Tp {
	case TpColumn:
		from := linq.GetFrom(c)
		return utility.Append(from.As(), c.Up(), ".")
	case TpAtrib:
		from := linq.GetFrom(c)
		col := utility.Append(from.As(), c.Column.Up(), ".")
		return utility.Format(`%s#>>'{%s}'`, col, c.Low())
	case TpClone:
		return utility.Append(utility.Uppcase(c.from), c.Up(), ".")
	case TpReference:
		from := linq.GetFrom(c)
		as := linq.GetAs()
		as = utility.Format(`A%s`, as)
		fn := utility.Format(`%s.%s`, as, c.Reference.Reference.Up())
		fm := utility.Format(`%s AS %s`, c.Reference.Reference.Model.Name, as)
		key := utility.Format(`%s.%s`, as, c.Reference.Key)
		Fkey := utility.Append(from.As(), c.Reference.Fkey, ".")
		return utility.Format(`(SELECT %s FROM %s WHERE %s=%v LIMIT 1)`, fn, fm, key, Fkey)
	case TpDetail:
		return utility.Format(`%v`, j.Quoted(c.Default))
	case TpFunction:
		def := FunctionDef(linq, c)
		return utility.Append(def, c.Up(), " AS ")
	case TpField:
		as := linq.As(c)
		if len(as) > 0 {
			as = utility.Format(`%s.`, as)
		}
		def := utility.Format(`(%v)`, c.Definition)
		return utility.ReplaceAll(def, []string{"{AS}.", "{as}.", "{AS}", "{as}"}, as)
	default:
		return utility.Format(`%s`, c.Up())
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
			return utility.Format(`'%s', %s`, c.Low(), def)
		case TpAtrib:
			def := c.As(linq)
			def = utility.Format(`COALESCE(%s, %v)`, def, j.Quoted(c.Default))
			return utility.Format(`'%s', %s`, c.Low(), def)
		case TpClone:
			return c.As(linq)
		case TpReference:
			Fkey := utility.Append(from.As(), c.Reference.Fkey, ".")
			def := c.As(linq)
			def = utility.Format(`jsonb_build_object('_id', %s, 'name', %s)`, Fkey, def)
			return utility.Format(`'%s', %s`, c.Title, def)
		case TpDetail:
			def := j.Quoted(c.Default)
			return utility.Format(`'%s', %s`, c.Low(), def)
		case TpFunction:
			def := FunctionDef(linq, c)
			return utility.Format(`'%s', %s`, c.Low(), def)
		case TpField:
			def := c.As(linq)
			return utility.Format(`'%s', %s`, c.Low(), def)
		default:
			def := j.Quoted(c.Default)
			return utility.Format(`'%s', %s`, c.Low(), def)
		}
	}

	switch c.Tp {
	case TpColumn:
		return utility.Append(from.As(), c.Up(), ".")
	case TpAtrib:
		col := utility.Append(from.As(), c.Column.Up(), ".")
		def := utility.Format(`%s#>>'{%s}'`, col, c.Low())
		def = utility.Format(`COALESCE(%s, %v)`, def, j.Quoted(c.Default))
		return utility.Format(`%s AS %s`, def, c.Up())
	case TpClone:
		return utility.Append(utility.Uppcase(c.from), c.Up(), ".")
	case TpReference:
		def := c.As(linq)
		return utility.Format(`%s AS %s`, def, utility.Uppcase(c.Title))
	case TpDetail:
		def := j.Quoted(c.Default)
		return utility.Format(`%v AS %s`, def, c.Up())
	case TpFunction:
		def := FunctionDef(linq, c)
		return utility.Format(`%s AS %s`, def, c.Up())
	case TpField:
		def := c.As(linq)
		return utility.Format(`%s AS %s`, def, c.Up())
	default:
		return utility.Format(`%s`, c.Up())
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
