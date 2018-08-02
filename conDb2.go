package main

import (
	_ "bitbucket.org/phiggins/db2cli"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"
	"ConDB2/godb2"
	"ConDB2/modules"
	"github.com/axgle/mahonia"

	"github.com/astaxie/beedb"
)

var (
	connStr = flag.String("conn", "", "connection string to use")
	repeat  = flag.Uint("repeat", 1, "number of times to repeat query")
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: %s [options]

%s connects to DB2 and executes a simple SQL statement a configurable
number of times.

Here is a sample connection string:

DATABASE=MYDBNAME; HOSTNAME=localhost; PORT=50000; PROTOCOL=TCPIP; UID=username; PWD=password;
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func execQuery(st *sql.Stmt) error {
	rows, err := st.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var t time.Time
		err = rows.Scan(&t)
		if err != nil {
			return err
		}
		fmt.Printf("Time: %v\n", t)
	}
	return rows.Err()
}

func dbOperations() error {
	db, err := sql.Open("db2-cli", *connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// Attention: If you have to go through DB2-Connect you have to terminate SQL-statements with ';'
	st, err := db.Prepare("select current timestamp from sysibm.sysdummy1;")
	if err != nil {
		return err
	}
	defer st.Close()

	for i := 0; i < int(*repeat); i++ {
		err = execQuery(st)
		if err != nil {
			return err
		}
	}

	return nil
}

//创建Schema
func CreateSchema(st *sql.DB) {
	Sql := `create schema gft authorization db2inst1 `
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.Exec-->", err)
	}
	fmt.Println("create schema Success!")
}

//设置Schema
func SETSchema(st *sql.Tx, schema string) {
	Sql := `set current  schema  = ` + schema
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.Exec-->", err)
	}
	fmt.Println("current schema is ", schema)
}
func GetCurrentSchema(st *sql.Tx) (sc string) {
	Sql := `select  current  schema from sysibm.dual`
	row := st.QueryRow(Sql)
	row.Scan(&sc)
	fmt.Println(sc)
	return
}

//创建表
func createTable(st *sql.DB) {
	creatSql := `create  table tbl_user(ID CHARACTER (10)  NOT NULL,
		NAME CHARACTER (20),
		primary key (ID)
	)`
	_, err := st.Exec(creatSql)
	if err != nil {
		fmt.Println("st.creatSql-->", err)
	}
	fmt.Println("create Success!")
}

func Insert(st *sql.Tx) {
	Sql := `insert into tbl_user(ID, NAME) values('5', '国境')`
	//Sql := `insert into tbl_user(ID, NAME) values('1','东邪西毒'),('2','南拳北腿')`
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.insert-->", err)
	}
	//st.Commit()
}

func Delete(st *sql.Tx) {
	//Sql := `insert into tbl_user(ID, NAME) values('3','国境')`
	Sql := `delete from  tbl_user where id = '3'`
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.DELETE-->", err)
	}
	//st.Commit()
}

func Update(st *sql.Tx) {
	Sql := `update tbl_user set NAME ='郭靖' where id='2'`
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.Update-->", err)
	}
	//st.Commit()
}

//查询
func Query(st *sql.Tx) {
	sql := "select * from TBL_USER" //TBL_DB2_CLR"// tbl_user"

	rows, err := st.Query(sql)
	if err != nil {
		fmt.Println("st.Query-->", err)
	}

	for rows.Next() {
		id := ""
		na := ""
		err = rows.Scan(&id, &na)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, na)
		//fmt.Println(utf8.Valid([]byte(na)))
		//dec := mahonia.NewDecoder("gbk")
		//na = dec.ConvertString(na)

		//fmt.Println(id, na)
	}
}

//查询
func QueryRow(st *sql.Tx, args ...interface{}) {
	//sql := "select * from TBL_USER where id = ? and name = ?" //TBL_DB2_CLR"// tbl_user"
	//sql := "select * from TBL_MCHT_INF where MCHT_CD =  ?"//'881350141310000'"
	sql := "select * from TBL_TXN_STLM_CFG where INS_ID_CD =  ?" //'99991002   '"
	row := st.QueryRow(sql, args...)

	tb := &modules.TBL_TXN_STLM_CFG{}

	err := row.Scan(&tb.INS_ID_CD, &tb.CARD_TP, &tb.TXN_NUM, &tb.BUS_CD, &tb.STLM_FLG, &tb.STLM_DESC)
	if err != nil {
		panic(err)
	}
	fmt.Println(tb)
	//fmt.Println(utf8.Valid([]byte(na)))
	dec := mahonia.NewDecoder("gbk")
	na := dec.ConvertString(tb.STLM_DESC)

	fmt.Println(na)
}

