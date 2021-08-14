package main

import (
	"encoding/json"
	"net/http"
)

type Api struct{}

func (a *Api) episodes(w http.ResponseWriter, r *http.Request) {
	db := new(DB).getConnection()

	episodes := []Episode{}
	db.Find(&episodes)

	json.NewEncoder(w).Encode(episodes)
}
