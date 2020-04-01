package errors

import "time"

// Error ...
type Error struct {
	when   time.Time // 错误发生时间
	code   int32     // 错误码
	what   string    // 错误原因
	detail string    // 错误详细信息
}

// New ...
func New(code int32, what, detail string) *Error {
	return &Error{
		when:   time.Now(),
		code:   code,
		what:   what,
		detail: detail,
	}
}

// NewError ...
func NewError(err Error) *Error {
	err.SetWhen(time.Now())
	return &err
}

// Error ...
func (e *Error) Error() string {
	return e.what
}

// SetWhen ...
func (e *Error) SetWhen(when time.Time) {
	e.when = when
}

// When ...
func (e *Error) When() time.Time {
	return e.when
}

// WhenString ...
func (e *Error) WhenString() string {
	return e.when.Format("2016-01-02 15:04:05.000000")
}

// Detail ...
func (e *Error) Detail() string {
	return e.detail
}

// What ...
func (e *Error) What() string {
	return e.what
}

// Code ...
func (e *Error) Code() int32 {
	return e.code
}
