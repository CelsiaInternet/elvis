package linq

import (
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

/**
*
**/
func (c *Model) Consolidate(linq *Linq) *Linq {
	var col *Column
	var source Json = Json{}
	var dta Json = Json{}
	var new Json = c.Model()

	setValue := func(key string, val interface{}) {
		dta.Set(key, val)
		new.Set(key, val)
	}

	for k, v := range linq.data {
		k = Lowcase(k)
		idxCol := c.ColIdx(k)

		if idxCol == -1 {
			idx := c.TitleIdx(k)
			if idx != -1 && ContainsInt([]int{TpReference}, c.Definition[idx].Tp) {
				col = c.Definition[idx]
				reference := linq.data.Json(k)
				setValue(col.name, reference.Key(col.Reference.Key))
				continue
			}
		}

		if idxCol == -1 && !c.integrityAtrib {
			source.Set(k, v)
			continue
		} else if idxCol == -1 {
			continue
		} else {
			col = c.Definition[idxCol]
		}

		if ContainsInt([]int{TpField, TpFunction, TpDetail}, col.Tp) {
			continue
		} else if k == Lowcase(c.SourceField) {
			atribs := linq.data.Json(k)

			if c.integrityAtrib {
				for ak, av := range atribs {
					ak = Lowcase(ak)
					if idx := c.AtribIdx(ak); idx != -1 {
						source[ak] = av
					}
				}
			} else {
				source = atribs
			}
		} else if ContainsInt([]int{TpColumn}, col.Tp) {
			setValue(k, v)
			col := c.Column(k)
			if col.PrimaryKey || col.ForeignKey {
				linq.references = append(linq.references, &ReferenceValue{c.Schema, c.Table, v, 1})
			}
		} else if ContainsInt([]int{TpAtrib}, col.Tp) {
			source.Set(k, v)
		}
	}

	if c.UseSource && len(source) > 0 {
		setValue(c.SourceField, source)
	}

	linq.dta = &dta
	linq.new = &new

	return linq
}

func (c *Model) Changue(current Json, linq *Linq) *Linq {
	var change bool
	dta := c.Consolidate(linq).dta

	for k, v := range current {
		linq.dta.Set(k, v)
	}

	for k, v := range *dta {
		k = Lowcase(k)
		idxCol := c.ColIdx(k)

		if idxCol != -1 {
			ch := current.Str(k) != dta.Str(k)
			if !change {
				change = ch
			}
			if ch {
				linq.dta.Set(k, v)
				col := c.Column(k)
				if col.PrimaryKey || col.ForeignKey {
					linq.references = append(linq.references, &ReferenceValue{c.Schema, c.Table, current.Str(k), -1})
					linq.references = append(linq.references, &ReferenceValue{c.Schema, c.Table, v, 1})
				}
			}
		}
	}

	linq.change = change

	return linq
}

/**
*
**/
func (c *Linq) PrepareInsert() {
	model := c.from[0].model
	model.Consolidate(c)
	now := Now()

	if model.UseDateMake {
		c.dta.Set(model.DateMakeField, now)
	}

	if model.UseDateUpdate {
		c.dta.Set(model.DateUpdateField, now)
	}
}

func (c *Linq) PrepareUpdate(current Json) bool {
	model := c.from[0].model
	model.Changue(current, c)

	if !c.change {
		return c.change
	}

	if model.UseDateMake {
		delete(*c.dta, Lowcase(model.DateMakeField))
	}

	now := Now()
	if model.UseDateUpdate {
		c.dta.Set(model.DateUpdateField, now)
	}

	return c.change
}

func (c *Linq) PrepareDelete(current Json) {
	model := c.from[0].model

	for k, v := range current {
		col := model.Column(k)
		if col.PrimaryKey || col.ForeignKey {
			c.references = append(c.references, &ReferenceValue{model.Schema, model.Table, v, -1})
		}
	}
}
