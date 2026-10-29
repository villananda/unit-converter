package main

import (
	"errors"
	"fmt"
	"math"
)

// absZeroTolerance absorbs float error when a temperature sits exactly at absolute zero.
const absZeroTolerance = 1e-9

// ErrBelowZero reports a value whose base quantity is negative (or below 0 K for temperature).
var ErrBelowZero = errors.New("result is below absolute zero or negative")

// Unit converts between its own scale and the category's base unit.
type Unit struct {
	Key      string
	Label    string
	toBase   func(float64) float64
	fromBase func(float64) float64
}

// Category groups units that share a base unit. Base quantities must be >= 0.
type Category struct {
	Key   string
	Label string
	Units []Unit
}

// scale returns conversion funcs for a unit whose base is `factor` times the base unit.
func scale(factor float64) (func(float64) float64, func(float64) float64) {
	return func(v float64) float64 { return v * factor },
		func(v float64) float64 { return v / factor }
}

func linear(key, label string, factor float64) Unit {
	to, from := scale(factor)
	return Unit{Key: key, Label: label, toBase: to, fromBase: from}
}

var categories = []Category{
	{
		Key:   "length",
		Label: "Length",
		Units: []Unit{
			linear("millimeter", "Millimeter (mm)", 0.001),
			linear("centimeter", "Centimeter (cm)", 0.01),
			linear("meter", "Meter (m)", 1),
			linear("kilometer", "Kilometer (km)", 1000),
			linear("inch", "Inch (in)", 0.0254),
			linear("foot", "Foot (ft)", 0.3048),
			linear("yard", "Yard (yd)", 0.9144),
			linear("mile", "Mile (mi)", 1609.344),
		},
	},
	{
		Key:   "weight",
		Label: "Weight",
		Units: []Unit{
			linear("milligram", "Milligram (mg)", 0.001),
			linear("gram", "Gram (g)", 1),
			linear("kilogram", "Kilogram (kg)", 1000),
			linear("ounce", "Ounce (oz)", 28.349523125),
			linear("pound", "Pound (lb)", 453.59237),
		},
	},
	{
		// Base unit is Kelvin, so the non-negative rule means "not below absolute zero".
		Key:   "temperature",
		Label: "Temperature",
		Units: []Unit{
			{
				Key:      "celsius",
				Label:    "Celsius (°C)",
				toBase:   func(v float64) float64 { return v + 273.15 },
				fromBase: func(v float64) float64 { return v - 273.15 },
			},
			{
				Key:      "fahrenheit",
				Label:    "Fahrenheit (°F)",
				toBase:   func(v float64) float64 { return (v-32)*5/9 + 273.15 },
				fromBase: func(v float64) float64 { return (v-273.15)*9/5 + 32 },
			},
			{
				Key:      "kelvin",
				Label:    "Kelvin (K)",
				toBase:   func(v float64) float64 { return v },
				fromBase: func(v float64) float64 { return v },
			},
		},
	},
}

func findCategory(key string) (Category, bool) {
	for _, c := range categories {
		if c.Key == key {
			return c, true
		}
	}
	return Category{}, false
}

func (c Category) unit(key string) (Unit, bool) {
	for _, u := range c.Units {
		if u.Key == key {
			return u, true
		}
	}
	return Unit{}, false
}

// Convert converts value from one unit to another within a category.
func Convert(categoryKey, from, to string, value float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, errors.New("value must be a finite number")
	}
	cat, ok := findCategory(categoryKey)
	if !ok {
		return 0, fmt.Errorf("unknown category: %q", categoryKey)
	}
	src, ok := cat.unit(from)
	if !ok {
		return 0, fmt.Errorf("unknown source unit: %q", from)
	}
	dst, ok := cat.unit(to)
	if !ok {
		return 0, fmt.Errorf("unknown target unit: %q", to)
	}

	base := src.toBase(value)
	if base < -absZeroTolerance {
		return 0, ErrBelowZero
	}
	if base < 0 {
		base = 0
	}
	return dst.fromBase(base), nil
}
