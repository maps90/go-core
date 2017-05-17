package datasource

import (
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/maps90/go-core/log"
)

var (
	dbReadConn  *gorm.DB
	dbWriteConn *gorm.DB
)

type Mysql struct {
	logMode                    bool
	writerConn, readerConn     string
	maxOpenConns, maxIdleConns int
	err                        error
}

func NewMysql(write, read string) *Mysql {
	m := new(Mysql)
	m.writerConn = write
	m.readerConn = read

	return m
}
func (d *Mysql) Error() error {
	return d.err
}
func (d *Mysql) SetDebug(debug bool) {
	d.logMode = debug
}

func (d *Mysql) SetOpenConn(oc int) {
	d.maxOpenConns = oc
}

func (d *Mysql) SetIdleConn(ic int) {
	d.maxIdleConns = ic
}

func (d *Mysql) Write() *gorm.DB {
	if d.writerConn == "" {
		d.err = errors.New("missing writer configuration.")
	}
	if dbWriteConn == nil {
		dbWriteConn = d.createMysqlConn(d.writerConn)
	}

	return dbWriteConn
}

func (d *Mysql) Read() *gorm.DB {
	if d.readerConn == "" {
		d.err = errors.New("missing reader configuration.")
	}
	if dbReadConn == nil {
		dbReadConn = d.createMysqlConn(d.readerConn)
	}

	return dbReadConn
}

func (d *Mysql) createMysqlConn(descriptor string) *gorm.DB {
	db, err := gorm.Open("mysql", descriptor)
	if err != nil {
		log.New(log.ErrorLevelLog, "DB Connection Error: ", err.Error())
		os.Exit(1)
	}
	db.DB().SetMaxIdleConns(d.maxIdleConns)
	db.LogMode(d.logMode)
	db.DB().SetMaxOpenConns(d.maxOpenConns)
	return &db
}
