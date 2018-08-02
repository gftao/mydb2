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
	"github.com/jinzhu/copier"
	"runtime"
)

type EngineDb2 struct {
	db *sql.DB
}

type EngStmt struct {
	sync.Locker
	Tx        *sql.Tx
	query     interface{}
	Args      []interface{}
	tableName string
}

type levelmap map[int]map[string]map[string]string

func opendb(file string) (*EngineDb2, error) {
	talkChan := loadConfig(file)
	dbdata := make(map[string]string)
 	for t := range talkChan {
		for k := range talkChan[t] {
			switch k {
			case "db2":
				dbdata = talkChan[t][k]
			}
		}
	}
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

func (tx *EngStmt) Exec(query string, args ...interface{}) error {
	_, err := tx.Tx.Exec(query, args ...)
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
	results, err := e.db.Query(query.(string), args...)
	return results, err
}

func (e *EngineDb2) QueryRow(query interface{}, args ...interface{}) (*sql.Row) {
	results := e.db.QueryRow(query.(string), args...)
	return results
}
func (tx *EngStmt) QueryRow(query interface{}, args ...interface{}) (*sql.Row) {
	results := tx.Tx.QueryRow(query.(string), args...)
	return results
}

func (tx *EngStmt) Query(query interface{}, args ...interface{}) (*sql.Rows, error) {
	results, err := tx.Tx.Query(query.(string), args...)
	return results, err
}

func (tx *EngStmt) CheckExit(query interface{}, args ...interface{}) (bool, error) {
	rows, err := tx.Tx.Query(query.(string), args...)
	if rows.Next() {
		return true, err
	}
	return false, err
}

func (tx *EngStmt) Where(query interface{}, args ...interface{}) *EngStmt {
	tx.query = query
	tx.Args = args
	return tx //&EngStmt{Query: query, Args: args}
}

func (tx *EngStmt) Get(bean interface{}) (b bool, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
			b = false
		}
	}()

	v := reflect.ValueOf(bean).Elem()
	t := reflect.TypeOf(bean).Elem()
	if !v.CanAddr() {
		return false, errors.New("copy to value is unaddressable")
	}
	if tb, ok := v.Interface().(TableName); ok {
		tx.tableName = tb.TableName()
	}
	qs := "select * from " + tx.tableName + " where " + tx.query.(string)
	//fmt.Println(qs, tx.Args)
	rows, err := tx.Tx.Query(qs, tx.Args...)
	if err != nil {
		return false, err
	}
	count := t.NumField()
	//fmt.Println(count)
	ci := 0
	for rows.Next() {
		ci ++
		id := make([]interface{}, count)
		it := make([]interface{}, count)
		for j := 0; j < count; j++ {
			k := v.Field(j).Type()
			id[j] = reflect.New(k)
			it[j] = &id[j]
		}
		err = rows.Scan(it...)
		if err != nil {
			return false, err
		}
		dec := mahonia.NewDecoder("gbk")
		for j := 0; j < count; j++ {
			//k := reflect.ValueOf(id[j]).Kind()
			k := v.Field(j).Kind()
			//fmt.Printf("%v->%s\n", k, t.Field(j).Name)
			switch k {
			case reflect.Int, reflect.Int32, reflect.Int64:
				//str := strconv.FormatInt(reflect.ValueOf(id[j]).Int(), 10)
				//i, _ := strconv.Atoi(str)
				//v.Field(j).Set(reflect.ValueOf(i))
				v.Field(j).Set(reflect.ValueOf(id[j]))
			case reflect.String:
				data := reflect.ValueOf(id[j]).Interface().([]byte)
				str := string(data)
				if runtime.GOOS == "windows" {
					r := dec.ConvertString(str)
					v.Field(j).Set(reflect.ValueOf(r))
				} else {
					v.Field(j).Set(reflect.ValueOf(str))
				}
				//fmt.Println(r)

			case reflect.Slice:
				data := reflect.ValueOf(id[j]).Interface().([]byte)
				str := string(data)
				r := dec.ConvertString(str)
				//fmt.Println(r)
				v.Field(j).Set(reflect.ValueOf(r))
			case reflect.Struct:
				v.Field(j).Set(reflect.ValueOf(id[j]))
			case reflect.Float64:
				v.Field(j).Set(reflect.ValueOf(id[j]))
			default:
				v.Field(j).Set(reflect.ValueOf(id[j]))
				//reflect.Copy(v.Field(j),reflect.ValueOf(id[j]))
				fmt.Println("----default----", reflect.TypeOf(id[j]).Kind(), id[j])
			}
			//v.Field(j).Set(reflect.ValueOf(r))
		}
	}
	//fmt.Println(ci)
	if ci == 0 {
		return true, ErrRecordNotFound
	}
	return true, nil
}

