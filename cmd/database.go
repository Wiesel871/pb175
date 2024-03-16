package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)


type Contact struct {
    id uint
    Name string
    Email string
}

type Contacts = []Contact

func newContact(id uint, name string, email string) (*Contact) {
    return &Contact{
        id: id,
        Name: name,
        Email: email,
    }
}

func insertContact(db *sql.DB, con *Contact) error {
    _, err := db.Exec(fmt.Sprintf("INSERT OR REPLACE INTO Contacts (id, Name, Email) VALUES (%d, \"%s\", \"%s\")", con.id, con.Name, con.Email))
    return err
}

func (st *State) getContacts() Contacts {
    rows, err := st.DB.Query("SELECT Name, Email FROM Contacts")
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
    }
	defer rows.Close()

	var contacts Contacts
	for rows.Next() {
		var contact Contact
		rows.Scan(&contact.Name, &contact.Email)
		contacts = append(contacts, contact)
	}
	return contacts
}



func initDB() (*sql.DB, error) {
	var err error
    db, err := sql.Open("sqlite3", "contacts.db")
	if err != nil {
        fmt.Printf("open err.Error(): %v\n", err.Error())
		return nil, err
	}
	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Contacts (
			id UNSIGNED INTEGER PRIMARY KEY,
			Name TEXT,
            Email TEXT NOT NULL UNIQUE
		)`); 
    err != nil {
        fmt.Printf("init err.Error(): %v\n", err.Error())
		return nil, err
	}
	return db, nil
}
