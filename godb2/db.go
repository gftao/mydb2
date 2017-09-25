package godb2

import (
	"database/sql"
	_ "bitbucket.org/phiggins/db2cli"
	"fmt"
	"github.com/pkg/errors"
	"github.com/widuu/goini"
	"reflect"
	"github.com/axgle/mahonia"
	"strconv"
	"strings"
	"sync"
	"golib/modules/logr"
	"time"
	"golib/modules/config"
)

type EngineDb2 struct {
	db *sql.DB
}

type EngStmt struct {
	sync.Locker
	Tx        *sql.Tx
	Query     interface{}
	Args      []interface{}
	tableName string
}

type levelmap map[int]map[string]map[string]string

func opendb(file string) (*EngineDb2, error) {
	talkChan := loadConfig(file)
	dbdata := make(map[string]string)

	dbdata["DATABASE"] = talkChan[0]["db"]["DATABASE"]
	dbdata["HOSTNAME"] = talkChan[0]["db"]["HOSTNAME"]
	dbdata["PORT"] = talkChan[0]["db"]["PORT"]
	dbdata["UID"] = talkChan[0]["db"]["UID"]
	dbdata["PWD"] = talkChan[0]["db"]["PWD"]
	dbdata["CurrentSchema"] = talkChan[0]["db"]["CurrentSchema"]
	dbdata["idlcon"] = talkChan[0]["db"]["idlcon"]
	dbdata["maxcon"] = talkChan[0]["db"]["maxcon"]
	conStr := "DATABASE=" + dbdata["DATABASE"] + "; HOSTNAME=" + dbdata["HOSTNAME"] + "; PORT=" + dbdata["PORT"] + "; PROTOCOL=TCPIP; " + "CurrentSchema=" + dbdata["CurrentSchema"] + "; UID=" + dbdata["UID"] + "; PWD=" + dbdata["PWD"] + ";"
	//conStr := "DATABASE=rcbank; HOSTNAME=192.168.20.78; PORT=56000; PROTOCOL=TCPIP; CurrentSchema=APSTFR;  UID=apstfr; PWD=apstfr;"
	fmt.Println(conStr)
	d, err := sql.Open("db2-cli", conStr)
	if err != nil {
		fmt.Println("open->", err)
		return nil, errors.New("open db2 failed:" + conStr)
	}
	idlCon, err := strconv.Atoi(dbdata["idlcon"])
	if err != nil || idlCon <= 0 {
		fmt.Println("空闲连接默认10")
		idlCon = 10
	}

	maxCon, err := strconv.Atoi(dbdata["maxcon"])
	if err != nil || maxCon <= 0 {
		fmt.Println("最大连接默认100")
		maxCon = 100
	}
	d.Ping()
	d.SetMaxIdleConns(idlCon)
	d.SetMaxOpenConns(maxCon)
	return &EngineDb2{db: d}, nil
}

func (tx *EngStmt) SETSchema(schema string) error {
	Sql := `set current  schema  = ` + schema
	_, err := tx.Tx.Exec(Sql)
	if err != nil {
		return err
	}
	return nil
}
func (tx *EngStmt) GetCurrentSchema() (sc string) {
	Sql := `select  current  schema from sysibm.dual`
	row := tx.Tx.QueryRow(Sql)
	row.Scan(&sc)
	return
}
func (tx *EngStmt) HasSchema(tblName string) (bool, []string) {
	nl := strings.Split(tblName, ".")
	if len(nl) == 2 {
		return true, nl
	}
	return false, nl
}

func (e *EngineDb2) Begin() (*EngStmt, error) {
	s, err := e.db.Begin()
	return &EngStmt{Tx: s}, err
}
func (e *EngineDb2) Close() error {
	return e.db.Close()
}
func (tx *EngStmt) Commit() error {
	return tx.Tx.Commit()
}

func (tx *EngStmt) Rollback() error {
	return tx.Tx.Rollback()
}

func (e *EngineDb2) Query(query interface{}, args ...interface{}) (*sql.Rows, error) {
	results, err := e.db.Query(query.(string))
	return results, err
}

