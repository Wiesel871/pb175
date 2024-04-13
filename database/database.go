package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)




func InitDB() (*sql.DB, error) {
	var err error
    db, err := sql.Open("sqlite3", "bazos.db")
	if err != nil {
        fmt.Printf("open err.Error(): %v\n", err.Error())
		return nil, err
	}
    if err = initUsers(db); err != nil {
        fmt.Printf("init Users err.Error(): %v\n", err.Error())
		return nil, err
    }

    if err = initOffers(db); err != nil {
        fmt.Printf("init Offers err.Error(): %v\n", err.Error())
		return nil, err

    }
	return db, nil
}
