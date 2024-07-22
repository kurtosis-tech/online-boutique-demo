package cartstore

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	db *gorm.DB
}

func NewDb(
	dbConnInfo *connectionInfo,
) (*Db, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbConnInfo.host, dbConnInfo.username, dbConnInfo.password, dbConnInfo.databaseName, dbConnInfo.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("An error occurred opening the connection to the database with dsn %s", dsn))
	}

	return &Db{
		db: db,
	}, nil
}

func (db *Db) Close() error {
	sqlDb, err := db.db.DB()
	if err != nil {
		return errors.Wrap(err, "An error occurred closing the database connection")
	}

	if err = sqlDb.Close(); err != nil {
		return errors.Wrap(err, "An error occurred closing the database connection")
	}

	return nil
}
