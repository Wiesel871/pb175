package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	comp "wiesel/pb175/components"
	data "wiesel/pb175/database"

	_ "github.com/a-h/templ"
)



func (st *GlobalState) AddUser(w http.ResponseWriter, r *http.Request) {
    println("signup")
    ID, err := data.SmallestMissingID(st.DBH.DB, st.DBH.Users)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("err: %v\n", err)
        comp.SignUpForm(err.Error()).Render(r.Context(), w)
        return
    }
    name := r.FormValue("name")
    email := r.FormValue("email")
    password := r.FormValue("psw")
    println(strconv.Itoa(ID) + name + email + password)

    user, err := data.NewUser(ID, name, email, password)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("err: %v\n", err)
        comp.SignUpForm(err.Error()).Render(r.Context(), w)
        return
    }
    err = st.DBH.InsertUser(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("err: %v\n", err)
        comp.SignUpForm(err.Error()).Render(r.Context(), w)
        return
    }

    fmt.Printf("login user: %v\n", user)
    fmt.Printf("This resource has been moved permanently to /profile/" + strconv.Itoa(user.ID))
    http.Redirect(w, r, "/profile/" + strconv.Itoa(user.ID), http.StatusMovedPermanently)
}

func (st *GlobalState) Login(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    //password := r.FormValue("psw")
    user, err := st.DBH.GetUserByEmail(email)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Printf("err: %v\n", err)
        comp.SignUpForm(err.Error()).Render(r.Context(), w)
        return
    }
    /*
    pswHash, err := data.HashPassword(password)
    if pswHash != user.Password {
        w.WriteHeader(400)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }
    */

    fmt.Printf("login user: %v\n", user)
    fmt.Printf("This resource has been moved permanently to /profile/" + strconv.Itoa(user.ID))
    http.Redirect(w, r, "/profile/" + strconv.Itoa(user.ID), http.StatusMovedPermanently)
}


func (st *GlobalState) Profile(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        w.WriteHeader(404)
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    user, err := st.DBH.GetUserById(id)
    if err != nil {
        w.WriteHeader(404)
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    comp.UserPage(user).Render(r.Context(), w)
}

func (st *GlobalState) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := st.DBH.GetUsers()
    if err != nil {

    }
    comp.Users(users).Render(r.Context(), w)
}
