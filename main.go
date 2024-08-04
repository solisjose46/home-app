package main

import (
    "net/http"
    "home-app/app/dao"
    "home-app/app/handle"
    "home-app/app/util"
    "github.com/solisjose46/pretty-print/debug"
    "github.com/gorilla/sessions"
)

func main() {
    debug.PrintInfo(main, "Starting server")
    // Initialize the database
    dao := &dao.Dao{}
    err := dao.InitDB()
    if err != nil {
        debug.PrintError(main, err)
        return
    }
    defer dao.CloseDB()

    cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))

	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production
	}

    handler := handle.NewHandler(dao, cookieStore)
    
    // Route definitions
    http.HandleFunc(util.LoginEndpoint, handler.LoginHandler)
    http.HandleFunc(util.HomeEndpoint, handler.HomeHandler)
    http.HandleFunc(util.FinanceEndpoint, handler.FinanceHandler)
    http.HandleFunc(util.FinanceTrackEndpoint, handler.FinanceTrackHandler)
    http.HandleFunc(util.FinanceTrackConfirmEndpoint, handler.FinanceTrackConfirmHandler)
    http.HandleFunc(util.FinanceFeedEndpoint, handler.FinanceFeedHandler)
    http.HandleFunc(util.FinanceFeedConfirmEndpoint, handler.FinanceFeedConfirmHandler)
    http.HandleFunc(util.FinanceFeedEditEndpoint, handler.FinanceFeedEditHandler)
    http.HandleFunc(util.LogoutEndpoint, handler.LogoutHandler)

    // Serve static files from the web/static directory
    staticFileDirectory := http.Dir(util.WebStaticDir)
    staticFileHandler := http.StripPrefix(util.StaticDir, http.FileServer(staticFileDirectory))
    http.Handle(util.StaticDir, staticFileHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
    })

    debug.PrintSucc(main, "Listening on :8080")
    http.ListenAndServe(":8080", nil)
}
