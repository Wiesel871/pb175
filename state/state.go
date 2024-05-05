package state

import (
	"net/http"
	data "wiesel/pb175/database"
)

type GlobalState struct {
    DBH *data.DBHandler
    Anonym *data.User
    SRV *http.Server
}

func GetAnonym() *data.User {
    return &data.User{
        ID: -1,
        Name: "",
        Email: "",
        Password: []byte(""),
        IsAdmin: false,
        HasPFP: false,
    }
}
