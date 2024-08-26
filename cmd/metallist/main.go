package main

import (
	"log"
	"fmt"
	"net/http"
	// "metallist/internal/services/anilist"
	"metallist/internal/services"
	// "metallist/internal/config"
)

func main() {
	// Server: services auth
	mux := http.NewServeMux()
	auth_services := services.AuthServices()
	for _, service := range auth_services {
	    mux.HandleFunc(service.LoginPathURL(), services.LoginHandler(service))
	    mux.HandleFunc(service.CallbackPathURL(), services.CallbackHandler(service))
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":1212"),
		Handler: mux,
	}

	// err := services.TestSaveTokens()
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// err := services.TestLoadTokens()
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	// Initialize cache database
	db, err := services.OpenBadger("/config/db/cache.db")
	if err != nil {
		fmt.Println("Error opening badger database: ", err)
		// Handle error
	}
	defer db.Close()

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
	    log.Printf("%v", err)
	} else {
	    log.Println("Server closed!")
	}
}
