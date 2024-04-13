package database

import (
	"database/sql"
	"fmt"

    "golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID int
    Name string
    Email string
    Password []byte
    hasPFP bool
}

type Users = []User

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func NewUser(id int, name string, email string, password string) (*User, error) {
    hashedPassword, error := HashPassword(password)
    if error != nil {
        return nil, error
    }
    return &User{
        ID: id,
        Name: name,
        Email: email,
        Password: hashedPassword,
    }, nil
}

func (dbh *dbHandler) InsertUser(con *User) error {
    _, err := dbh.db.Exec(`
    INSERT INTO Users (ID, Name, Email, Password) 
    VALUES (?, ?, ?, ?, ?)`, 
    con.ID, con.Name, con.Email, con.Password, false)
    fmt.Printf("con.id: %v\n", con.ID)
    return err
}

func (dbh *dbHandler) GetUsers() Users {
    rows, err := dbh.db.Query("SELECT Name, Email FROM" + dbh.Users)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
    }
	defer rows.Close()

	var contacts Users
	for rows.Next() {
		var contact User
		rows.Scan(&contact.Name, &contact.Email)
		contacts = append(contacts, contact)
	}
	return contacts
}

func initUsers(db *sql.DB) (string, error) {
    name := "Users"
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + name +` (
			ID INTEGER PRIMARY KEY,
			Name TEXT NOT NULL,
            Email TEXT NOT NULL UNIQUE,
            Password BLOB NOT NULL,
            hasPFP BOOL,
            isAdmin BOOL,
		)`); 
    return name, err
}
