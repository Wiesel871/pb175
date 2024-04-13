package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Offer struct {
    ID int
    ID_owner int
    Name string
    Description string
}


type Offers = []Offer



func initOffers(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS Offers (
            ID INTEGER PRIMARY KEY,
            ID_owner INTEGER,
            Name TEXT NOT NULL,
            Description TEXT,
            FOREIGN KEY (ID_owner) REFERENCES Users(ID)
        )`) 
    return err
}
