package main

import (
	"flag"
	"fmt"
	"os"
	"golib/modules/config"
	"ConDB2/godb2"
	"golib/modules/logr"
	"ConDB2/modules"
)

var confFile = flag.String("confFile", "./etc/ConDB2.ini", "配置文件")

func main() {
	flag.Parse()

	flag.Parse()
	if *confFile == "" {
		fmt.Println("配置文件不存在", *confFile)
		flag.Usage()
		os.Exit(-1)
	}
	err := config.InitModuleByParams(*confFile)
	if err != nil {
		fmt.Println("加载配置文件失败", err)
		os.Exit(-1)
	}
	//初始化日志
	fmt.Println("开始初始化日志")
	err = logr.InitModules()
	if err != nil {
		fmt.Println("初始化日志失败", err)
		return
	}
	fmt.Println("开始初始化数据库")
	/////////db2
	err = godb2.InitModel()
	if err != nil {
		fmt.Println("InitModel->", err)
		return
	}
	dbConn := godb2.GetInstance()
	tx, err := dbConn.Begin()
	if err != nil {
		fmt.Println("获取DB事物失败:", err)
	}
	//ok, err := tx.Where("STATE = ?", "1").GetAll(&MCHT_SYNC)
	MCHT_INF := &modules.TBL_MCHT_INF{}
	ok, err := tx.Where("MCHT_CD = ? ", "").Get(MCHT_INF)
	if err != nil || !ok {
		tx.Rollback()
		fmt.Errorf("更新商户信息表失败:%v,%s", ok, err)
		return
	}
	tx.Commit()


}
