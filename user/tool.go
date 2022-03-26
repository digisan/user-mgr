package user

import (
	"reflect"
	"strings"

	. "github.com/digisan/go-generics/v2"
)

func ListField(objs ...any) (fields []string) {
	for _, obj := range objs {
		typ := reflect.TypeOf(obj)
		// fmt.Println("Type:", typ.Name(), "Kind:", typ.Kind())
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fields = append(fields, field.Name)
		}
	}
	return
}

func FieldValue(ptr any, field string) any {
	r := reflect.ValueOf(ptr).Elem()
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

func ListValidator(objs ...any) (tags []string) {
	for _, obj := range objs {
		typ := reflect.TypeOf(obj)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("validate")
			// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
			tags = append(tags, strings.Split(tag, ",")...)
		}
	}
	tags = Settify(tags...)
	return Filter(&tags,
		func(i int, e string) bool {
			return len(e) > 0 && NotIn(e, "required", "email") // exclude internal validate tags
		},
	)
}