func (tx *EngStmt) GetAll(rowsSlicePtr interface{}) (b bool, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
			b = false
		}
	}()
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	//fmt.Printf("%+v\n", sliceValue)
	if sliceValue.Kind() != reflect.Slice {
		return false, errors.New("needs a pointer to a slice")
	}
	sliceElementType := sliceValue.Type().Elem()
	beans := reflect.MakeSlice(reflect.SliceOf(sliceElementType), 0, 4)
	bean := reflect.New(sliceElementType)
	v := reflect.Indirect(bean)
	if tb, ok := bean.Interface().(TableName); ok {
		tx.tableName = tb.TableName()
	}
	qs := "select * from " + tx.tableName + " where " + tx.query.(string)
	//fmt.Println(qs, tx.Args)
	rows, err := tx.Tx.Query(qs, tx.Args...)
	if err != nil {
		return false, err
	}
	count := sliceElementType.NumField()

	id := make([]interface{}, count)
	it := make([]interface{}, count)
	for j := 0; j < count; j++ {
		k := v.Field(j).Type()
		id[j] = reflect.New(k)
		it[j] = &id[j]
	}
	ci := 0
	for rows.Next() {
		ci ++
		if ci == 1000 {
			break
		}
		err = rows.Scan(it...)
		if err != nil {
			return false, err
		}
		dec := mahonia.NewDecoder("gbk")
		for j := 0; j < count; j++ {
			k := v.Field(j).Kind()
			switch k {
			case reflect.Int, reflect.Int32:
				str := strconv.FormatInt(reflect.ValueOf(id[j]).Int(), 10)
				i, _ := strconv.Atoi(str)
				v.Field(j).Set(reflect.ValueOf(i))
			case reflect.String:
				data := reflect.ValueOf(id[j]).Interface().([]byte)
				str := string(data)
				if runtime.GOOS == "windows" {
					r := dec.ConvertString(str)
					v.Field(j).Set(reflect.ValueOf(r))
				} else {
					v.Field(j).Set(reflect.ValueOf(str))
				}
			case reflect.Slice:
				data := reflect.ValueOf(id[j]).Interface().([]byte)
				str := string(data)
				r := dec.ConvertString(str)
				v.Field(j).Set(reflect.ValueOf(r))
			case reflect.Struct:
				v.Field(j).Set(reflect.ValueOf(id[j]))
			default:
				v.Field(j).Set(reflect.ValueOf(id[j]))
				fmt.Println("----default----", reflect.TypeOf(id[j]).Kind(), id[j])
			}

		}
		beans = reflect.Append(beans, v)
	}
	if ci == 0 {
		return true, ErrRecordNotFound
	}
	copier.Copy(rowsSlicePtr, beans.Interface())
	return true, nil
}

func (s *EngStmt) FindOne(links interface{}, querystring string, args ...interface{}) (bool, error) {
	has, err := s.Where(querystring, args...).Get(links)
	return has, err
}

func (tx *EngStmt) Uptade(bean interface{}) (bool, error) {

	v := reflect.ValueOf(bean).Elem()
	t := reflect.TypeOf(bean).Elem()

	tblName := ""
	if tb, ok := v.Interface().(TableName); ok {
		tblName = tb.TableName()
	}
	//fmt.Println(tblName)
	qs := strings.Split(tx.query.(string), "?")
	if (len(qs) - 1) != len(tx.Args) {
		return false, errors.New("参数不匹配")
	}
	for i, s := range tx.Args {
		qs[i] += "'" + reflect.ValueOf(s).String() + "'"
	}
	//fmt.Println(qs)
	q := strings.Join(qs, "")
	///for update
	if true {
		SqlUp := "select * from " + tblName + " where " + q + " for update"
		//fmt.Println(SqlUp)
		r, err := tx.Tx.Exec(SqlUp)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		i, err := r.RowsAffected()
		if i == 0 {
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
		case reflect.Int, reflect.Int64, reflect.Float64:
			vv := fmt.Sprintf("%s = %d", t.Field(i).Name, v.Field(i).Interface())
			ns = append(ns, vv)
		default:
			if v.Field(i).String() != "" {
				vs := ""
				if strings.Contains(v.Field(i).String(), "'") {
					s := strings.Replace(v.Field(i).String(), "'", `''`, -1)
					vs = fmt.Sprintf("%s = '%s'", t.Field(i).Name, s)

				} else {
					vs = fmt.Sprintf("%s = '%s'", t.Field(i).Name, v.Field(i).String())
				}
				ns = append(ns, vs)
			}
		}
	}

	Sql := "update " + tblName + " set " + strings.Join(ns, ",") + " where " + q
	//fmt.Println(Sql)
	r, err := tx.Tx.Exec(Sql)
	if err != nil {
		//fmt.Println(err)
		return false, err
	}
	i, err := r.RowsAffected()
	if i == 0 {
		return false, err
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
		case reflect.Int, reflect.Int64, reflect.Float64:
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
			if t.Field(i).Name == "REC_UPD_TS" || t.Field(i).Name == "REC_CRT_TS" {
				vvv := "(VALUES TIMESTAMP(CURRENT TIMESTAMP))"
				nv = append(nv, vvv)
			} else {
				vvv := "(select current date from sysibm.sysdummy1)"
				nv = append(nv, vvv)
			}

		default:
			if strings.Contains(v.Field(i).String(), "'") {
				vs := strings.Replace(v.Field(i).String(), "'", `''`, -1)
				//fmt.Println(vs)
				nv = append(nv, "'"+vs+" '")

			} else {
				nv = append(nv, "'"+v.Field(i).String()+"'")
			}
		}

	}
	su := tblName + " (" + strings.Join(ns, ", ") + ")"
	sv := " VALUES (" + strings.Join(nv, ", ") + ")"
	Sql := `INSERT INTO ` + su + sv //+" VALUES (" + ")"
	//fmt.Println(Sql)
	_, err := s.Tx.Exec(Sql)
	if err != nil {
		fmt.Println(err)
		return false, err
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
