package csv

import (
	"reflect"
	"strings"
)

const tagKey = "csv"

func getFieldInfo(field reflect.StructField) fieldInfo {

	tagList := strings.Split(field.Tag.Get(tagKey), ",")

	return fieldInfo{
		index:      field.Index,
		Name:       getOrDefault(tagList, 0, field.Name),
		Unmarshal:  getOrDefault(tagList, 1, ""),
		Marshal:    getOrDefault(tagList, 2, ""),
		IsOptional: getOrDefault(tagList, 3, "optional") != "required",
		Type:       field.Type,
	}
}

type fieldInfo struct {
	index      []int
	Name       string
	Unmarshal  string
	Marshal    string
	IsOptional bool
	Type       reflect.Type
}

func getOrDefault(tags []string, i int, def string) string {
	if len(tags) > i && tags[i] != "" {
		return strings.Trim(tags[i], " ")
	}

	return def
}
