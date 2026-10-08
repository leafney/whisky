/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 16:46
 * @Description:
 */

package parsex

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

/*
注意：`validate` 在使用时，最好不要将零值作为required，否则会被认为是空值，导致校验失败。例如 `status=0` 时，即使status是required，也会被认为是0,未输入内容，导致校验失败。

*/

// ParseJson structs with tags `json` for json parameters
//
//	Example: {"name": "leafney", "age": 18}
func ParseJson(c *fiber.Ctx, obj interface{}) error {
	if err := c.BodyParser(obj); err != nil {
		return err
	}

	if v := validate.Struct(obj); !v.Validate() {
		return v.Errors
	}
	return nil
}

// ParseForm structs with tags `form` for form parameters
//
//	Example: name=leafney&age=18
func ParseForm(c *fiber.Ctx, obj interface{}) error {
	return ParseJson(c, obj)
}

// ParseParam structs with tags `params` for path parameters
//
//	Example: /users/:id
func ParseParam(c *fiber.Ctx, obj interface{}) error {
	if err := c.ParamsParser(obj); err != nil {
		return err
	}

	if v := validate.Struct(obj); !v.Validate() {
		return v.Errors
	}
	return nil
}

// ParseQuery structs with tags `query` for query parameters
//
//	Example: /users?id=1&name=leafney
func ParseQuery(c *fiber.Ctx, obj interface{}) error {
	if err := c.QueryParser(obj); err != nil {
		return err
	}

	if v := validate.Struct(obj); !v.Validate() {
		return v.Errors
	}
	return nil
}

// ParseAll structs with tags `params`, `query` and `json` for all parameters
func ParseAll(c *fiber.Ctx, obj interface{}) error {
	// 解析 URL 参数
	_ = c.ParamsParser(obj)

	// 解析 Query 参数
	_ = c.QueryParser(obj)

	// 解析 JSON 请求体
	_ = c.BodyParser(obj)

	// 校验参数
	if v := validate.Struct(obj); !v.Validate() {
		return v.Errors
	}

	return nil
}

// ParsePage 处理分页参数 page and size
//
//	Example: /users?page=1&size=10
func ParsePage(c *fiber.Ctx) (*PageParam, error) {
	items := new(PageParam)
	if err := c.QueryParser(items); err != nil {
		return nil, err
	}
	if items.Page < 1 {
		items.Page = PageMinDef
	}
	if items.Size < 1 {
		items.Size = PageSizeDef
	}
	if items.Size > 20 {
		items.Size = PageSizeMax
	}

	return items, nil
}

// ParseInt64Page 处理分页参数 page and size
func ParseInt64Page(c *fiber.Ctx) (*PageInt64Param, error) {
	items := new(PageInt64Param)
	if err := c.QueryParser(items); err != nil {
		return nil, err
	}
	if items.Page < 1 {
		items.Page = int64(PageMinDef)
	}
	if items.Size < 1 {
		items.Size = int64(PageSizeDef)
	}
	if items.Size > 20 {
		items.Size = int64(PageSizeMax)
	}

	return items, nil
}
