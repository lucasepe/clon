package parser

import (
	"fmt"
	"os"
	"strconv"
)

type Any interface{}

type Value struct {
	value Any
}

func (v Value) Do() Any {
	return v.value
}

func (v Value) Merge(g Generator) Generator {
	// Values cannot be merged. Return the new one
	return g
}

type Generator interface {
	Do() Any
	Merge(Generator) Generator
}

type Obj struct {
	fields map[string]Generator
}

func NewObj() Obj {
	return Obj{
		fields: map[string]Generator{},
	}
}

func (obj Obj) Do() Any {
	res := map[string]Any{}
	for field, valueGen := range obj.fields {
		res[field] = valueGen.Do()
	}
	return res
}

func (obj Obj) Merge(g Generator) Generator {
	switch g := g.(type) {
	case Obj:
		// Objects can be merged together
		res := NewObj()
		for f, v := range obj.fields {
			res.Add(f, v)
		}
		for f, v := range g.fields {
			res.Add(f, v)
		}
		return res
	default:
		// other types, less so, return the new one
		return g
	}
}

func (obj Obj) Add(field string, value Generator) Obj {
	existingGenerator, found := obj.fields[field]
	if found {
		value = existingGenerator.Merge(value)
	}
	obj.fields[field] = value
	return obj
}

type Arr []Generator

func (arr Arr) Merge(g Generator) Generator {
	// arrays can' t be merged with other generators
	return g
}

func (arr Arr) Do() Any {
	res := make([]Any, len(arr))
	for idx, elemGen := range arr {
		res[idx] = elemGen.Do()
	}
	return res
}

func (arr *Arr) Add(g Generator) Generator {
	*arr = append(*arr, g)
	return *arr
}

func parseRawValue(value string) (Any, error) {
	switch {
	case value == "true":
		return true, nil
	case value == "false":
		return false, nil
	case value == "null":
		return nil, nil
	default:
		var v Any
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			v, err = strconv.ParseFloat(value, 64)
		}
		if err != nil {
			return nil, fmt.Errorf("Invalid raw literal %q: isn't any of true, false, null or a numeric", value)
		}
		return v, nil
	}
}

func parseEnvValue(value string) (Any, error) {
	res, ok := os.LookupEnv(value)
	if !ok {
		return "", nil
	}

	switch {
	case value == "true":
		return true, nil
	case value == "false":
		return false, nil
	case value == "null":
		return nil, nil
	default:
		var v Any
		v, err := strconv.ParseInt(res, 10, 64)
		if err != nil {
			v, err = strconv.ParseFloat(res, 64)
		}
		if err != nil {
			v = res
			//return nil, fmt.Errorf("Invalid env literal: %q", value)
		}

		return v, nil
	}
}
