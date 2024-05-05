package database

import (
	"database/sql"
    "fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Offer struct {
    ID          int64
    OwnerID    int64
    Name        string
    Description string
}


type Offers = []Offer

func NewOffer(id, id_owner int64, name, desc string) (*Offer) {
    return &Offer{
        ID: id,
        OwnerID: id_owner,
        Name: name,
        Description: desc,
    }
}

func (dbh *DBHandler) InsertOffer(of *Offer) error {
    _, err := dbh.DB.Exec(`
    INSERT INTO Offers (ID, ID_owner, Name, Description) 
    VALUES (?, ?, ?, ?)`, 
    of.ID, of.OwnerID, of.Name, of.Description)
    return err
}

func (dbh *DBHandler) GetOffers(by, sc, fil string) (Offers, error) {
    rows, err := dbh.DB.Query("SELECT * FROM " + dbh.Offers + " WHERE Name LIKE '%" + fil + "%' COLLATE NOCASE ORDER BY " + by + " " + sc)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
        return nil, err
    }
	defer rows.Close()

	var offers Offers
	for rows.Next() {
		var offer Offer
        _ = rows.Scan(&offer.ID, &offer.OwnerID, &offer.Name, &offer.Description)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (dbh *DBHandler) GetOffersByOwner(id int64) (Offers, error) {
    rows, err := dbh.DB.Query("SELECT * FROM " + dbh.Offers + " WHERE ID_owner = ?", id)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
        return nil, err
    }
	defer rows.Close()

	var offers Offers
	for rows.Next() {
		var offer Offer
		rows.Scan(&offer.ID, &offer.OwnerID, &offer.Name, &offer.Description)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (dbh *DBHandler) GetOfferById(id int64) (*Offer, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Offers + " WHERE ID = ?", id)
    var offer = new(Offer)
    if err := row.Scan(&offer.ID, &offer.OwnerID, &offer.Name, &offer.Description); err != nil {
        return nil, err
    }
    return offer, nil
}

func (dbh *DBHandler) AdjustOffer(o *Offer, name, desc string) error {
    _, err := dbh.DB.Exec(`
    UPDATE ` + dbh.Offers + ` 
    SET Name = ?, Description = ?
    WHERE ID = ?;`, 
    name, desc, o.ID)
    return err
}

func (dbh *DBHandler) DeleteOffer(id int64) error {
    _, err := dbh.DB.Exec(`
    DELETE FROM ` + dbh.Offers + `
    WHERE ID = ?
    `, id)
    return err
}


func initOffers(db *sql.DB) (string, error) {
    name := "Offers"
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + name + ` (
            ID INTEGER PRIMARY KEY,
            ID_owner INTEGER,
            Name TEXT NOT NULL,
            Description TEXT,
            FOREIGN KEY (ID_owner) REFERENCES Users(ID)
        )`) 
    _, _ = db.Exec(`
    INSERT INTO ` + name + ` (ID, ID_owner, Name, Description) 
    VALUES (0, 0, "test1", "idk")`)
    
    _, _ = db.Exec(`
    INSERT INTO ` + name + ` (ID, ID_owner, Name, Description) 
    VALUES (1, 0, "test2", "idk")`)

    return name, err
}
