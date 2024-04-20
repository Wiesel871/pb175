package client

import (
    "net/http"

    comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"

    _"github.com/a-h/templ"
)

func SignUpForm(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        if id > -1 {
            http.SetCookie(w, ut.NewSession(id))
            http.Redirect(w, r, "/home", http.StatusMovedPermanently)
            println("already signed up")
            return
        }
        comp.Page(comp.SignUpForm(""), st.Anonym, comp.All).Render(r.Context(), w)
    }
}

func LogInForm(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        if id > -1 {
            http.SetCookie(w, ut.NewSession(id))
            http.Redirect(w, r, "/home", http.StatusMovedPermanently)
            println("already signed up")
            return
        }
        comp.Page(comp.LogInForm(""), st.Anonym, comp.All).Render(r.Context(), w)
    }
}
