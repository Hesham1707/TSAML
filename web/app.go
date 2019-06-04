package web

import (
	"fmt"
	"github.com/hesham/tsaml/web/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	http.HandleFunc("/request.html", app.RequestHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/request.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:4000/) ...")
	http.ListenAndServe(":4000", nil)
}



