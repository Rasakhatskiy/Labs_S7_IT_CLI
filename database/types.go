package database

import "reflect"

type TypeInteger struct {
}

type TypeReal struct {
}

type TypeChar struct {
}

type TypeString struct {
}

type TypeHTML struct {
}

type TypeStringRange struct {
}

var (
	TypeIntegerTS     = reflect.TypeOf(TypeInteger{}).String()
	TypeRealTS        = reflect.TypeOf(TypeReal{}).String()
	TypeCharTS        = reflect.TypeOf(TypeChar{}).String()
	TypeStringTS      = reflect.TypeOf(TypeString{}).String()
	TypeHTMLTS        = reflect.TypeOf(TypeHTML{}).String()
	TypeStringRangeTS = reflect.TypeOf(TypeStringRange{}).String()
)

var TypesListStr []string
