/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-22
 * @Description: HTTP 错误处理
 */

package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/pkg/response"
)

// HandleHTTPError 处理HTTP错误
func HandleHTTPError(c *fiber.Ctx, err error) error {
	if bizErr, ok := err.(*BizError); ok {
		return response.Fail(c, bizErr.Message)
	}

	// 默认返回系统错误
	return response.Fail(c, "系统错误")
}

// HandleHTTPErrorWithCode 处理HTTP错误并返回错误码
func HandleHTTPErrorWithCode(c *fiber.Ctx, err error) error {
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return response.FailWith400Msg(c, bizErr.Message)
	}

	// 默认返回系统错误
	return response.FailWithCode(c, 500, "系统错误")
}

// LogAndHandleError 记录并处理错误
func LogAndHandleError(c *fiber.Ctx, err error, logger interface{ Errorf(string, ...interface{}) }) error {
	if bizErr, ok := err.(*BizError); ok {
		logger.Errorf("业务错误 [%s]: %s", bizErr.Code, bizErr.Message)
		if bizErr.Cause != nil {
			logger.Errorf("错误原因: %v", bizErr.Cause)
		}
		return response.Fail(c, bizErr.Message)
	}

	logger.Errorf("系统错误: %v", err)
	return response.Fail(c, "系统错误")
}
