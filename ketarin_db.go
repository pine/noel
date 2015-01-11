package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type KetarinDb struct {
    db *sql.DB
}

func NewKetarinDb(path string) *KetarinDb {
    kdb := KetarinDb{}
    
    if db, err := sql.Open("sqlite3", path); err != nil {
        return nil
    } else {
        kdb.db = db
        return &kdb
    }
}

func (kdb *KetarinDb) SetSetting(name string, value string) error {
    tx, err := kdb.db.Begin()
	if err != nil {
		return err
	}
    
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO settings VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
    
    if _, err = stmt.Exec(name, value); err != nil {
        return err
    }
    
	return tx.Commit()
}

func (kdb *KetarinDb) Close() error {
    return kdb.db.Close()
}