package types

import "github.com/go-playground/validator"

type Pager struct {
	Offset int `query:"offset" validate:"min=0"`
	Limit  int `query:"limit" validate:"max=10"`
}

func (p *Pager) Validate() error {
	if p.Limit <= 0 {
		p.Limit = 10
	}

	var structValidator = validator.New()
	structError := structValidator.Struct(p)

	if structError != nil {
		return structError
	}

	return nil
}

type PagerResult struct {
	Next     int         `json:"next"`
	Previous int         `json:"previous"`
	Total    int         `json:"total"`
	Data     interface{} `json:"data"`
}

func (p *PagerResult) GetPagerResult(pager *Pager, totalRows int, data interface{}) *PagerResult {
	var current int

	if pager.Offset == 0 {
		current = 1 * pager.Limit
	} else {
		current = pager.Offset * pager.Limit
	}

	if current < totalRows {
		p.Next = pager.Offset + 1
	} else {
		p.Next = 0
	}

	if current > pager.Limit {
		p.Previous = pager.Offset - 1
	} else {
		p.Previous = 0
	}

	p.Total = totalRows
	p.Data = data

	return p
}
