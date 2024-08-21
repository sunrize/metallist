package main

import (
	"log"
	"fmt"
	"net/http"
	// "metallist/internal/services/anilist"
	"metallist/internal/services/auth"
	// "metallist/internal/config"
)

func main() {
	// Server: services auth
	mux := http.NewServeMux()
	auth_services := auth.AuthServices()
	for _, service := range auth_services {
	    mux.HandleFunc(service.LoginPathURL(), auth.LoginHandler(service))
	    mux.HandleFunc(service.CallbackPathURL(), auth.CallbackHandler(service))
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":1212"),
		Handler: mux,
	}

	// err := auth.TestSaveTokens()
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// err := auth.TestLoadTokens()
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
	    log.Printf("%v", err)
	} else {
	    log.Println("Server closed!")
	}
}
