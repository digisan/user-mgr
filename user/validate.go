package user

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"gopkg.in/go-playground/validator.v9"
)

type ValRst struct {
	OK      bool
	FailErr error
}

func NewValRst(ok bool, failMsg string) ValRst {
	if ok {
		return ValRst{ok, nil}
	}
	return ValRst{ok, fmt.Errorf("%v", failMsg)}
}

var (
	vTags           = ListValidator(User{}.Core, User{}.Profile, User{}.Admin)
	mFieldValidator = &sync.Map{}
)

func RegisterValidator(tag string, f func(o, v any) ValRst) {
	mFieldValidator.Store(tag, f)
}

func fnValidator(tag string) func(o, v any) ValRst {
	f, ok := mFieldValidator.Load(tag)
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("missing [%s] validator", tag))
	return f.(func(o, v any) ValRst)
}

func (user *User) Validate(exclTags ...string) error {
	mIfFail := make(map[string]error)
	v := validator.New()
	for _, vTag := range vTags {
		if In(vTag, exclTags...) {
			v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool { return true })
			continue
		}
		fn := fnValidator(vTag) // [fn] must be valued here !
		tag := vTag             // [tag] must be out of callback
		err := v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool {
			rst := fn(user, fl.Field().Interface())
			mIfFail[tag] = rst.FailErr
			return rst.OK
		})
		lk.FailOnErr("%v", err)
	}
	err := v.Struct(user)
	if err != nil {
		_, tag := ErrField(err)
		for _, e := range err.(validator.ValidationErrors) {
			if mIfFail[tag] != nil {
				return mIfFail[tag]
			}
			return fmt.Errorf("%v", e)
		}
	}
	lk.FailOnErr("%v", err)
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
