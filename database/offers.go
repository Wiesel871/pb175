package database

import (
	"database/sql"
    "fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Offer struct {
    ID int
    ID_owner int
    Name string
    Description string
    HasPhoto bool
}


type Offers = []Offer

func NewOffer(id, id_owner int, name, desc string) (*Offer) {
    return &Offer{
        ID: id,
        ID_owner: id_owner,
        Name: name,
        Description: desc,
        HasPhoto: false,
    }
}

func (dbh *DBHandler) InsertOffer(of *Offer) error {
    _, err := dbh.DB.Exec(`
    INSERT INTO Offers (ID, ID_owner, Name, Description) 
    VALUES (?, ?, ?, ?, ?)`, 
    of.ID, of.ID_owner, of.Name, of.Description, of.HasPhoto)
    fmt.Printf("of.id: %v\n", of.ID)
    return err
}

func (dbh *DBHandler) GetOffers() (Offers, error) {
    rows, err := dbh.DB.Query("SELECT (*) FROM" + dbh.Offers)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
        return nil, err
    }
	defer rows.Close()

	var offers Offers
	for rows.Next() {
		var offer Offer
		rows.Scan(&offer.ID, &offer.ID_owner, &offer.Name, &offer.Description, &offer.HasPhoto)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (dbh *DBHandler) GetOffersByOwner(id int) (Offers, error) {
    rows, err := dbh.DB.Query("SELECT * FROM " + dbh.Offers + " WHERE ID_owner = ?", id)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
        return nil, err
    }
	defer rows.Close()

	var offers Offers
	for rows.Next() {
		var offer Offer
		rows.Scan(&offer.ID, &offer.ID_owner, &offer.Name, &offer.Description)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (dbh *DBHandler) GetOfferById(id int) (*Offer, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE ID = ?", id)
    var offer = new(Offer)
    if err := row.Scan(&offer.ID, &offer.ID_owner, &offer.Name, &offer.Description, &offer.HasPhoto); err != nil {
        return nil, err
    }
    return offer, nil
}


func initOffers(db *sql.DB) (string, error) {
    name := "Offers"
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + name + ` (
            ID INTEGER PRIMARY KEY,
            ID_owner INTEGER,
            Name TEXT NOT NULL,
            Description TEXT,
            HasPhoto BOOL,
            FOREIGN KEY (ID_owner) REFERENCES Users(ID)
        )`) 
    return name, err
}