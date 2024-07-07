/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 11:52
 * @Description:
 */

package utils

import (
	"github.com/leafney/rose"
	"github.com/leafney/whisky/pkgs/cmds"
	"testing"
)

func TestGetTypeName(t *testing.T) {
	a := rose.Md5HashBuf(cmds.ScriptTest)
	t.Logf("name %v", a)
}
