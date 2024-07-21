package main

import (
    "log"
    "fmt"
    "net/http"
    "home-app/app/dao"
    "home-app/app/handlers"
    "home-app/app/util"
)

func main() {
    // Initialize the database
    err := dao.InitDB()
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer dao.CloseDB()

    // Route definitions
    http.HandleFunc(util.LoginEndpoint, handlers.LoginHandler)
    http.HandleFunc(util.HomeEndpoint, handlers.HomeHandler)
    http.HandleFunc(util.FinanceEndpoint, handlers.FinanceHandler)
    http.HandleFunc(util.FinanceTrackEndpoint, handlers.FinanceTrackHandler)
    http.HandleFunc(util.FinanceTrackConfirmEndpoint, handlers.FinanceTrackConfirmHandler)
    http.HandleFunc(util.FinanceFeedEndpoint, handlers.FinanceFeedHandler)
    http.HandleFunc(util.FinanceFeedConfirmEndpoint, handlers.FinanceFeedConfirmHandler)
    http.HandleFunc(util.FinanceFeedEditEndpoint, handlers.FinanceFeedEditHandler)
    http.HandleFunc(util.LogoutEndpoint, handlers.LogoutHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
    })

    fmt.Println("Server listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
