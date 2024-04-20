package utility

import "net/http"
import "time"
import "strconv"
import "wiesel/pb175/state"
import db "wiesel/pb175/database"

type Response = func(w http.ResponseWriter, r *http.Request)

type GSP = *state.GlobalState


func GetClientID(r *http.Request) int {
    cookie, err := r.Cookie(Session)
    if err != nil {
        return -1
    }
    id, err := strconv.Atoi(cookie.Value)
    if err != nil {
        return -1
    }
    return id
}

const (
    Session     = "NeUplneBazosSKSessionIdentifierCookie"
    SessionTimeout = 24 * time.Hour 
    CookieMaxAge   = int(SessionTimeout / time.Second)
)

func NewSession(id int) *http.Cookie {
    expiration := time.Now().Add(SessionTimeout)
    return &http.Cookie{
        Name:    Session,
        Value:   strconv.Itoa(id),
        Expires: expiration,
        MaxAge:  CookieMaxAge,
        Path:    "/",
        HttpOnly: true,
    }
}

func GetUser(st GSP, id int) *db.User {
    user, err := st.DBH.GetUserById(id)
    if err != nil {
        user = st.Anonym
    }
    return user
}

