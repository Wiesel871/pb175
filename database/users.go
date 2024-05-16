package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID          int64
    Name        string
    Email       string
    Password    []byte
    Details     string
    HasPFP      bool
    IsAdmin     bool
}

type Users = []User


func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password string, hash []byte) error {
    return bcrypt.CompareHashAndPassword(hash, []byte(password))
}


func NewUser(id int64, name string, email string, password string) (*User, error) {
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
    INSERT INTO ` + dbh.Users + ` 
    (ID, Name, Email, Password, Details, HasPFP, IsAdmin) 
    VALUES (?, ?, ?, ?, ?, ?, ?)`, 
    con.ID, con.Name, con.Email, con.Password, "", false, false)
    return err
}

func (dbh *DBHandler) GetUserById(id int64) (*User, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE ID = ?", id)
    var u = new(User)
    if err := row.Scan(
        &u.ID, &u.Name, &u.Email, &u.Password, 
        &u.Details, &u.HasPFP, &u.IsAdmin); 
    err != nil {
        return nil, err
    }
    return u, nil
}

func (dbh *DBHandler) GetUserByEmail(email string) (*User, error) {
    row := dbh.DB.QueryRow("SELECT * FROM " + dbh.Users + " WHERE Email = ?", email)
    var u = new(User)
    if err := row.Scan(
        &u.ID, &u.Name, &u.Email, &u.Password, 
        &u.Details, &u.HasPFP, &u.IsAdmin); 
    err != nil {
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
        if err := rows.Scan(
            &user.ID, &user.Name, &user.Email, 
            &user.Password, &user.Details, &user.HasPFP, &user.IsAdmin);
        err != nil {
            return nil, err
        }
		users = append(users, user)
	}
	return users, nil
}

func (dbh *DBHandler) AdjustUser(
    u *User, name, email, 
    details string, hasPFP bool) error {
    _, err := dbh.DB.Exec(`
    UPDATE ` + dbh.Users + ` 
    SET Name = ?, Email = ?, Password = ?, Details = ?, HasPFP = ?
    WHERE ID = ?`, 
    name, u.Email, u.Password, details, hasPFP, u.ID)
    return err
}

func (dbh *DBHandler) Promote(id int64) error {
    _, err := dbh.DB.Exec(`
    UPDATE ` + dbh.Users + ` 
    SET IsAdmin = True
    WHERE ID = ?`, 
    id)
    return err
}

func (dbh *DBHandler) Demote(id int64) error {
    _, err := dbh.DB.Exec(`
    UPDATE ` + dbh.Users + ` 
    SET IsAdmin = False
    WHERE ID = ?`, 
    id)
    return err
}

func (dbh *DBHandler) DeleteUser(id int64) error {
    tx, err := dbh.DB.Begin()
    if err != nil {
        return err
    }
    offers, err := dbh.GetOffersByOwner(id)
    if err != nil {
        tx.Rollback()
        return err
    }
    for _, offer := range offers {
        _, err = tx.Exec(`
        DELETE FROM ` + dbh.Offers + `
        WHERE ID = ?
        `, offer.ID)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    _, err = tx.Exec(`
    DELETE FROM ` + dbh.Users + `
    WHERE ID = ?
    `, id)
    if err == nil {
        tx.Commit()
    }
    return err
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
    return name, err
}
