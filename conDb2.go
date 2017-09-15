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
	sql := "select * from TBL_USER where id = ? and name = ?" //TBL_DB2_CLR"// tbl_user"
	row := st.QueryRow(sql,args...)

	id := ""
	na := ""
	err := row.Scan(&id, &na)
	if err != nil {
		panic(err)
	}
	fmt.Println(id, na)
	//fmt.Println(utf8.Valid([]byte(na)))
	//dec := mahonia.NewDecoder("gbk")
	//na = dec.ConvertString(na)

	//fmt.Println(id, na)
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
	conStr := "DATABASE=rcbank;  HOSTNAME=192.168.20.12; PORT=56000; PROTOCOL=TCPIP;  UID=db2inst1; PWD=db2inst1;"

	if false {

		db, err := sql.Open("db2-cli", conStr)
		if err != nil {
			fmt.Println("open->", err)
			return
		}
		defer db.Close()

		st, err := db.Begin()
		if err != nil {
			fmt.Println("Begin->", err)
			return
		}
		defer st.Commit()
		SETSchema(st, "gft")
		//Query(st)
		//Update(st)
		QueryRow(st, "1" ,"东邪西毒")
	}

	if true {
		err := godb2.InitModel()
		if err != nil {
			fmt.Println("InitModel->", err)
			return
		}
		fmt.Println("--------1---------")
		engine := godb2.GetInstance()
		engine.SETSchema("gft")

		tb := &modules.Tbl_user{}
		//engine.Where("id = ? and name = ?","1","东邪西毒").Get(tb)
		ok, err := engine.Where("id = ? ","4").Get(tb)
		if err != nil {
			fmt.Println(err)
		}
		if !ok {
			fmt.Println("not find")
		}
		tb1 := &modules.Tbl_user{}
		ok, err  = engine.FindOne(tb1,"id = ? ","1")
		if err != nil {
			fmt.Println(err)
		}
		if !ok {
			fmt.Println("not find")
		}
		fmt.Printf("%+v\n",tb)
		fmt.Printf("%+v\n",tb1)
		//fmt.Println(len(tb.ID))
	}
	//stmt, err := st.Prepare("insert into tbl_user (ID, NAME) values (?, ?)")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer stmt.Close()
	//re,err:= stmt.Exec(6,`黄蓉`)
	//if err != nil {
	//	fmt.Println("1->",err)
	//}
	//fmt.Println(re.RowsAffected())

	if true {
		//QueryRow(st)
	}
	//st.Commit()
	if false {
		//SETSchema(st,"gft")
		//GetCurrentSchema(st)
		//CreateSchema(db)
		//createTable(db)
		//QueryRow(db)
		//Insert(db)
		//Delete(st)

	}

	if false {
		//MaxConnect(db)
	}

	time.Sleep(2 * time.Second)

	//db.SetMaxIdleConns(2)
	//sql := "select * from TBL_DB2_CLR"
	//rows, err := db.Query(sql)
	//if err != nil {
	//	fmt.Println("Query->",err)
	//	return
	//}
	//defer rows.Close()
	//rs,_:=rows.Columns()
	//fmt.Println(rs)
	//
	//for rows.Next() {
	//	na := ""
	//	id := ""
	//	err = rows.Scan(&id,&na)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(utf8.FullRuneInString(na))
	//	dec := mahonia.NewDecoder("gbk")
	//	na = dec.ConvertString(na)
	//
	//	fmt.Println(id, na)
	//}
	//
	//if err := dbOperations(); err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//}
	//var i ,n string
	//row:= db.QueryRow(sql)
	//row.Scan(&i ,&n )
	//fmt.Println(i,n)

	//insql:="insert into TBL_DB2_CLR values('3','西毒')"
	//res, err := db.Exec(insql)
	//fmt.Println(res)
	//if err != nil {
	//	fmt.Println("Query->",err)
	//	return
	//}
	//rows, err := db.Query(sql)
	//if err != nil {
	//	fmt.Println("Query->",err)
	//	return
	//}
	//for rows.Next() {
	//	na := ""
	//	id := ""
	//	err = rows.Scan(&id,&na)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(utf8.FullRuneInString(na))
	//	//dec := mahonia.NewDecoder("gbk")
	//	//na = dec.ConvertString(na)
	//
	//	fmt.Println(id, na)
	//}
}

func MaxConnect(db *sql.DB) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()
	//db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	for i := 0; i < 17; i++ {
		st, err := db.Begin()
		if err != nil {
			fmt.Println("Begin->", err)
			return
		}
		fmt.Println(i)
		if i <= 2 {
			st.Rollback()
		}
	}
}
