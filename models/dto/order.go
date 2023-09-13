package dto

import "fmt"

type OrderDirection string

const (
	OrderByASC      OrderDirection = "ASC"
	OrderByDESC     OrderDirection = "DESC"
	OrderDefaultJey                = "record_id"
)

type OrderParam struct {
	Key       string         `query:"order_key"`
	Direction OrderDirection `query:"order_direction"`
}

func (o *OrderParam) ParseOrder() string {
	if o.Key == "" {
		o.Key = OrderDefaultJey
	}

	key := o.Key
	direction := "DESC"
	if o.Direction == OrderByASC {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", key, direction)
}
