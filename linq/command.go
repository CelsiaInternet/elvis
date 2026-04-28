package linq

import (
	"context"

	"github.com/celsiainternet/elvis/et"
)

func (c *Linq) Debug() *Linq {
	c.debug = true

	return c
}

/**
* Command
* @return et.Items
* @return error
**/
func (c *Linq) Command() (et.Items, error) {
	if c.Act == ActInsert {
		return c.commandInsert()
	}

	if c.Act == ActUpdate {
		return c.commandUpdate()
	}

	if c.Act == ActUpsert {
		return c.commandUpsert()
	}

	if c.Act == ActDelete {
		return c.commandDelete()
	}

	return et.Items{}, nil
}

/**
* CommandOne
* @return et.Item
* @return error
**/
func (c *Linq) CommandOne() (et.Item, error) {
	result, err := c.Command()
	if err != nil {
		return et.Item{}, err
	}

	if result.Count == 0 {
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* Go
* @return et.Item
* @return error
**/
func (c *Linq) Go() (et.Item, error) {
	return c.CommandOne()
}

/**
* commandInsert uses ON CONFLICT DO NOTHING in the SQL so duplicate
* detection is atomic and requires only one round-trip to the database.
* result.Ok == false means the record already existed (conflict).
**/
func (c *Linq) commandInsert() (et.Items, error) {
	if err := c.prepareInsertData(); err != nil {
		return et.Items{}, err
	}

	result, err := c.insert()
	if err != nil {
		return et.Items{}, err
	}

	if !result.Ok {
		return et.Items{Ok: false, Count: 0}, nil
	}

	return et.Items{
		Ok:     true,
		Count:  1,
		Result: []et.Json{result.Result},
	}, nil
}

func (c *Linq) commandUpdate() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareUpdate()
	if err != nil {
		return et.Items{}, err
	}

	if currents.Count == 0 {
		return result, nil
	}

	tx, err := c.db.BeginTx(context.Background())
	if err != nil {
		return et.Items{}, err
	}
	c.tx = tx

	rollback := func() {
		tx.Rollback()
		c.tx = nil
	}

	model := c.from[0].model
	for _, current := range currents.Result {
		model.Changue(current, c)
		if c.change {
			item, err := c.update(current)
			if err != nil {
				rollback()
				return et.Items{}, err
			}
			result.Result = append(result.Result, item.Result)
			result.Ok = true
			result.Count++
		} else {
			result.Result = append(result.Result, current)
			result.Ok = true
			result.Count++
		}
	}

	if err := tx.Commit(); err != nil {
		rollback()
		return et.Items{}, err
	}
	c.tx = nil

	return result, nil
}

func (c *Linq) commandDelete() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareDelete()
	if err != nil {
		return et.Items{}, err
	}

	if currents.Count == 0 {
		return result, nil
	}

	tx, err := c.db.BeginTx(context.Background())
	if err != nil {
		return et.Items{}, err
	}
	c.tx = tx

	rollback := func() {
		tx.Rollback()
		c.tx = nil
	}

	for _, current := range currents.Result {
		item, err := c.delete(current)
		if err != nil {
			rollback()
			return et.Items{}, err
		}
		result.Result = append(result.Result, item.Result)
		result.Ok = true
		result.Count++
	}

	if err := tx.Commit(); err != nil {
		rollback()
		return et.Items{}, err
	}
	c.tx = nil

	return result, nil
}

func (c *Linq) commandUpsert() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareUpsert()
	if err != nil {
		return et.Items{}, err
	}

	if currents.Count == 0 {
		item, err := c.insert()
		if err != nil {
			return et.Items{}, err
		}

		if item.Ok {
			result.Result = append(result.Result, item.Result)
			result.Ok = true
			result.Count++
		}
	} else {
		model := c.from[0].model
		for _, current := range currents.Result {
			model.Changue(current, c)
			if c.change {
				item, err := c.update(current)
				if err != nil {
					return et.Items{}, err
				} else {
					result.Result = append(result.Result, item.Result)
					result.Ok = true
					result.Count++
				}
			} else {
				result.Result = append(result.Result, current)
				result.Ok = true
				result.Count++
			}
		}
	}

	return result, nil
}

/**
*
**/
func (c *Linq) Current() (et.Items, error) {
	c.SqlCurrent()
	return c.db.Query(c.sql)
}

/**
* Basic operation
**/
func (c *Linq) insert() (et.Item, error) {
	model := c.from[0].model

	for _, trigger := range model.BeforeInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlInsert()
	items, err := c.command()
	if err != nil {
		return et.Item{}, err
	}

	item := items.First()
	if !item.Ok {
		return item, nil
	}

	new := &item.Result

	for _, trigger := range model.AfterInsert {
		err := trigger(model, nil, new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.Details(new)

	return item, nil
}

func (c *Linq) update(current et.Json) (et.Item, error) {
	model := c.from[0].model
	c.idT = current.ValStr("-1", IdTFiled.Low())

	for _, trigger := range model.BeforeUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlUpdate()
	items, err := c.command()
	if err != nil {
		return et.Item{}, err
	}

	item := items.First()
	if !item.Ok {
		return item, nil
	}

	new := &item.Result

	for _, trigger := range model.AfterUpdate {
		err := trigger(model, &current, new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.Details(new)

	return item, nil
}

func (c *Linq) delete(current et.Json) (et.Item, error) {
	model := c.from[0].model
	c.idT = current.ValStr("-1", IdTFiled.Low())

	for _, trigger := range model.BeforeDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlDelete()
	items, err := c.command()
	if err != nil {
		return et.Item{}, err
	}

	item := items.First()
	if !item.Ok {
		return item, nil
	}

	for _, trigger := range model.AfterDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	return et.Item{
		Ok:     true,
		Result: current,
	}, nil
}
