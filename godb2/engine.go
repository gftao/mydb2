package godb2

import (
	//"golib/modules/config"
)

var instance *EngineDb2

func InitModel() error {
	var err error
	instance, err = initDb()
	return err

}

func initDb() (*EngineDb2, error) {
	//dbconf := config.StringDefault("dbconf", "")
	//if dbconf == "" {
	//	return nil, errors.New("数据库配置文件为空"+dbconf)
	//}

	//db, err := opendb(dbconf)
	db, err := opendb("")
	if err != nil {
		return nil, err
	}
	return db, err
}