//1208	N-1	UTF-8 编码的
//1386	D-4	GBK
func SetCodePage(st *sql.DB, encode string) {
	Sql := "db2set db2codepage=" + encode
	_, err := st.Exec(Sql)
	if err != nil {
		fmt.Println("st.db2codepage-->", err)
	}

	_, err = st.Exec("terminate")
	if err != nil {
		fmt.Println("st.db2codepage-->", err)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	//conStr := `Driver={IBM DB2 ODBC Driver};Hostname=localhost;Port=50000;Protocol=TCPIP;Database=center;CurrentSchema=GUFT;UID=guft;PWD=gft;`
	//conStr := "DATABASE=center; HOSTNAME=192.168.127.21; PORT=50000;PROTOCOL=TCPIP;CurrentSchema=GUFT;  UID=guft; PWD=gft;"
	//conStr := "DATABASE=rcbank;  HOSTNAME=192.168.20.74; PORT=56000; PROTOCOL=TCPIP;CurrentSchema=TEST; UID=db2inst1; PWD=db2inst1;"
	//conStr := "DATABASE=rcbank;  HOSTNAME=192.168.20.12; PORT=56000; PROTOCOL=TCPIP;  UID=db2inst1; PWD=db2inst1;"
	conStr := "DATABASE=rcbank; HOSTNAME=192.168.20.78; PORT=56000; PROTOCOL=TCPIP; CurrentSchema=apsmgm;  UID=apstfr; PWD=apstfr;"

	if true {
		db, err := sql.Open("db2-cli", conStr)
		if err != nil {
			fmt.Println("open->", err)
			return
		}
		defer db.Close()
		//st, err := db.Begin()
		//if err != nil {
		//	fmt.Println("Begin->", err)
		//	return
		//}
		st, err := db.Prepare("select MCHT_CD from TBL_MCHT_INF")
		if err != nil {
			fmt.Println("st.Query-->", err)
		}

		rows, err := st.Query()
		if err != nil {
			fmt.Println("st.Query-->", err)
		}
		defer rows.Close()
		for rows.Next() {
			mc :=""
			rows.Scan(&mc)
			fmt.Println(mc)
		}

		//count := 59
		//mc:= &modules.TBL_MCHT_INF{}
		//v := reflect.ValueOf(mc).Elem()
		//
		//for rows.Next() {
		//
		//	id := make([]interface{}, count)
		//	it := make([]interface{}, count)
		//	for j := 0; j < count; j++ {
		//		k := v.Field(j).Type()
		//		id[j] = reflect.New(k)
		//		it[j] = &id[j]
		//	}
		//	err = rows.Scan(it...)
		//
		//	if err != nil {
		//		fmt.Println("st1->Query-->", err)
		//	}
		//	fmt.Println(it)
		//	fmt.Println("st2Query-->", rows.Err())
		//}
		//fmt.Println("st3Query-->", rows.Err())
		//defer st.Commit()
		//rows, err := st.Query("select * from TBL_MCHT_INF where MCHT_CD = ?", "110170920645418")
		//rows, err := st.Query("select   COUNT(*)  from TBL_MCHT_INF")
		//if err != nil {
		//	fmt.Println("st.Query-->", err)
		//}
		//fmt.Println(rows.Columns())
		//cl,_ := rows.Columns()
		//l := len(cl)
		//t, err := rows.ColumnTypes()
		//if err != nil {
		//	fmt.Println("st.ColumnTypes-->", err)
		//}
		//for i := 0;i < l ; i++  {
		//	fmt.Println(t[i].Name(),t[i].DatabaseTypeName())
		//}
		//ok := rows.Next()
		//fmt.Println(ok)
		//c := 0
		//
		//for rows.Next() {
		//	rows.Scan(&c)
		//	fmt.Println(c)
		//}

		//c = 0
		////for ; it > 0; it-- {
		//	row := st.QueryRow("select * from TBL_MCHT_INF")
		//	mc := "q"
		//	row.Scan(&mc)
		//	fmt.Println("mc: ", mc)
		//	c ++
		////}

		//fmt.Println( rst.print())
		//MaxConnect(db)
		//SETSchema(st, "gft")
		//Query(st)
		//Update(st)
		//QueryRow(st, "99991002   ")
		//for {
		//	time.Sleep(5 * time.Second)
		//}
	}
	//查询
	if false {
		err := godb2.InitModel()
		if err != nil {
			fmt.Println("InitModel->", err)
			return
		}
		fmt.Println("--------find---------")
		e := godb2.GetInstance()
		defer e.Close()
		engine, err := e.Begin()
		//engine.SETSchema("gft")
		if err != nil {
			fmt.Println("Begin->", err)
			return
		}

		tb := &modules.TBL_TXN_STLM_CFG{}
		ok, err := engine.Where("INS_ID_CD = ? ", "99991002").Get(tb)
		if err != nil {
			fmt.Println(err)
			engine.Rollback()
			return
		}
		if !ok {
			fmt.Println("not find")
			engine.Rollback()
			return
		}
		fmt.Printf("%+v\n", tb)
		//tb1 := &modules.Tbl_mcht_inf{}
		//ok, err = engine.Where("MCHT_CD = ? ", "99991002").Get(tb1)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//if !ok {
		//	fmt.Println("not find")
		//}
		//fmt.Printf("%+v\n", tb1)
		tb2 := &modules.TBL_MCHT_BIZ_DEAL{}
		ok, err = engine.Where("MCHT_CD = ? and PROD_CD = ? and BIZ_CD = ? and TRANS_CD = ?", "999120241110000", "1151", "0000007", "1131").Get(tb2)
		if err != nil {
			fmt.Println(err)
			engine.Rollback()
			return
		}
		if !ok {
			fmt.Println("not find")
			engine.Rollback()
			return
		}
		fmt.Printf("%+v\n", tb2)
		//==========Insert=======
		if false {
			tb3 := &modules.TBL_MCHT_BIZ_DEAL{PROD_CD: "1151", BIZ_CD: "0000007", TRANS_CD: "1131", OPER_IN: "I", REC_UPD_OPR: "90"}
			tb3.MCHT_CD = "999120241110003"
			//tb3.REC_CRT_TS = time.Now()
			//tb3.REC_UPD_TS = time.Now()
			ok, err = engine.Insert(tb3)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !ok {
				fmt.Println("not find")
			}
			fmt.Printf("%+v\n", tb3)
		}
		//========================
		//==========Uptade=======
		if true {
			tb := &modules.TBL_MCHT_BIZ_DEAL{OPER_IN: "I", REC_OPR_ID: "3"}
			ok, err = engine.Where("MCHT_CD = ? and PROD_CD = ? and BIZ_CD = ? and TRANS_CD = ?", "999120241110003", "1151", "0000007", "1131").Uptade(tb)
			if err != nil {
				fmt.Println(err)
			}
			if !ok {
				fmt.Println("not find")
				engine.Rollback()
				return
			}

			fmt.Printf("%+v\n", tb)
		}
		//========================
		engine.Commit()
	}
	if false {
		fmt.Println("--------------------------------------------")
		db, err := sql.Open("db2-cli", conStr)
		if err != nil {
			fmt.Println("open->", err)
			return
		}
		defer db.Close()
		orm := beedb.New(db, "pg")
		beedb.OnDebug = true
		ormtb2 := &modules.TBL_TXN_STLM_CFG{}
		err = orm.Where("INS_ID_CD = ? ", "79991201   ").Find(ormtb2)
		//res,err:=orm.Exec(`INSERT INTO TBL_MCHT_BIZ_DEAL(MCHT_CD, PROD_CD, BIZ_CD, TRANS_CD, OPER_IN, REC_OPR_ID, REC_UPD_OPR, REC_CRT_TS, REC_UPD_TS) VALUES ('999120241110004', '1151', '0000007', '1131', 'I', '', '14402', '2016-11-18 10:28:04', '2017-03-03 17:22:12')`)
		if err != nil {
			fmt.Println("--", err)
		}
		fmt.Println(ormtb2)
		//fmt.Println(res.RowsAffected())

	}
	time.Sleep(2 * time.Second)
}

func MaxConnect(db *sql.DB) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	for i := 0; i < 10; i++ {
		_, err := db.Begin()
		if err != nil {
			fmt.Println("Begin->", err)
			return
		}
		fmt.Println(i)
		//st.Rollback()
	}
}
