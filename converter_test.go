package main

import (
	"errors"
	"math"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		category string
		from     string
		to       string
		value    float64
		want     float64
	}{
		{"meter to foot", "length", "meter", "foot", 1, 1 / 0.3048},
		{"mile to kilometer", "length", "mile", "kilometer", 1, 1.609344},
		{"pound to kilogram", "weight", "pound", "kilogram", 1, 0.45359237},
		{"celsius to fahrenheit", "temperature", "celsius", "fahrenheit", 100, 212},
		{"fahrenheit to kelvin", "temperature", "fahrenheit", "kelvin", -40, 233.15},
		{"kelvin to celsius", "temperature", "kelvin", "celsius", 0, -273.15},
		{"absolute zero edge", "temperature", "fahrenheit", "kelvin", -459.67, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.category, tt.from, tt.to, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertRejects(t *testing.T) {
	tests := []struct {
		name     string
		category string
		from     string
		to       string
		value    float64
		wantErr  error
	}{
		{"below absolute zero", "temperature", "celsius", "kelvin", -300, ErrBelowZero},
		{"negative length", "length", "meter", "meter", -1, ErrBelowZero},
		{"NaN value", "length", "meter", "meter", math.NaN(), nil},
		{"unknown category", "volume", "meter", "meter", 1, nil},
		{"unknown unit", "length", "parsec", "meter", 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(tt.category, tt.from, tt.to, tt.value)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
