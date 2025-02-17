/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 16:28
 * @Description:
 */

package errc

import "strings"

const (
	LangDefault = LangZh
	LangZh      = "zh"
	LangEn      = "en"
)

// 通用
const (
	Success   = 200  // 成功
	Failed    = 3000 // 操作错误
	ErrClient = 4000 // 客户端错误
	ErrServer = 5000 // 服务端错误
)

const (
	ErrUnAuthorized = 4001 // 鉴权失败，token无效 invalid token 需要通过refresh_token重新获取access_token
	ErrAuthExpired  = 4002 // "授权已过期"
	ErrForbidden    = 4003 // "权限不足，请重新登录"
)

var ErrMsgZh = map[int]string{}

var ErrMsgEn = map[int]string{}

// 根据统一错误码获取相应的错误提示信息
func GetErrMsg(code int, lang string) string {
	var errList = ErrMsgEn
	switch lang {
	case LangZh:
		errList = ErrMsgZh
	case LangEn:
		errList = ErrMsgEn
	default:
		// 判断默认语言
		if strings.EqualFold(LangDefault, LangZh) {
			errList = ErrMsgZh
		} else {
			errList = ErrMsgEn
		}
	}

	msg, ok := errList[code]
	if ok {
		return msg
	}
	return errList[ErrServer]
}
