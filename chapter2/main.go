package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type TableRecord struct {
	id        int
	modify_at int
}

var DB *sql.DB
var EmptyRecord *TableRecord = new(TableRecord)
var ErrorNotFound = errors.New("record not exist")

func main() {
	err := InitDatabase()
	if err != nil {
		fmt.Printf("Connect Fail: %+v \n", err)
		os.Exit(1)
	}

	record, err := GetRow(100)
	if err != nil {
		fmt.Printf("Query Error: %+v \n", err)
		os.Exit(1)
	}

	fmt.Println(*record)

}

//数据库连接
func InitDatabase() error {
	DB, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test_go_client?parseTime=true")
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		return errors.Wrapf(err, "failed to connect database")
	}
	return nil
}

//查询
func GetRow(id int) (*TableRecord, error) {
	record := new(TableRecord)
	rows, err := DB.Query("select id,modify_at from test_select where  id=?", 1)

	//sql错误
	if err != nil {
		return nil, errors.Wrapf(err, "sql error")
	}

	//关闭rows
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&record.id, &record.modify_at)
	}

	//空记录
	if *record == *EmptyRecord {
		return nil, errors.Wrap(ErrorNotFound, "empty data")
	}

	return record, rows.Err()
}