func (e *EngineDb2) QueryRow(query interface{}, args ...interface{}) (*sql.Row) {
	results := e.db.QueryRow(query.(string))
	return results
}
func (tx *EngStmt) QueryRow(query interface{}, args ...interface{}) (*sql.Row) {
	results := tx.Tx.QueryRow(query.(string))
	return results
}
func (tx *EngStmt) Where(query interface{}, args ...interface{}) *EngStmt {
	tx.Query = query
	tx.Args = args
	return tx //&EngStmt{Query: query, Args: args}
}
func (tx *EngStmt) Get(bean interface{}) (bool, error) {
	v := reflect.ValueOf(bean).Elem()
	t := reflect.TypeOf(bean).Elem()
	if tb, ok := v.Interface().(TableName); ok {
		tx.tableName = tb.TableName()
	}
	qs := "select * from " + tx.tableName + " where " + tx.Query.(string)
	fmt.Println(qs)
	fmt.Println(tx.Args)
	rows, err := tx.Tx.Query(qs, tx.Args...)
	//fmt.Println("---------")
	if err != nil {
		//fmt.Println(err)
		return false, err
	}
	count := v.NumField()
	fmt.Println(count)

	if true {
		for rows.Next() {
			id := make([]interface{}, count)
			for j := 0; j < count; j++ {
				k := v.Field(j).Kind()
				switch k {
				case reflect.Int:
					id[j] = 0
				case reflect.String:
					id[j] = ""
				case reflect.Struct:
					id[j] = reflect.New(reflect.TypeOf(reflect.Struct))
				default:
					fmt.Println("init type->", t.Field(j).Name, k)
					id[j] = ""
				}
			}
			//fmt.Printf("id->%+v", id)
			it := make([]interface{}, count)
			for j := 0; j < count; j++ {
				it[j] = &id[j]
			}
			err = rows.Scan(it...)
			if err != nil {
				panic(err)
			}
			dec := mahonia.NewDecoder("gbk")
			for j := 0; j < count; j++ {
				k := reflect.ValueOf(id[j]).Kind()
				//fmt.Println("----String--1---", id[j])
				switch k {
				case reflect.Int, reflect.Int32:
					str := strconv.FormatInt(reflect.ValueOf(id[j]).Int(), 10)
					i, _ := strconv.Atoi(str)
					v.Field(j).Set(reflect.ValueOf(i))
				case reflect.String:
					v.Field(j).Set(reflect.ValueOf(id[j]))
				case reflect.Slice:
					data := reflect.ValueOf(id[j]).Interface().([]byte)
					str := string(data)
					r := dec.ConvertString(str)
					//fmt.Println("----String--2---", r)
					v.Field(j).Set(reflect.ValueOf(r))
				case reflect.Struct:
					v.Field(j).Set(reflect.ValueOf(id[j]))
				default:
					//reflect.Copy(v.Field(j),reflect.ValueOf(id[j]))
					fmt.Println("----default----", reflect.TypeOf(id[j]).Kind(), id[j])
				}
				//v.Field(j).Set(reflect.ValueOf(r))
			}
		}
	}
	return true, nil
}

func (s *EngStmt) FindOne(links interface{}, querystring string, args ...interface{}) (bool, error) {
	has, err := s.Where(querystring, args...).Get(links)
	return has, err
}

func (tx *EngStmt) Uptade(bean interface{}) (bool, error) {
	fmt.Println("--------------update---------------")

	v := reflect.ValueOf(bean).Elem()
	t := reflect.TypeOf(bean).Elem()

	tblName := ""
	if tb, ok := v.Interface().(TableName); ok {
		tblName = tb.TableName()
	}

	fmt.Println(tblName)
	qs := strings.Split(tx.Query.(string), "?")
	if (len(qs) - 1) != len(tx.Args) {
		return false, errors.New("参数不匹配")
	}
	for i, s := range tx.Args {
		qs[i] += "'" + reflect.ValueOf(s).String() + "'"
	}
	fmt.Println(qs)
	q := strings.Join(qs, "")
	///for update
	SqlUp := "select * from " + tblName + " where " + q + " for update"
	fmt.Println(SqlUp)
	if true {
		//Sql := `update tbl_user set NAME ='郭靖' where id='2'`
		_, err := tx.Tx.Exec(SqlUp)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
	}
	///
	count := v.NumField()
	ns := []string{}
	for i := 0; i < count; i++ {
		k := v.Field(i).Kind()
		switch k {
		case reflect.Struct:
			if t.Field(i).Name == "REC_UPD_TS" {
				vvv := "(VALUES TIMESTAMP(CURRENT TIMESTAMP))"
				s := fmt.Sprintf("%s = %s", t.Field(i).Name, vvv)
				ns = append(ns, s)
			}

		default:
			if v.Field(i).String() != "" {
				s := fmt.Sprintf("%s = '%s'", t.Field(i).Name, v.Field(i).String())
				ns = append(ns, s)
			}
		}
	}

	Sql := "update " + tblName + " set " + strings.Join(ns, ",") + " where " + q
	fmt.Println(Sql)
	//qs := "select * from " + tblName + " where " + tx.Query.(string)
	//update TBL_MCHT_BIZ_DEAL  set OPER_IN = 'U', REC_UPD_OPR = '1' where   mcht_cd = '999120241110001'
	if true {
		//Sql := `update tbl_user set NAME ='郭靖' where id='2'`
		_, err := tx.Tx.Exec(Sql)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
	}

	return true, nil
}

