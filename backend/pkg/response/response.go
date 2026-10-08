/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 13:05
 * @Description:
 */

package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/pkg/errc"
)

type Response struct {
	Code  int         `json:"code"`  // 自定义错误码
	Msg   string      `json:"msg"`   // 给用户看的错误信息
	Data  interface{} `json:"data"`  // 返回的数据
	Error string      `json:"error"` // 给开发者看的详细错误信息
}

type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func jsonResult(c *fiber.Ctx, statusCode int, errCode int, msg string, data interface{}, err string) error {
	return c.Status(statusCode).JSON(Response{
		Code:  errCode,
		Msg:   msg,
		Data:  data,
		Error: err,
	})
}

func Ok(c *fiber.Ctx) error {
	msg := errc.GetErrMsg(errc.Success, "")
	return jsonResult(c, fiber.StatusOK, errc.Success, msg, map[string]interface{}{}, "")
}

func OkWithData(c *fiber.Ctx, data interface{}) error {
	msg := errc.GetErrMsg(errc.Success, "")
	return jsonResult(c, fiber.StatusOK, errc.Success, msg, data, "")
}

func OkWithPage(c *fiber.Ctx, list interface{}, total int64) error {
	msg := errc.GetErrMsg(errc.Success, "")
	pageData := PageData{
		List:  list,
		Total: total,
	}
	return jsonResult(c, fiber.StatusOK, errc.Success, msg, pageData, "")
}

func OkWithMsg(c *fiber.Ctx, msg string) error {
	return jsonResult(c, fiber.StatusOK, errc.Success, msg, map[string]interface{}{}, "")
}

func OkWithDetailed(c *fiber.Ctx, msg string, data interface{}) error {
	return jsonResult(c, fiber.StatusOK, errc.Success, msg, data, "")
}

// *******************

// Fail 操作失败
func Fail(c *fiber.Ctx, err string) error {
	msg := errc.GetErrMsg(errc.Failed, "")
	return jsonResult(c, fiber.StatusOK, errc.Failed, msg, map[string]interface{}{}, err)
}

// FailWithCode 使用统一错误消息
func FailWithCode(c *fiber.Ctx, errCode int, err string) error {
	msg := errc.GetErrMsg(errCode, "")
	return jsonResult(c, fiber.StatusOK, errCode, msg, map[string]interface{}{}, err)
}

// FailWithMsg 自定义错误消息
func FailWithMsg(c *fiber.Ctx, msg, err string) error {
	return jsonResult(c, fiber.StatusOK, errc.Failed, msg, map[string]interface{}{}, err)
}

// FailWith400Msg client error
func FailWith400Msg(c *fiber.Ctx, err string) error {
	msg := errc.GetErrMsg(errc.ErrClient, "")
	return jsonResult(c, fiber.StatusBadRequest, errc.ErrClient, msg, map[string]interface{}{}, err)
}

// FailWith401Msg Unauthorized 认证相关
func FailWith401Msg(c *fiber.Ctx, err string) error {
	msg := errc.GetErrMsg(errc.ErrUnAuthorized, "")
	return jsonResult(c, fiber.StatusUnauthorized, errc.ErrUnAuthorized, msg, map[string]interface{}{}, err)
}

// FailWith403Msg Forbidden 权限相关
func FailWith403Msg(c *fiber.Ctx, err string) error {
	msg := errc.GetErrMsg(errc.ErrForbidden, "")
	return jsonResult(c, fiber.StatusUnauthorized, errc.ErrForbidden, msg, map[string]interface{}{}, err)
}

// FailWith500Msg server error
func FailWith500Msg(c *fiber.Ctx, err string) error {
	msg := errc.GetErrMsg(errc.ErrServer, "")
	return jsonResult(c, fiber.StatusInternalServerError, errc.ErrServer, msg, map[string]interface{}{}, err)
}
