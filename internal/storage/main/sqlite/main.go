package sqlite

import (
	"database/sql"

	"github.com/Izumra/SKUD_OKEI/lib/config"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewConnetion(cfg *config.Config) *Storage {
	db, err := sql.Open(cfg.Db.DriverName, cfg.Db.SourcePath)
	if err != nil {
		panic(err)
	}

	return &Storage{
		db,
	}
}
