package user

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	v2 "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	ur "github.com/digisan/user-mgr/user/registered"
	. "github.com/digisan/user-mgr/util"
	"gopkg.in/go-playground/validator.v9"
)

type ValidateResult struct {
	OK  bool
	Err error
}

func NewValidateResult(ok bool, failMsg string) ValidateResult {
	if ok {
		return ValidateResult{ok, nil}
	}
	return ValidateResult{ok, fmt.Errorf("%v", failMsg)}
}

var (
	mFieldValidator = &sync.Map{}
)

func RegisterValidator(tag string, f func(o, v any) ValidateResult) {
	mFieldValidator.Store(tag, f)
}

func fnValidator(tag string) func(o, v any) ValidateResult {
	f, ok := mFieldValidator.Load(tag)
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("missing [%s] validator", tag))
	return f.(func(o, v any) ValidateResult)
}

func Validate(user *ur.User, exclTags ...string) error {
	vTags := ListValidator(user.Core, user.Profile, user.Admin)
	mIfFail := make(map[string]error)
	v := validator.New()
	for _, vTag := range vTags {
		if v2.In(vTag, exclTags...) {
			v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool { return true })
			continue
		}
		fn := fnValidator(vTag) // [fn] must be valued here !
		tag := vTag             // [tag] must be out of callback
		err := v.RegisterValidation(vTag, func(fl validator.FieldLevel) bool {
			rst := fn(user, fl.Field().Interface())
			mIfFail[tag] = rst.Err
			return rst.OK
		})
		lk.FailOnErr("%v", err)
	}
	err := v.Struct(user)
	if err != nil {
		_, tag := ErrField(err)
		for _, e := range err.(validator.ValidationErrors) {
			if err, ok := mIfFail[tag]; ok && err != nil {
				return err
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
	fieldTag := r.FindAllString(es, -1)
	field, tag := fieldTag[0], fieldTag[1]
	field = strings.Trim(field, "'")
	tag = strings.Trim(tag, "'")
	return field, tag
}
