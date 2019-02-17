package csv

import (
	"reflect"
	"strings"
)

const tagKey = "csv"

func getTags(field reflect.StructField) tags {

	tagList := strings.Split(field.Tag.Get(tagKey), ",")
	
	return tags{
		index: field.Index,
		Name: getOrDefault(tagList, 0, field.Name),
		Unmarshal: getOrDefault(tagList, 1, ""),
		Marshal: getOrDefault(tagList, 2, ""),
		IsOptional: getOrDefault(tagList, 3, "") != "required",
	}
}


type tags struct {
	index []int
	Name string
	Unmarshal string
	Marshal string
	IsOptional bool
}

func getOrDefault(tags []string, i int, def string) (string) {
	if len(tags) > i && tags[i] != "" {
		return strings.Trim(tags[i], " ")
	}

	return def
}