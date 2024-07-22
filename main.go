package main

import (
    "net/http"
    "home-app/app/dao"
    "home-app/app/handlers"
    "home-app/app/util"
)

func main() {
    // Initialize the database
    err := dao.InitDB()
    if err != nil {
        util.PrintError(err)
    }
    defer dao.CloseDB()

    handlers.StoreInit()

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

    // Serve static files from the web/static directory
    staticFileDirectory := http.Dir(util.WebStaticDir)
    staticFileHandler := http.StripPrefix(util.StaticDir, http.FileServer(staticFileDirectory))
    http.Handle(util.StaticDir, staticFileHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
    })

    util.PrintSuccess("Listening on :8080")
    http.ListenAndServe(":8080", nil)
}
