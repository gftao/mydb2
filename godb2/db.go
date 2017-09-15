package godb2

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/widuu/goini"
	"reflect"
	"github.com/axgle/mahonia"
	"strconv"
)

type EngineDb2 struct {
	db *sql.DB
	//Cond
}
type EngStmt struct {
	Tx        *sql.Tx
	Query     interface{}
	Args      []interface{}
	tableName string
}

type levelmap map[int]map[string]map[string]string

func opendb(file string) (*EngineDb2, error) {
	//talkChan := loadConfig(file)
	//dbdata := make(map[string]string)
	//
	//dbdata["DATABASE"] = talkChan[0]["db"]["DATABASE"]
	//dbdata["HOSTNAME"] = talkChan[0]["db"]["HOSTNAME"]
	//dbdata["PORT"] = talkChan[0]["db"]["PORT"]
	//dbdata["UID"] = talkChan[0]["db"]["UID"]
	//dbdata["PWD"] = talkChan[0]["db"]["PWD"]
	//dbdata["idlcon"] = talkChan[0]["db"]["idlcon"]
	//dbdata["maxcon"] = talkChan[0]["db"]["maxcon"]
	//CurrentSchema=GUFT;
	conStr := "DATABASE=rcbank;  HOSTNAME=192.168.20.12; PORT=56000; PROTOCOL=TCPIP; UID=db2inst1; PWD=db2inst1;"

	d, err := sql.Open("db2-cli", conStr)
	if err != nil {
		fmt.Println("open->", err)
		return nil, errors.New("open db2 failed:" + conStr)
	}
	return &EngineDb2{db: d}, nil
}

func (e *EngineDb2) SETSchema(schema string) error {
	Sql := `set current  schema  = ` + schema
	_, err := e.db.Exec(Sql)
	if err != nil {
		return err
	}
	return nil
}
func (e *EngineDb2) Query(query interface{}, args ...interface{}) (*sql.Rows, error) {
	results, err := e.db.Query(query.(string))
	return results, err
}

func (e *EngineDb2) Where(query interface{}, args ...interface{}) *EngStmt {

	s, _ := e.db.Begin()

	return &EngStmt{Tx: s, Query: query, Args: args}
}
func (tx *EngStmt) Get(bean interface{}) (bool, error) {
	v := reflect.ValueOf(bean).Elem()
	if tb, ok := v.Interface().(TableName); ok {
		tx.tableName = tb.TableName()
	}
	qs := "select * from " + tx.tableName + " where " + tx.Query.(string)
	//fmt.Println(qs)
	//fmt.Println(tx.Args)
	rows, err := tx.Tx.Query(qs, tx.Args...)
	//fmt.Println("---------")
	if err != nil {
		//fmt.Println(err)
		return false, err
	}
	count := v.NumField()
	//fmt.Println(count)
	for rows.Next() {
		id := make([]interface{}, count)
		for j := 0; j < count; j++ {
			k := v.Field(j).Kind()
			switch k {
			case reflect.Int:
				id[j] = 0
			case reflect.String:
				id[j] = ""
			default:
				id[j] = ""
			}
		}
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
			case reflect.String, reflect.Slice:
				data := reflect.ValueOf(id[j]).Interface().([]byte)
				str := string(data)
				r := dec.ConvertString(str)
				//fmt.Println("----String--2---", r)
				v.Field(j).Set(reflect.ValueOf(r))
			default:
				//fmt.Println("----default----", reflect.ValueOf(id[j]).Kind(), id[j])
			}
			//v.Field(j).Set(reflect.ValueOf(r))
		}
	}
	return true, nil
}

func (s *EngineDb2) FindOne(links interface{}, querystring string, args ...interface{}) (bool, error) {
	has, err := s.Where(querystring, args...).Get(links)
	return has, err
}

func (s *EngStmt) Uptade(bean interface{}) (bool, error) {

	return true, nil
}

func GetInstance() *EngineDb2 {
	//err := instance.db.Ping()
	//for err != nil {
	//	fmt.Printf("-->",err)
	//}
	//dbconf := config.StringDefault("dbconf", "")
	//
	//for err != nil {
	//	logr.Error("数据库连接已经断开，重新连接", err)
	//	instance, err = opendb(dbconf)
	//	time.Sleep(time.Duration(5) * time.Second)
	//}
	//fmt.Printf("instance")
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
