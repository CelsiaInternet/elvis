package linq

import (
	"github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/msg"
)

func (c *Linq) Debug() *Linq {
	c.debug = true

	return c
}

/**
* Executors
**/
func (c *Linq) Command() (json.Item, error) {
	if c.Act == ActInsert {
		return c.commandInsert()
	}

	if c.Act == ActUpdate {
		return c.commandUpdate()
	}

	if c.Act == ActDelete {
		return c.commandDelete()
	}

	if c.Act == ActUpsert {
		return c.commandUpsert()
	}

	return json.Item{}, nil
}

/**
* Exec
**/
func (c *Linq) commandInsert() (json.Item, error) {
	if len(c.where) > 0 {
		current, err := c.Current()
		if err != nil {
			return json.Item{}, err
		}

		if current.Ok {
			return json.Item{
				Ok: !current.Ok,
				Result: json.Json{
					"message": msg.RECORD_FOUND,
				},
			}, nil
		}
	}

	return c.insert()
}

func (c *Linq) commandUpdate() (json.Item, error) {
	current, err := c.Current()
	if err != nil {
		return json.Item{}, err
	}

	if !current.Ok {
		return json.Item{
			Ok: current.Ok,
			Result: json.Json{
				"message": msg.RECORD_NOT_FOUND,
			},
		}, nil
	}

	return c.update(current.Result)
}

func (c *Linq) commandDelete() (json.Item, error) {
	current, err := c.Current()
	if err != nil {
		return json.Item{}, err
	}

	if !current.Ok {
		return json.Item{
			Ok: current.Ok,
			Result: json.Json{
				"message": msg.RECORD_NOT_FOUND,
			},
		}, nil
	}

	return c.delete(current.Result)
}

func (c *Linq) commandUpsert() (json.Item, error) {
	current, err := c.Current()
	if err != nil {
		return json.Item{}, err
	}

	if current.Ok {
		return c.update(current.Result)
	}

	return c.insert()
}

/**
*
**/
func (c *Linq) Current() (json.Item, error) {
	c.sql = c.SqlCurrent()

	return c.QueryOne()
}

/**
* Basic operation
**/
func (c *Linq) insert() (json.Item, error) {
	c.PrepareInsert()
	model := c.from[0].model

	for _, trigger := range model.BeforeInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.SqlInsert()

	item, err := c.QueryOne()
	if err != nil {
		return json.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.Details(&item.Result)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return item, nil
}

func (c *Linq) update(current json.Json) (json.Item, error) {
	changue := c.PrepareUpdate(current)
	if !changue {
		return json.Item{
			Ok: changue,
			Result: json.Json{
				"message": msg.RECORD_NOT_CHANGE,
			},
		}, nil
	}

	model := c.from[0].model

	for _, trigger := range model.BeforeUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.SqlUpdate()

	item, err := c.QueryOne()
	if err != nil {
		return json.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.Details(&item.Result)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return item, nil
}

func (c *Linq) delete(current json.Json) (json.Item, error) {
	c.PrepareDelete(current)
	model := c.from[0].model

	for _, trigger := range model.BeforeDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.SqlDelete()

	item, err := c.QueryOne()
	if err != nil {
		return json.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	for _, trigger := range model.AfterDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return json.Item{}, err
		}
	}

	c.Details(&item.Result)

	if model.AfterReferences != nil {
		go model.AfterReferences(c.references)
	}

	return item, nil
}
