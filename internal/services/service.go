package services

import (
	"context"
	"fmt"
	// "log"
	"net/http"
	// "os"
    // "time"
	// "errors"
	// "encoding/json"
    // "io"

	// "metallist/internal/config"
	// "golang.org/x/oauth2"
	"metallist/internal/services"
)

// type Service struct {
//     ID   string `json:"id"`
//     Name string `json:"name"`
// }

// func devCache() {
//     // Save some data
//     err = cache.SaveJSON(db, "my_unique_key", myData)
//     if err != nil {
//         // Handle error
//     }

//     // Load some data
//     var loadedData interface{}
//     err = cache.LoadJSON(db, "my_unique_key", &loadedData)
//     if err != nil {
//         // Handle error
//     }

//     // Do something with loadedData
//     fmt.Println("Loaded data: ", loadedData)
// }