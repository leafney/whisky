/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 13:05
 * @Description:
 */

package response

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/global/respcode"
)

type Response struct {
	Code  int         `json:"code"`  // 自定义错误码
	Msg   string      `json:"msg"`   // 给用户看的错误信息
	Data  interface{} `json:"data"`  // 返回的数据
	Error string      `json:"error"` // 给开发者看的详细错误信息
}

func jsonResult(c fiber.Ctx, statusCode int, errCode int, msg string, data interface{}, err string) error {
	return c.Status(statusCode).JSON(Response{
		Code:  errCode,
		Msg:   msg,
		Data:  data,
		Error: err,
	})
}

func Ok(c fiber.Ctx) error {
	return jsonResult(c, fiber.StatusOK, respcode.Success, "操作成功", map[string]interface{}{}, "")
}

func OkWithData(c fiber.Ctx, data interface{}) error {
	return jsonResult(c, fiber.StatusOK, respcode.Success, "操作成功", data, "")
}

// Fail 操作失败
func Fail(c fiber.Ctx, err string) error {
	return jsonResult(c, fiber.StatusOK, respcode.Failed, "操作失败", map[string]interface{}{}, err)
}
