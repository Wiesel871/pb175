package handlers

import (
	_"fmt"
	"net/http"
	"strconv"

	comp "wiesel/pb175/components"
	data "wiesel/pb175/database"

	_ "github.com/a-h/templ"
)



func (st *GlobalState) SignUp(w http.ResponseWriter, r *http.Request) {
    ID, err := data.SmallestMissingID(st.DBH.DB, st.DBH.Users)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        comp.SignUpForm(err.Error()).Render(r.Context(), w)
        return
    }
    name := r.FormValue("name")
    email := r.FormValue("email")
    password := r.FormValue("psw")

    user, err := data.NewUser(ID, name, email, password)
    if err != nil {
        w.WriteHeader(500)
        return
    }
    err = st.DBH.InsertUser(user)
    if err != nil {
        w.WriteHeader(422)
        comp.SignUpForm("email or name in use").Render(r.Context(), w)
        return
    }

    http.SetCookie(w, NewSession(user.ID))
    http.Redirect(w, r, "/profile/" + strconv.Itoa(user.ID), http.StatusFound)
}

func (st *GlobalState) Login(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("psw")
    user, err := st.DBH.GetUserByEmail(email)
    if err != nil {
        w.WriteHeader(422)
        comp.LogInForm("Email not found").Render(r.Context(), w)
        return
    }
    if err = data.CheckPasswordHash(password, user.Password); err != nil {
        w.WriteHeader(422)
        comp.LogInForm("Incorrect password").Render(r.Context(), w)
        return
    }

    http.SetCookie(w, NewSession(user.ID))
    http.Redirect(w, r, "/profile/" + strconv.Itoa(user.ID), http.StatusFound)
}

func (st *GlobalState) Profile(w http.ResponseWriter, r *http.Request) {
    user := GetClientID(r)
    if user > -1 {
        http.SetCookie(w, NewSession(user))
    }

    client, err := st.DBH.GetUserById(user)
    if err != nil {
        client = st.Anonym
    }

    id, err := strconv.Atoi(r.PathValue("id"))

    if err != nil {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }
    owner, err := st.DBH.GetUserById(id)
    if err != nil {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }
    curr_page := comp.All
    if client.ID == owner.ID {
        curr_page = comp.ProfileP
    }
    comp.Page(comp.Profile(owner, client), client, curr_page).Render(r.Context(), w)
}

func (st *GlobalState) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    user := GetClientID(r)
    if user > -1 {
        http.SetCookie(w, NewSession(user))
    }

    client, err := st.DBH.GetUserById(user)
    if err != nil {
        client = st.Anonym
    }

    if err != nil || !client.IsAdmin {
        w.WriteHeader(403)
        comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)

    }
    users, err := st.DBH.GetUsers()
    if err != nil {

    }
    comp.Page(comp.Users(users, client), client, comp.All).Render(r.Context(), w)
}
