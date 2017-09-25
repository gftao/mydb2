package main

import (
	"flag"
	"fmt"
	"os"
	"golib/modules/config"
	"ConDB2/godb2"
)

var confFile = flag.String("confFile", "./etc/ConDB2.ini", "配置文件")

func main() {
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
	//row:= tx.QueryRow("VALUES NEXTVAL FOR SEQ_IND")
	row:= tx.QueryRow("VALUES PREVVAL FOR SEQ_IND")
	v:=""

	row.Scan(&v)
	fmt.Println("-->",v)
	//tb := &modules.TBL_MCHT_INF{}
	//ok, err := tx.Where("", "").Get(tb)
	//if err != nil {
	//	fmt.Println("Where->", err)
	//	return
	//}
	//if ok {
	//	fmt.Println("not find")
	//	tx.Rollback()
	//	return
	//}
}
