package linq

import (
	j "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/msg"
)

func (c *Linq) Debug() *Linq {
	c.debug = true

	return c
}

/**
* Executors
**/
func (c *Linq) Command() (j.Item, error) {
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

	return j.Item{}, nil
}

/**
* Exec
**/
func (c *Linq) commandInsert() (j.Item, error) {
	if len(c.where) > 0 {
		current, err := c.Current()
		if err != nil {
			return j.Item{}, err
		}

		if current.Ok {
			return j.Item{
				Ok: !current.Ok,
				Result: j.Json{
					"message": msg.RECORD_FOUND,
				},
			}, nil
		}
	}

	return c.insert()
}

func (c *Linq) commandUpdate() (j.Item, error) {
	current, err := c.Current()
	if err != nil {
		return j.Item{}, err
	}

	if !current.Ok {
		return j.Item{
			Ok: current.Ok,
			Result: j.Json{
				"message": msg.RECORD_NOT_FOUND,
			},
		}, nil
	}

	return c.update(current.Result)
}

func (c *Linq) commandDelete() (j.Item, error) {
	current, err := c.Current()
	if err != nil {
		return j.Item{}, err
	}

	if !current.Ok {
		return j.Item{
			Ok: current.Ok,
			Result: j.Json{
				"message": msg.RECORD_NOT_FOUND,
			},
		}, nil
	}

	return c.delete(current.Result)
}

func (c *Linq) commandUpsert() (j.Item, error) {
	current, err := c.Current()
	if err != nil {
		return j.Item{}, err
	}

	if current.Ok {
		return c.update(current.Result)
	}

	return c.insert()
}

/**
*
**/
func (c *Linq) Current() (j.Item, error) {
	c.sql = c.SqlCurrent()

	return c.QueryOne()
}

/**
* Basic operation
**/
func (c *Linq) insert() (j.Item, error) {
	c.PrepareInsert()
	model := c.from[0].model

	for _, trigger := range model.BeforeInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.SqlInsert()

	item, err := c.QueryOne()
	if err != nil {
		return j.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.Details(&item.Result)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return item, nil
}

func (c *Linq) update(current j.Json) (j.Item, error) {
	changue := c.PrepareUpdate(current)
	if !changue {
		return j.Item{
			Ok: changue,
			Result: j.Json{
				"message": msg.RECORD_NOT_CHANGE,
			},
		}, nil
	}

	model := c.from[0].model

	for _, trigger := range model.BeforeUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.SqlUpdate()

	item, err := c.QueryOne()
	if err != nil {
		return j.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.Details(&item.Result)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return item, nil
}

func (c *Linq) delete(current j.Json) (j.Item, error) {
	c.PrepareDelete(current)
	model := c.from[0].model

	for _, trigger := range model.BeforeDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.SqlDelete()

	item, err := c.QueryOne()
	if err != nil {
		return j.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	for _, trigger := range model.AfterDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return j.Item{}, err
		}
	}

	c.Details(&current)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return j.Item{
		Ok: true,
		Result: current,
	}, nil
}
