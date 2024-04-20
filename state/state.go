package state

import data "wiesel/pb175/database"

type GlobalState struct {
    DBH *data.DBHandler
    Anonym *data.User
}
