//go:generate gopherjs build github.com/antha-lang/units/js -m -o lib.js
package main

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)
import "github.com/antha-lang/units"

func main() {
	js.Global.Set("Antha", map[string]interface{}{
		"units": js.MakeWrapper(&ulib{}),
	})
}

type ulib struct{}

// Measurement
type Measurement struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

func (a Measurement) Quantity() float64 {
	return a.Value
}

func (a Measurement) MeasurementUnit() string {
	return a.Unit
}

func wrapError(err error) *js.Object {
	if err == nil {
		return nil
	}
	return js.MakeWrapper(err)
}

func (*ulib) Measurement(v float64, unit string) *js.Object {
	return js.MakeWrapper(Measurement{
		Value: v,
		Unit:  unit,
	})
}

func (*ulib) New(unitString string, ms ...Measurement) (*js.Object, *js.Object) {
	if len(ms) == 0 {
		return js.MakeWrapper(Measurement{}), wrapError(errors.New("not enough arguments"))
	}

	var casts []units.Measurement
	for _, v := range ms {
		casts = append(casts, v)
	}

	m, err := units.New(unitString, casts[0], casts[1:]...)

	obj := js.MakeWrapper(Measurement{
		Value: m.Quantity(),
		Unit:  m.MeasurementUnit(),
	})

	return obj, wrapError(err)
}

func (*ulib) Parse(quantity float64, unitString string) (*js.Object, *js.Object) {
	m, err := units.Parse(quantity, unitString)

	obj := js.MakeWrapper(Measurement{
		Value: m.Quantity(),
		Unit:  m.MeasurementUnit(),
	})

	return obj, wrapError(err)
}
