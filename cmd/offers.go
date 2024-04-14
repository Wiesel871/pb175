package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	comp "wiesel/pb175/components"
	_ "wiesel/pb175/database"
)


func (st *GlobalState) AddOffer(w http.ResponseWriter, r *http.Request) {
}

func (st *GlobalState) GetUserOffers(w http.ResponseWriter, r *http.Request) {
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

    offers, err := st.DBH.GetOffersByOwner(id)
    id_viewer, own := LogedInOwned(r, user)
    if id_viewer > -1 {
        http.SetCookie(w, NewSession(id_viewer))
    }
    comp.Offers(offers, own, id_viewer).Render(r.Context(), w)
}

func (st *GlobalState) GetOffer(w http.ResponseWriter, r *http.Request) {
    id_owner, err := strconv.Atoi(r.PathValue("id_owner"))
    if err != nil {
        w.WriteHeader(404)
        println("owner err")
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        w.WriteHeader(404)
        println("id err")
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    user, err := st.DBH.GetUserById(id_owner)
    if err != nil {
        w.WriteHeader(404)
        fmt.Printf("user err: %v\n", err)
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }


    offer, err := st.DBH.GetOfferById(id)
    if err != nil || offer.ID_owner != id_owner {
        w.WriteHeader(404)
        println("incorrect combo", id_owner, err.Error())
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    id_viewer, own := LogedInOwned(r, user)
    if id_viewer > -1 {
        http.SetCookie(w, NewSession(id_viewer))
    }
    comp.Offer(offer, own, id_viewer).Render(r.Context(), w)
}


func (st *GlobalState) GetOffers(w http.ResponseWriter, r *http.Request) {
    offers, err := st.DBH.GetOffers()
    if err != nil {
        w.WriteHeader(404)
        //comp.PageNotFound().Render(r.Context(), w)
        return
    }
    id, _ := LogedInOwned(r, nil)
    if id > -1 {
        http.SetCookie(w, NewSession(id))
    }
    comp.OffersPage(offers, false, id).Render(r.Context(), w)
}
