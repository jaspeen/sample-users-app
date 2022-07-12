package services

import (
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
)

var Int2StringConverter = copier.TypeConverter{
	SrcType: copier.Int,
	DstType: copier.String,
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(int)

		if !ok {
			return nil, errors.New("src type not matching")
		}

		return strconv.Itoa(s), nil
	},
}

var String2IntConverter = copier.TypeConverter{
	SrcType: copier.String,
	DstType: copier.Int,
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(string)

		if !ok {
			return nil, errors.New("src type not matching")
		}

		return strconv.Atoi(s)
	},
}
