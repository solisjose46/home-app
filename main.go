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
    http.HandleFunc(handlers.loginEndpoint, handlers.LoginHandler)
    http.HandleFunc(handlers.homeEndpoint, handlers.HomeHandler)
    http.HandleFunc(handlers.financeEndpoint, handlers.FinanceHandler)
    http.HandleFunc(handlers.financeTrackEndpoint, handlers.FinanceTrackHandler)
    http.HandleFunc(handlers.financeTrackConfirmEndpoint, handlers.FinanceTrackConfirmHandler)
    http.HandleFunc(handlers.financeFeedEndpoint, handlers.FinanceFeedHandler)
    http.HandleFunc(handlers.financeFeedConfirmPoint, handlers.FinanceFeedConfirmHandler)
    http.HandleFunc(handlers.financeFeedEditEndpoint, handlers.FinanceFeedEditHandler)
    http.HandleFunc("/logout", handlers.LogoutHandler)

    // Default route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, handlers.loginEndpoint, http.StatusSeeOther)
    })

    fmt.Println("Server listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
