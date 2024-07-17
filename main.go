package main

import (
    "log"
    "fmt"
    "net/http"
	"home-app/app/dao"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
	
	err := dao.InitDB() 
    
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer dao.CloseDB()

    http.HandleFunc("/", handler)
    fmt.Println("Server listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}
