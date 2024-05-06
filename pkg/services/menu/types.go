package menu

import "encoding/json"

type Dish struct {
	Name   string
	Amount float32
}

type Meal struct {
	Name   string
	Dishes []Dish `json:",omitempty"`
}

type Day struct {
	Name  string
	Meals []Meal `json:",omitempty"`
}

type Menu struct {
	Days []Day `json:",omitempty"`
}

func (menu *Menu) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &menu.Days)
}

func (menu Menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(menu.Days)
}

type ProductData struct {
	Name     string
	Amount   float32 `json:",omitempty"`
	UnitCost float32 `json:"unit_cost,omitempty"`
}
