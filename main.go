package main

import (
    "log"
    "fmt"
    "net/http"
    "home-app/app/dao"
    "home-app/app/handlers"
)

func main() {
    // Initialize the database
    err := dao.InitDB()
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer dao.CloseDB()

    // Route definitions
    http.HandleFunc(handlers.LoginEndpoint, handlers.LoginHandler)
    http.HandleFunc(handlers.HomeEndpoint, handlers.HomeHandler)
    http.HandleFunc(handlers.FinanceEndpoint, handlers.FinanceHandler)
    http.HandleFunc(handlers.FinanceTrackEndpoint, handlers.FinanceTrackHandler)
    http.HandleFunc(handlers.FinanceTrackConfirmEndpoint, handlers.FinanceTrackConfirmHandler)
    http.HandleFunc(handlers.FinanceFeedEndpoint, handlers.FinanceFeedHandler)
    http.HandleFunc(handlers.FinanceFeedConfirmPoint, handlers.FinanceFeedConfirmHandler)
    http.HandleFunc(handlers.FinanceFeedEditEndpoint, handlers.FinanceFeedEditHandler)
    http.HandleFunc(handlers.LogoutEndpoint, handlers.LogoutHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, handlers.LoginEndpoint, http.StatusSeeOther)
    })

    fmt.Println("Server listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
