package main

import (
	"database/sql"
	"fmt"
)

type MysqlUtil struct {
	db *sql.DB
}

func checkErr(err error) {
	if err != nil {
		err.Error()
	}
}

func NewMysqlUtil() *MysqlUtil {
	fmt.Println("open the database, redpacket")
	db, err := sql.Open("mysql", "test:12345678@/magicmaze")
	checkErr(err)
	return &MysqlUtil{
		db: db,
	}
}

func (util *MysqlUtil) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return util.db.Query(sql, args...)
}

func (util *MysqlUtil) Insert(sql string, args ...interface{}) (int64, error) {
	stmt, err := util.db.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(args...)
	checkErr(err)
	return res.LastInsertId()
}

func (util *MysqlUtil) Update(sql string, args ...interface{}) {
	stmt, err := util.db.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(args...)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
}
