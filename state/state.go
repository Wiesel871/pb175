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
