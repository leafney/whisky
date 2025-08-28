/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-22
 * @Description: 统一错误处理
 */

package errors

import (
	"fmt"
)

// ErrorCode 错误码定义
type ErrorCode string

const (
	// 系统错误
	SystemError   ErrorCode = "SYSTEM_ERROR"
	InvalidInput  ErrorCode = "INVALID_INPUT"
	NotFoundError ErrorCode = "NOT_FOUND"
	Unauthorized  ErrorCode = "UNAUTHORIZED"
	Forbidden     ErrorCode = "FORBIDDEN"

	// 业务错误
	ConfigError   ErrorCode = "CONFIG_ERROR"
	NetworkError  ErrorCode = "NETWORK_ERROR"
	MonitorError  ErrorCode = "MONITOR_ERROR"
	CronError     ErrorCode = "CRON_ERROR"
	DatabaseError ErrorCode = "DATABASE_ERROR"
)

// BizError 业务错误
type BizError struct {
	Code    ErrorCode
	Message string
	Details interface{}
	Cause   error
}

func (e *BizError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewBizError 创建业务错误
func NewBizError(code ErrorCode, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// NewBizErrorWithDetails 创建带详细信息的业务错误
func NewBizErrorWithDetails(code ErrorCode, message string, details interface{}) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewBizErrorWithCause 创建带原因的业务错误
func NewBizErrorWithCause(code ErrorCode, message string, cause error) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Wrap 包装错误
func Wrap(err error, code ErrorCode, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// IsBizError 检查是否为业务错误
func IsBizError(err error) bool {
	_, ok := err.(*BizError)
	return ok
}

// GetCode 获取错误码
func GetCode(err error) ErrorCode {
	if bizErr, ok := err.(*BizError); ok {
		return bizErr.Code
	}
	return SystemError
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	if bizErr, ok := err.(*BizError); ok {
		return bizErr.Message
	}
	return err.Error()
}

// 预定义错误创建函数
func SystemErr(message string) *BizError {
	return NewBizError(SystemError, message)
}

func InvalidInputErr(message string) *BizError {
	return NewBizError(InvalidInput, message)
}

func NotFoundErr(message string) *BizError {
	return NewBizError(NotFoundError, message)
}

func ConfigErr(message string) *BizError {
	return NewBizError(ConfigError, message)
}

func NetworkErr(message string) *BizError {
	return NewBizError(NetworkError, message)
}

func MonitorErr(message string) *BizError {
	return NewBizError(MonitorError, message)
}

func CronErr(message string) *BizError {
	return NewBizError(CronError, message)
}

func DatabaseErr(message string) *BizError {
	return NewBizError(DatabaseError, message)
}
