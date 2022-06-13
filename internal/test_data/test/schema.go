// Code generated by https://github.com/meerkat-lib/disorder; DO NOT EDIT.
package test

import (
	"fmt"
	"github.com/meerkat-lib/disorder/internal/test_data/test/sub"
	"time"
)

type Color int

const (
	ColorRed Color = iota
	ColorGreen
	ColorBlue
)

func (*Color) Enum() {}

func (enum *Color) FromString(value string) error {
	switch {
	case value == "red":
		*enum = ColorRed
		return nil
	case value == "green":
		*enum = ColorGreen
		return nil
	case value == "blue":
		*enum = ColorBlue
		return nil
	}
	return fmt.Errorf("invalid enum value")
}

func (enum *Color) ToString() (string, error) {
	switch *enum {
	case ColorRed:
		return "red", nil
	case ColorGreen:
		return "green", nil
	case ColorBlue:
		return "blue", nil
	default:
		return "", fmt.Errorf("invalid enum value")
	}
}

type Animal int

const (
	AnimalDog Animal = iota
	AnimalCat
)

func (*Animal) Enum() {}

func (enum *Animal) FromString(value string) error {
	switch {
	case value == "dog":
		*enum = AnimalDog
		return nil
	case value == "cat":
		*enum = AnimalCat
		return nil
	}
	return fmt.Errorf("invalid enum value")
}

func (enum *Animal) ToString() (string, error) {
	switch *enum {
	case AnimalDog:
		return "dog", nil
	case AnimalCat:
		return "cat", nil
	default:
		return "", fmt.Errorf("invalid enum value")
	}
}

type Object struct {
	ObjArray    []*sub.SubObject          `disorder:"obj_array"`
	ObjMap      map[string]*sub.SubObject `disorder:"obj_map"`
	Time        *time.Time                `disorder:"time"`
	IntField    int32                     `disorder:"int_field"`
	StringField string                    `disorder:"string_field"`
	EnumField   *Color                    `disorder:"enum_field"`
	IntArray    []int32                   `disorder:"int_array"`
	IntMap      map[string]int32          `disorder:"int_map"`
}
