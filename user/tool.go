package user

import (
	. "github.com/digisan/go-generics/v2"
	st "github.com/digisan/gotk/struct-tool"
)

func ListField(objs ...any) (fields []string) {
	for _, obj := range objs {
		fields = append(fields, st.Fields(obj)...)
	}
	return
}

func ListValidator(objs ...any) (tags []string) {
	for _, obj := range objs {
		tags = append(tags, st.ValidatorTags(obj, "required", "email")...)
	}
	return Settify(tags...)
}
