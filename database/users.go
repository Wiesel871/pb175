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
    isAdmin bool
}

type Users = []User



/* Hashes given password */
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}


/* Creates new user struct from parameters */
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

func (dbh *DBHandler) InsertUser(con *User) error {
    _, err := dbh.DB.Exec(`
    INSERT INTO Users (ID, Name, Email, Password, hasPFP, isAdmin) 
    VALUES (?, ?, ?, ?, ?, ?)`, 
    con.ID, con.Name, con.Email, con.Password, false, false)
    fmt.Printf("con.id: %v\n", con.ID)
    return err
}

func (dbh *DBHandler) GetUserById(id int) (*User, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE ID = ?", id)
    var u = new(User)
    if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.hasPFP, &u.isAdmin); err != nil {
        return nil, err
    }
    return u, nil
}

func (dbh *DBHandler) GetUsers() (Users, error) {
    rows, err := dbh.DB.Query("SELECT * FROM " + dbh.Users)
    if err != nil {
        fmt.Printf("get err.Error(): %v\n", err.Error())
        return nil, err
    }
	defer rows.Close()

	var users Users
	for rows.Next() {
		var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, user.Password, &user.hasPFP, &user.isAdmin); err != nil {
            return nil, err
        }
		users = append(users, user)
	}
	return users, nil
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
