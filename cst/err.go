package cst

import (
	"errors"
	"fmt"

	lk "github.com/digisan/logkit"
)

const (
	ERR_ON_FALSE = iota
	ERR_ON_TRUE
	ERR_3RD_LIB
	ERR_DB_NOT_INIT
	ERR_SEND_EMAIL
	ERR_TIMEOUT
	ERR_INV_PARAM
	ERR_USER_NOT_REG
	ERR_USER_ALREADY_REG
	ERR_USER_NOT_EXISTS
	ERR_USER_DORMANT
	ERR_USER_PWD_INCORRECT
	ERR_USER_INV_FIELD
	ERR_VCODE_MISSING
	ERR_VCODE_EXP
	ERR_VCODE_VERIFY_FAIL
	ERR_TOKEN_MISSING
	ERR_VALIDATOR_MISSING
	ERR_INV_DATA_FMT
	ERR_INPUT_INCONSISTENT
	ERR_MARSHAL_UNMARSHAL
	ERR_TYPE_CVT
	ERR_UNKNOWN
	ERR_COUNT
)

type Code int

var (
	mCodeErr = map[Code]error{
		ERR_ON_FALSE:           errors.New("on false condition, raise an error"),
		ERR_ON_TRUE:            errors.New("on true condition, raise an error"),
		ERR_3RD_LIB:            errors.New("3rd party library error"),
		ERR_DB_NOT_INIT:        errors.New("database is not initialized"),
		ERR_SEND_EMAIL:         errors.New("send email error"),
		ERR_TIMEOUT:            errors.New("timeout"),
		ERR_INV_PARAM:          errors.New("invalid parameter"),
		ERR_USER_NOT_REG:       errors.New("user is not registered"),
		ERR_USER_ALREADY_REG:   errors.New("user is already registered"),
		ERR_USER_NOT_EXISTS:    errors.New("user doesn't exist"),
		ERR_USER_DORMANT:       errors.New("user is dormant"),
		ERR_USER_PWD_INCORRECT: errors.New("user password is incorrect"),
		ERR_USER_INV_FIELD:     errors.New("user's field is invalid"),
		ERR_VCODE_MISSING:      errors.New("verification code is missing"),
		ERR_VCODE_EXP:          errors.New("verification code is expired"),
		ERR_VCODE_VERIFY_FAIL:  errors.New("verification code cannot be verified"),
		ERR_TOKEN_MISSING:      errors.New("token is missing"),
		ERR_VALIDATOR_MISSING:  errors.New("validator is missing"),
		ERR_INV_DATA_FMT:       errors.New("invalid data format"),
		ERR_INPUT_INCONSISTENT: errors.New("input values or content are inconsistent"),
		ERR_MARSHAL_UNMARSHAL:  errors.New("marshal or unmarshal error"),
		ERR_TYPE_CVT:           errors.New("data type conversion error"),
		ERR_UNKNOWN:            errors.New("unknown error"),
	}
)

type UserMgrErr struct {
	message string
	code    Code
}

func (e UserMgrErr) Error() string {
	if err, ok := mCodeErr[e.code]; ok {
		return err.Error()
	}
	return mCodeErr[ERR_UNKNOWN].Error()
}

func (e UserMgrErr) Code() Code {
	return e.code
}

func (e UserMgrErr) Wrap(err any) error {
	if errBase, ok := mCodeErr[e.code]; ok {
		return fmt.Errorf("%w: %v", errBase, err)
	}
	return fmt.Errorf("%w: %v", mCodeErr[ERR_UNKNOWN], err)
}

func Err(code Code) UserMgrErr {
	return UserMgrErr{
		code:    code,
		message: mCodeErr[code].Error(),
	}
}

func GetErrCode(err error) Code {
	if err != nil {
		for code, errBase := range mCodeErr {
			if err.Error() == errBase.Error() || errors.Is(err, errBase) {
				return code
			}
		}
	}
	return ERR_UNKNOWN
}

func GetErrMsg(code Code) string {
	if err, ok := mCodeErr[code]; ok {
		return err.Error()
	}
	return mCodeErr[ERR_UNKNOWN].Error()
}

func init() {
	lk.FailOnErrWhen(ERR_COUNT != len(mCodeErr), "%v", errors.New("mCodeMsg description missing"))
}
