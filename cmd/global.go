package cmd

import (
	"fmt"
	_ "fmt"
	"net/http"
	comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"
    st "wiesel/pb175/state"
    ct "wiesel/pb175/cmd/client"
    us "wiesel/pb175/cmd/user"
    of "wiesel/pb175/cmd/offers"

	_"github.com/a-h/templ"
)



func Home(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        user := ut.GetUser(st, id)
        comp.Page(comp.IndexBody(), user, comp.HomeN).Render(r.Context(), w)
    }
}


func SetupUserHandler(mux *http.ServeMux, st *st.GlobalState) {
    fmt.Printf("st.Anonym: %v\n", st.Anonym)
    mux.HandleFunc("/", Home(st))
    mux.HandleFunc("GET /home", Home(st))

    fs := http.FileServer(http.Dir("images"))
    mux.Handle("GET /images/", http.StripPrefix("/images/", fs))

    mux.HandleFunc("POST /logout", ct.LogOut)

    mux.HandleFunc("/signup", ct.SignUpForm(st))
    mux.HandleFunc("POST /signup", ct.SignUp(st))

    mux.HandleFunc("/login", ct.LogInForm(st))
    mux.HandleFunc("POST /login", ct.LogIn(st))

    mux.HandleFunc("/profile", ct.RedirectToUser(st))

    mux.HandleFunc("GET /profile/{id}", us.Profile(st))
    mux.HandleFunc("GET /users", us.GetAllUsers(st))

    mux.HandleFunc("POST /change_details" , us.ChangeDetails(st))

    mux.HandleFunc("GET /profile/{id}/offers", of.GetUserOffers(st))
    mux.HandleFunc("GET /profile/{id_owner}/offers/{id}", of.GetOffer(st))

    mux.HandleFunc("GET /offers", of.GetOffers(st))

    mux.HandleFunc("GET /add_offer", of.AddOffer(st))
    mux.HandleFunc("POST /add_offer", of.UploadOffer(st))

}
