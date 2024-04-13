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

func SmallestMissingID(db *sql.DB, table string) (int, error) {
    var res int
    err := db.QueryRow(`
    SELECT MIN(ID + 1) AS smallest_missing_id
    FROM your_table t1
    WHERE NOT EXISTS (
        SELECT 1
        FROM ` + table + ` t2
        WHERE t2.ID = t1.ID + 1)
    `).Scan(&res)
    if err != nil {
        return 0, err
    }
    return res, nil
}

func InitDB() (*DBHandler, error) {
    res := new(DBHandler)
	var err error
    db, err := sql.Open("sqlite3", "bazos.db")
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
