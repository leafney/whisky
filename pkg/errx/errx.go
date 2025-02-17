/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 16:28
 * @Description:
 */

package errx

import (
	"errors"
	"fmt"
	"github.com/leafney/whisky/pkg/errc"
)

const defCode = errc.Failed

type XError struct {
	Code int
	Msg  string
}

func (e *XError) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s", e.Code, e.Msg)
}

func ErrorCM(code int, msg string) error {
	return &XError{
		Code: code,
		Msg:  msg,
	}
}

func ErrorCF(code int, format string, a ...interface{}) error {
	return &XError{
		Code: code,
		Msg:  fmt.Sprintf(format, a...),
	}
}

func ErrorCE(code int, err error) error {
	return &XError{
		Code: code,
		Msg:  err.Error(),
	}
}

func ErrorE(err error) error {
	return &XError{
		Code: defCode,
		Msg:  err.Error(),
	}
}

func ErrorM(msg string) error {
	return &XError{
		Code: defCode,
		Msg:  msg,
	}
}

func ErrorMF(format string, a ...interface{}) error {
	return &XError{
		Code: defCode,
		Msg:  fmt.Sprintf(format, a...),
	}
}

func GetError(err error) (int, string) {
	var s *XError
	if errors.As(err, &s) {
		return s.Code, s.Msg
	}

	return defCode, err.Error()
}

func GetCode(err error) int {
	var s *XError
	if errors.As(err, &s) {
		return s.Code
	}
	//// 不能解析的错误，返回默认 code
	return defCode
}

func GetMsg(err error) string {
	var s *XError
	if errors.As(err, &s) {
		return s.Msg
	}
	return err.Error()
}
