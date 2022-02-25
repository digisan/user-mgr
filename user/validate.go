package user

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/digisan/go-generics/str"
	lk "github.com/digisan/logkit"
	"gopkg.in/go-playground/validator.v9"
)

var (
	vTags           = ListUserValidator()
	mFieldValidator = &sync.Map{}
)

func RegisterValidator(tag string, f func(fv interface{}) bool) {
	mFieldValidator.Store(tag, f)
}

func fnValidator(tag string) func(fv interface{}) bool {
	f, ok := mFieldValidator.Load(tag)
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("missing [%s] validator", tag))
	return f.(func(fv interface{}) bool)
}

func (user *User) Validate(exclTags ...string) error {
	v := validator.New()
	for _, vTag := range vTags {
		if str.In(vTag, exclTags...) {
			v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool { return true })
			continue
		}
		fn := fnValidator(vTag) // [fn] must be valued here !
		v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool {
			return fn(fl.Field().String())
		})
	}
	if err := v.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%v", e)
		}
	}
	return nil
}

// return field string is one of user/valfield/ const
// from "Key: 'User.Addr' Error:Field validation for 'Addr' failed on the 'addr' tag"
func ErrField(err error) (string, string) {
	r := regexp.MustCompile(`'[^\.\s]+'`)
	es := fmt.Sprint(err)
	fieldtag := r.FindAllString(es, -1)
	field, tag := fieldtag[0], fieldtag[1]
	field = strings.Trim(field, "'")
	tag = strings.Trim(tag, "'")
	return field, tag
}
