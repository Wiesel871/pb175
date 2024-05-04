package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DBHandler struct {
    DB *sql.DB
    Users string
    Offers string
}

func InitDB(dbn string) (*DBHandler, error) {
    res := new(DBHandler)
	var err error
    db, err := sql.Open("sqlite3", dbn)
	if err != nil {
        fmt.Printf("open err.Error(): %v\n", err.Error())
		return nil, err
	}
    res.DB = db
    var name string
    if name, err = initUsers(db); err != nil {
        fmt.Printf("init Users err.Error(): %v\n", err.Error())
		return nil, err
    }
    res.Users = name

    if name, err = initOffers(db); err != nil {
        fmt.Printf("init Offers err.Error(): %v\n", err.Error())
		return nil, err

    }
    res.Offers = name
	return res, nil
}
