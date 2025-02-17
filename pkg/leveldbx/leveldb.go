/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 18:28
 * @Description:
 */

package leveldbx

import (
	"github.com/leafney/rose"
	rleveldb "github.com/leafney/rose-leveldb"
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/pkg/xlogx"
)

type LevelDBSvc struct {
	*rleveldb.LevelDB
}

func NewLevelDBSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *LevelDBSvc {

	dbPath := cfg.LevelDB.Path

	if err := rose.DirExistsEnsure(dbPath); err != nil {
		log.Fatalf("[Leveldb] dbPath exist [%v] error [%v]", dbPath, err)
	}

	db, err := rleveldb.NewLevelDB(dbPath)
	if err != nil {
		log.Fatalf("[Leveldb] OpenFile [%v] error [%v]", dbPath, err)
	}

	go func() {
		// 等待停止信号
		<-stop
		if err := db.Close(); err != nil {
			log.Errorf("[Leveldb] Closed error [%v]", err)
		} else {
			log.Infoln("[Leveldb] Exit successful")
		}
	}()

	log.Infoln("[Leveldb] Load successful")

	return &LevelDBSvc{db}
}
