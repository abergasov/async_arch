package database

import (
	"fmt"
	"time"

	"async_arch/internal/config"
	"async_arch/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DBConnector interface {
	Client() *sqlx.DB
}

type DBConnect struct {
	db *sqlx.DB
}

func InitDBConnect(cnf *config.DBConf) *DBConnect {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cnf.Address, cnf.Port, cnf.User, cnf.Pass, cnf.DBName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error connect to db %s %s@%s:%s", cnf.DBName, cnf.User, cnf.Address, cnf.Port), err)
	}

	if cnf.MaxConnections == 0 {
		db.SetMaxOpenConns(10)
	} else {
		db.SetMaxOpenConns(cnf.MaxConnections) // максимальное число коннектов одновременных к бд
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(10) // число коннектов в пуле

	err = db.Ping()
	dbAddr := fmt.Sprintf("%s %s@%s:%s", cnf.DBName, cnf.User, cnf.Address, cnf.Port)
	if err != nil {
		logger.Fatal("error ping db", err, zap.String("db", dbAddr))
	}
	logger.Info("ping ok", zap.String("db", dbAddr))
	return &DBConnect{db}
}

func (d *DBConnect) Client() *sqlx.DB {
	return d.db
}
