package client

import (
    "net/http"
    "strconv"
    "time"
    "os"
    "fmt"

    comp "wiesel/pb175/components"
    db "wiesel/pb175/database"
    ut "wiesel/pb175/cmd/utility"

    _ "github.com/a-h/templ"
)


func SignUp(st ut.GSP) ut.Response { 
    return func(w http.ResponseWriter, r *http.Request) {
        ID := time.Now().Unix()
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("psw")

        user, err := db.NewUser(ID, name, email, password)
        if err != nil {
            w.WriteHeader(500)
            return
        }

        err = os.Mkdir("images/" + strconv.FormatInt(ID, 10), 0755) 
        if err != nil {
            fmt.Println(err) 
            return
        }
        err = st.DBH.InsertUser(user)
        if err != nil {
            w.WriteHeader(422)
            comp.SignUpForm("email or name in use").Render(r.Context(), w)
            return
        }

        

        http.SetCookie(w, ut.NewSession(user.ID))
        http.Redirect(w, r, "/profile/" + strconv.FormatInt(user.ID, 10), http.StatusFound)
    }
}

func LogOut(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now().Add(-time.Hour)
    cookie := http.Cookie{
        Name:    ut.Session,
        Value:   "-1",
        Expires: expiration,
        MaxAge:  -1,
        Path:    "/",
    }
    http.SetCookie(w, &cookie)

    http.Redirect(w, r, "/", http.StatusFound)
}


func LogIn(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        email := r.FormValue("email")
        password := r.FormValue("psw")
        user, err := st.DBH.GetUserByEmail(email)
        if err != nil {
            w.WriteHeader(422)
            comp.LogInForm("Email not found").Render(r.Context(), w)
            return
        }
        if err = db.CheckPasswordHash(password, user.Password); err != nil {
            w.WriteHeader(422)
            comp.LogInForm("Incorrect password").Render(r.Context(), w)
            return
        }

        http.SetCookie(w, ut.NewSession(user.ID))
        http.Redirect(w, r, "/profile/" + strconv.FormatInt(user.ID, 10), http.StatusFound)
    }
}

func RedirectToUser(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        if id == -1 {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/profile/" + strconv.FormatInt(id, 10), http.StatusFound)
    }
}