func (s *EngStmt) Insert(bean interface{}) (bool, error) {

	v := reflect.ValueOf(bean).Elem()
	t := reflect.TypeOf(bean).Elem()
	tblName := ""
	crtSchema := s.GetCurrentSchema()
	defer s.SETSchema(crtSchema)

	if tb, ok := v.Interface().(TableName); ok {
		tblName = tb.TableName()
	}
	//fmt.Println(tblName,crtSchema)
	if h, n := s.HasSchema(tblName); h {
		//defer s.Unlock()
		tblName = n[1]
		//s.Lock()
		s.SETSchema(n[0])
	}
	count := v.NumField()
	//fmt.Println(count)
	ns := []string{}
	nv := []string{}
	for i := 0; i < count; i++ {
		ns = append(ns, t.Field(i).Name)
		k := v.Field(i).Kind()
		switch k {
		case reflect.Int,reflect.Int64, reflect.Float64:
			vv := fmt.Sprintf("%v", v.Field(i).Interface())
			nv = append(nv, vv)
		case reflect.Struct:
			/*
			//fmt.Println(k, t.Field(i).Name, v.Field(i).Elem())
			vv := fmt.Sprintf("%v",v.Field(i))
			//fmt.Println(k, t.Field(i).Name, v.Field(i).String(),vv)
			//fmt.Println(strings.Split(vv," +")[0])
			vvv := strings.Split(vv," +")[0]
			nv = append(nv, "'"+vvv+"'")
			*/
			if t.Field(i).Name == "REC_UPD_TS" ||t.Field(i).Name == "REC_CRT_TS" {
				vvv := "(VALUES TIMESTAMP(CURRENT TIMESTAMP))"
				nv = append(nv, vvv)
			}else {
				vvv := "(select current date from sysibm.sysdummy1)"
				nv = append(nv, vvv)
			}

		default:
			nv = append(nv, "'"+v.Field(i).String()+"'")
		}

	}
	su := tblName + " (" + strings.Join(ns, ", ") + ")"
	sv := " VALUES (" + strings.Join(nv, ", ") + ")"
	//fmt.Println(su)
	//fmt.Println(nv)
	//fmt.Println(sv)

	Sql := `INSERT INTO ` + su + sv //+" VALUES (" + ")"
	fmt.Println(Sql)
	//Sql := `INSERT INTO tbl_user(ID, NAME) values('5', '国境')`
	if true {
		_, err := s.Tx.Exec(Sql)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
	}

	return true, nil
}

func GetInstance() *EngineDb2 {
	err := instance.db.Ping()

	dbconf := config.StringDefault("dbconf", "")
	for err != nil {
		logr.Error("数据库连接已经断开，重新连接", err)
		instance, err = opendb(dbconf)
		time.Sleep(time.Duration(5) * time.Second)
	}
	return instance
}

func loadConfig(file string) (talkChan levelmap) {
	talkChan = make(levelmap)
	conf := goini.SetConfig(file)
	talkChan1 := conf.ReadList()

	for k, v := range talkChan1 {
		talkChan[k] = v
		for k1, v1 := range v {
			talkChan[k][k1] = v1
			for k2, v2 := range v1 {
				talkChan[k][k1][k2] = v2
			}

		}
	}

	return talkChan
}
