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

	if r.FormValue("submittedgettransactions") == "true" {
		ClientID := r.FormValue("ClientID")
		txid, err := app.Fabric.GetLogs(ClientID)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	if r.FormValue("submittedcheckanamoly") == "true" {
		ClientID := r.FormValue("ClientID")
		txid, err := app.Fabric.CheckAnomaly(ClientID)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	if r.FormValue("submittedlistusers") == "true" {

		txid, err := app.Fabric.GetClientNames()
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
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
	
	 if r.FormValue("submittedgetbalance") == "true" {
		ClientID := r.FormValue("ClientID")

		args := []string{ClientID}

		txid, err := app.Fabric.InvokeHello("getBalance",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 if r.FormValue("submittedwithdraw") == "true" {
		ClientID := r.FormValue("ClientID")
		amount := r.FormValue("Amount")

		args := []string{ClientID,amount}

		txid, err := app.Fabric.InvokeHello("withdraw",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	if r.FormValue("submittedtransfer") == "true" {
		SClientID := r.FormValue("SClientID")
		TClientID := r.FormValue("TClientID")
		amount := r.FormValue("Amount")

		args := []string{SClientID,TClientID,amount}

		txid, err := app.Fabric.InvokeHello("transfer",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 if r.FormValue("submitteddeposit") == "true" {
		ClientID := r.FormValue("ClientID")
		amount := r.FormValue("Amount")

		args := []string{ClientID,amount}

		txid, err := app.Fabric.InvokeHello("deposit",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 
	 if r.FormValue("submittedaddstock") == "true" {
		StockName := r.FormValue("StockName")
		SharesNumber := r.FormValue("SharesNumber")
		OwnerID := r.FormValue("OwnerID")
		Price := r.FormValue("Price")

		args := []string{StockName,SharesNumber,OwnerID,Price}

		txid, err := app.Fabric.InvokeHello("addStocks",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 if r.FormValue("submittedremovestock") == "true" {
		StockName := r.FormValue("StockName")
		OwnerID := r.FormValue("OwnerID")

		args := []string{StockName,OwnerID,}

		txid, err := app.Fabric.InvokeHello("removeStocks",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 if r.FormValue("submittedbuystock") == "true" {
		StockName := r.FormValue("StockName")
		SharesNumber := r.FormValue("SharesNumber")
		OwnerID := r.FormValue("OwnerID")
		Price := r.FormValue("Price")

		args := []string{OwnerID,StockName,Price,SharesNumber}

		txid, err := app.Fabric.InvokeHello("buyStocks",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	 if r.FormValue("submittedsellstock") == "true" {
		StockName := r.FormValue("StockName")
		SharesNumber := r.FormValue("SharesNumber")
		OwnerID := r.FormValue("OwnerID")
		Price := r.FormValue("Price")

		args := []string{StockName,OwnerID,Price,SharesNumber}

		txid, err := app.Fabric.InvokeHello("sellStocks",args)
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}

	renderTemplate(w, r, "request.html", data)
}
