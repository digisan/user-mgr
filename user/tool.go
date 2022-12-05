package user

import (
	. "github.com/digisan/go-generics/v2"
)

func ListField(objs ...any) (fields []string) {
	for _, obj := range objs {
		fields = append(fields, Fields(obj)...)
	}
	return
}

func ListValidator(objs ...any) (tags []string) {
	for _, obj := range objs {
		tags = append(tags, ValidatorTags(obj, "required", "email")...)
	}
	return Settify(tags...)
}
