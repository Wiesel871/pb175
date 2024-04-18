package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID int
    Name string
    Email string
    Password []byte
    Details string
    HasPFP bool
    IsAdmin bool
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

func CheckPasswordHash(password string, hash []byte) error {
    return bcrypt.CompareHashAndPassword(hash, []byte(password))
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
    INSERT INTO Users (ID, Name, Email, Password, Details, HasPFP, IsAdmin) 
    VALUES (?, ?, ?, ?, ?, ?, ?)`, 
    con.ID, con.Name, con.Email, con.Password, "<empty>", false, false)
    return err
}

func (dbh *DBHandler) GetUserById(id int) (*User, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE ID = ?", id)
    var u = new(User)
    if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Details, &u.HasPFP, &u.IsAdmin); err != nil {
        return nil, err
    }
    return u, nil
}

func (dbh *DBHandler) GetUserByEmail(email string) (*User, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE Email = ?", email)
    var u = new(User)
    if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Details, &u.HasPFP, &u.IsAdmin); err != nil {
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
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Details, &user.HasPFP, &user.IsAdmin); err != nil {
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
            Details TEXT NOT NULL,
            HasPFP BOOL NOT NULL,
            IsAdmin BOOL NOT NULL
		)`); 
    pass, _ := HashPassword("Foch258147")
    _, _ = db.Exec(`
    INSERT OR REPLACE INTO Users (ID, Name, Email, Password, Details, HasPFP, IsAdmin) 
    VALUES (0, "Filip", "filpavlovic@gmail.com", ?, "Him", false, true)`,
    pass)

    return name, err
}
