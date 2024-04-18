package handlers

import (
	_"fmt"
	"net/http"
	"strconv"

	comp "wiesel/pb175/components"
	_ "wiesel/pb175/database"
)


func (st *GlobalState) AddOffer(w http.ResponseWriter, r *http.Request) {

}

func (st *GlobalState) GetUserOffers(w http.ResponseWriter, r *http.Request) {
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

    offers, err := st.DBH.GetOffersByOwner(id)
    comp.Page(comp.Offers(offers, owner, client), client, comp.All)
}

func (st *GlobalState) GetOffer(w http.ResponseWriter, r *http.Request) {
    user := GetClientID(r)
    if user > -1 {
        http.SetCookie(w, NewSession(user))
    }

    client, err := st.DBH.GetUserById(user)
    if err != nil {
        client = st.Anonym
    }

    id_owner, err := strconv.Atoi(r.PathValue("id_owner"))
    if err != nil {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }

    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }

    owner, err := st.DBH.GetUserById(id_owner)
    if err != nil {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }


    offer, err := st.DBH.GetOfferById(id)
    if err != nil || offer.ID_owner != id_owner {
        w.WriteHeader(404)
        comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        return
    }
    comp.Page(comp.Offer(offer, owner, client), client, comp.All).Render(r.Context(), w)
}


func (st *GlobalState) GetOffers(w http.ResponseWriter, r *http.Request) {
    offers, _ := st.DBH.GetOffers()
    id := GetClientID(r)
    if id > -1 {
        http.SetCookie(w, NewSession(id))
    }
    client, err := st.DBH.GetUserById(id)
    if err != nil {
        client = st.Anonym
    }

    comp.Page(comp.Offers(offers, nil, client), client, comp.OffersN).Render(r.Context(), w)
}
