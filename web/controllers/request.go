package controllers

import (
	"net/http"
)

func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		ClientName := r.FormValue("ClientName")
		ClientID := r.FormValue("ClientID")
		ClientBalance := r.FormValue("ClientBalance")

		args := []string{ClientID,ClientName,ClientBalance}

		txid, err := app.Fabric.InvokeHello("createClient",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "request.html", data)
}
