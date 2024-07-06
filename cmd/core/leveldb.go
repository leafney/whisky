/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 12:36
 * @Description:
 */

package core

import (
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/whisky/global"
)

func InitLevelDB(stop chan struct{}) {
	dbPath := ".cache"
	//if rose.StrIsEmpty(dbPath) {
	//	dbPath = vars.DefLEVDBName
	//}

	//// 保证路径存在
	//if err := rose.DEnsurePathExist(dbPath); err != nil {
	//	global.GXLog.Fatalf("[Leveldb] dbPath exist [%v] error [%v]", dbPath, err)
	//}

	db, err := rleveldb.NewLevelDB(dbPath)
	if err != nil {
		global.GXLog.Fatalf("[Leveldb] OpenFile [%v] error [%v]", dbPath, err)
	}

	go func() {
		// 等待停止信号
		<-stop
		if err := db.Close(); err != nil {
			global.GXLog.Errorf("[Leveldb] Closed error [%v]", err)
		} else {
			global.GXLog.Infoln("[Leveldb] Exit successful")
		}
	}()

	global.GLevelDB = db

	global.GXLog.Infoln("[Leveldb] Load successful")
}
