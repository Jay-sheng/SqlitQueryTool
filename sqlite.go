package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var tablename = make([]string, 1)
var seltabname = make([]string, 1)
var textname = make([]string, 1)

func Dateinitialize(dbname string) {
	var name string
	//打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println("open2 db fail =", err)
		panic(err)
	}
	//查询表格名称
	rows, err1 := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err1 != nil {
		fmt.Println("query db fail =", err1)
		panic(err1)
	}
	i := 1
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("rows.scan fail =", err)
			panic(err)
		}
		tablename = append(tablename, name)
		seltabname = append(seltabname, "SELECT * FROM "+name)
		textname = append(textname, "./log/"+name)
		i++
	}
	rows.Close()
}
func InputData(dbname string, number int) {
	var lever string
	var updatestamp time.Time
	var upload int
	var body_data *string
	var body_text string
	//打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println("open1 dn fail =", err)
		panic(err)
	}
	var yue, ri int
	var t1, t2 time.Time
	fmt.Println("请输入要导出数据的日期（例如:4 25（表示：2018年4月25日） ）:")
	fmt.Scan(&yue, &ri)
	//根据输入的时间计算时间戳
	t1 = time.Date(2018, time.Month(yue), ri, 0, 0, 0, 0, t1.Location())
	t2 = time.Date(2018, time.Month(yue), ri, 24, 0, 0, 0, t2.Location())
	fmt.Println("数据查询写入中，请耐心等待...")
	//查询数据库中的数据
	rows1, err := db.Query(seltabname[number])
	if err != nil {
		fmt.Println("query db1 fail =", err)
		panic(err)
	}
	//根据要打印的表格新建相应的文件
	f, err := os.Create(textname[number] + ".txt")
	if err != nil {
		fmt.Println("os .create err =", err)
		return
	}
	defer f.Close()
	i := '1'
	//将读取到的数据写入文件
	for rows1.Next() {
		err = rows1.Scan(&lever, &updatestamp, &upload, &body_data, &body_text)
		if err != nil {
			fmt.Println("rows.sann1 fail =", err)
			panic(err)
		}
		if t1.Unix() < updatestamp.Unix() && t2.Unix() > updatestamp.Unix() {
			//获取文件大小
			FileSize, _ := f.Seek(0, os.SEEK_END)
			if FileSize/1048576 < 300 {
				//按照指定的格式写文件
				w := bufio.NewWriter(f)
				w.WriteString(lever)
				w.WriteString("\t")
				w.WriteString(updatestamp.Format("2006-01-02 15:04:05.000"))
				w.WriteString("\t")
				w.WriteString(body_text)
				w.WriteString("\n")
				w.Flush()
			} else {
				//当文件大小超过100M，重新生成新文件
				f.Close()
				f, err = os.Create(textname[number] + string(i) + ".txt")
				if err != nil {
					fmt.Println("os .create err =", err)
					return
				}
				i++
				w := bufio.NewWriter(f)
				w.WriteString(lever)
				w.WriteString("\t")
				w.WriteString(updatestamp.Format("2006-01-02 15:04:05.000"))
				w.WriteString("\t")
				w.WriteString(body_text)
				w.WriteString("\n")
				w.Flush()
			}
		}
	}

	rows1.Close()
}
func GetDbname(dirPth string) (files []string, dirs []string, err error) {
	//获取指定目录下的所有文件和目录
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Println("ReadDir fail =", err)
		return
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { //目录，递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetDbname(dirPth + PthSep + fi.Name())

		} else {
			//过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".db")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files, dirs, nil
}
func main() {
	files, dirs, err := GetDbname("./log")
	if err != nil {
		fmt.Println("GetDbname err =", err)
		return
	}
	for _, table := range dirs {
		tmp, _, _ := GetDbname(table)
		for _, tmp1 := range tmp {
			files = append(files, tmp1)
		}
	}
	i := 1
	for _, table1 := range files {
		fmt.Printf("数据库文件(%d):%s\n", i, table1)
		i++
	}
	var d int
	fmt.Println("请输入要打开的数据库文件序号:")
	fmt.Scan(&d)
	dbname := files[d-1]
	Dateinitialize(dbname)
	for {
		var a int
		fmt.Println("请输入表格序号:")
		for i := 1; i < len(tablename); i++ {
			fmt.Println(i, tablename[i])
		}
		fmt.Println("结束程序请输入：0")
		fmt.Scan(&a)
		if a == 0 {
			break
		} else {

			InputData(dbname, a)

		}
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
