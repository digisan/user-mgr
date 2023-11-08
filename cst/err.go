package cst

const (
	INVALID_PARAM = iota
)

type Code int

var (
	m = map[Code]string{
		INVALID_PARAM: "invalid parameter",
	}
)

type UserMgrErr struct {
	message string
	code    Code
}

func (e *UserMgrErr) Error() string {
	if err, ok := m[e.code]; ok {
		return err
	}
	return "unknown user-mgr error message"
}

func NewErr(code Code) UserMgrErr {
	return UserMgrErr{
		code:    code,
		message: m[code],
	}
}
