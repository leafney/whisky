/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 16:46
 * @Description:
 */

package parsex

type (
	PageParam struct {
		Page int `json:"page" form:"page" query:"page"`
		Size int `json:"size" form:"size" query:"size"`
	}

	PageInt64Param struct {
		Page int64 `json:"page" form:"page" query:"page"`
		Size int64 `json:"size" form:"size" query:"size"`
	}
)
