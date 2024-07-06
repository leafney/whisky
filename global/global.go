/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:34
 * @Description:
 */

package global

import (
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/rose/xlog"
)

var (
	GLevelDB *rleveldb.LevelDB
	GXLog    *xlog.Log
)
