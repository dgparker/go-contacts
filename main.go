package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgparker/go-contacts/handler"
	mw "github.com/dgparker/go-contacts/middleware"
	"github.com/dgparker/go-contacts/services"
	"github.com/julienschmidt/httprouter"
)

var dbURI string
var dbname string
var dbcoll string
var port string

func main() {
	fmt.Println(dbURI)
	fmt.Println(dbname)
	fmt.Println(dbcoll)
	fmt.Println(port)

	contactService := services.NewMongoService(dbURI, dbname, dbcoll)

	contactHandler := handler.NewContactHandler(contactService)

	r := httprouter.New()
	r.GET("/entry", mw.Add(contactHandler.HandleGetEntries, mw.SetHeaders))
	r.GET("/entry/:id", mw.Add(contactHandler.HandleGetEntryByID, mw.SetHeaders))
	r.POST("/entry", mw.Add(contactHandler.HandlePostEntry, mw.SetHeaders))
	r.PUT("/entry", mw.Add(contactHandler.HandlePutEntry, mw.SetHeaders))
	r.DELETE("/entry/:id", mw.Add(contactHandler.HandleDeleteEntryByID, mw.SetHeaders))
	r.GET("/csv/entry", mw.Add(contactHandler.HandleGetCSV, mw.SetHeaders))
	r.POST("/csv/entry", mw.Add(contactHandler.HandlePostCSV, mw.SetHeaders))

	fmt.Printf("Server listening on port: %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
