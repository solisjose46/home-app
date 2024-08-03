package handlers

import (
	"github.com/solisjose46/pretty-print/debug"
	"github.com/gorilla/sessions"
	"home-app/app/dao"
	"home-app/app/models"
	"home-app/app/templates"
	"home-app/app/util"
	"net/http"
	"strconv"
	"fmt"
)

// TODO: import from ignored file
var store = sessions.NewCookieStore([]byte("something-very-secret"))

const (
	SessionName						= "session-name"
	SessionUserId					= "user-id"
	SessionUsername 				= "username"
	FormUsername 					= "username"
	FormPassword					= "password"
	FormUserId						= "user-id"
	FormAmount 						= "amount"
	FormExpenseId 					= "expense-id"
	FormName 						= "name"
	FormCategory 					= "category"
	InternalServerError 			= "Internal Server Error"
	Redirect 						= "Redirect 303"
	MethodNotAllowed 				= "Method Not Allowed"
    GET 							= "GET "
    POST 							= "POST "
	Expire	 						= -1
)

func StoreInit() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] != nil {
		http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
		debug.PrintInfo(LoginHandler, GET, util.LoginEndpoint)

        response, err := templates.GetLogin()
        if err != nil {
			debug.PrintError(LoginHandler, err)
            http.Error(w, InternalServerError, http.StatusInternalServerError)
            return
        }
        w.Write([]byte(response))
		debug.PrintSucc(LoginHandler, GET, util.LoginEndpoint)
        return
    }

	debug.PrintInfo(LoginHandler, POST, util.LoginEndpoint)

	username := r.FormValue(FormUsername)
	password := r.FormValue(FormPassword)

	response, err := templates.PostLogin(username, password)

	if err != nil {
		debug.PrintError(LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if response != "" {
		debug.PrintSucc(LoginHandler, POST, util.LoginEndpoint, "Server response")
		w.Write([]byte(response))
		return
	}

	userId, err := dao.GetUserId(username)
	if err != nil {
		debug.PrintError(LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Values[SessionUsername] = username
	session.Values[SessionUserId] = userId
	err = session.Save(r, w)

	if err != nil {
		debug.PrintError(LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	fmt.Println("Session Values:")
	for key, value := range session.Values {
		fmt.Printf("%s: %v\n", key, value)
	}

	http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
	debug.PrintSucc(LoginHandler, POST, util.LoginEndpoint)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(LogoutHandler, POST, util.LogoutEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(LogoutHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = Expire
	session.Save(r, w)

	http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
	debug.PrintSucc(LogoutHandler, POST, util.LogoutEndpoint)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(HomeHandler, GET, util.HomeEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(HomeHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	for key, value := range session.Values {
		fmt.Printf("%s: %v\n", key, value)
	}

	if session.Values[SessionUsername] == nil {
		debug.PrintInfo(HomeHandler, "redirecting to login")
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	response, err := templates.GetHome()
	
	if err != nil {
		debug.PrintError(HomeHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	debug.PrintSucc(HomeHandler, GET, util.HomeEndpoint)
}

func FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(FinanceHandler, GET, util.FinanceEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		debug.PrintInfo(FinanceHandler, "redirecting to login")
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[SessionUserId].(string)
	response, err := templates.GetFinance(userId)

	if err != nil {
		debug.PrintError(FinanceHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))

	debug.PrintSucc(FinanceHandler, GET, util.FinanceTrackEndpoint)
}

func FinanceTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[SessionUserId].(string)

	if r.Method == http.MethodGet {
		debug.PrintInfo(FinanceTrackHandler, GET, util.FinanceTrackEndpoint)

		response, err := templates.GetFinanceTrack(userId)
		if err != nil {
			debug.PrintError(FinanceTrackHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		debug.PrintSucc(FinanceTrackHandler, GET, util.FinanceTrackEndpoint)
		return
	}

	debug.PrintInfo(FinanceTrackHandler, POST, util.FinanceTrackEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)

	if err != nil {
		debug.PrintError(FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		UserId:   userId,
		Name:     r.FormValue(FormName),
		Amount:   amount,
		Category: r.FormValue(FormCategory),
	}

	response, err := templates.PostFinanceTrack(expense)

	if err != nil {
		debug.PrintError(FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte(response))
	debug.PrintSucc(FinanceTrackHandler, POST, util.FinanceTrackEndpoint)
}

func FinanceTrackConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(FinanceTrackConfirmHandler, POST, util.FinanceTrackConfirmEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	
	if err != nil {
		debug.PrintError(FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		UserId:   session.Values[SessionUserId].(string),
		Name:     r.FormValue(FormName),
		Amount:   amount,
		Category: r.FormValue(FormCategory),
	}

	response, err := templates.PostFinanceTrackConfirm(expense)

	if err != nil {
		debug.PrintError(FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	debug.PrintSucc(FinanceTrackConfirmHandler, POST, util.FinanceTrackConfirmEndpoint)
}

func FinanceFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId, ok:= session.Values[SessionUserId].(string)
	if !ok {
		debug.PrintError(FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		debug.PrintInfo(FinanceFeedHandler, GET, util.FinanceFeedEndpoint)

		response, err := templates.GetFinanceFeed(userId)
		if err != nil {
			debug.PrintError(FinanceFeedHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		debug.PrintSucc(FinanceFeedHandler, GET, util.FinanceFeedEndpoint)
		w.Write([]byte(response))
		return
	}

	debug.PrintSucc(FinanceFeedHandler, POST, util.FinanceFeedEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUsername].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := templates.PostFinanceFeed(expense)

	if err != nil {
		debug.PrintError(FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	debug.PrintSucc(FinanceFeedHandler, POST, util.FinanceFeedEndpoint)
}

func FinanceFeedEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(FinanceFeedEditHandler, POST, util.FinanceFeedEditEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUserId].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := templates.PostFinanceFeedEdit(expense)
	if err != nil {
		debug.PrintError(FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
	debug.PrintSucc(FinanceFeedEditHandler, POST, util.FinanceFeedEditEndpoint)
}

func FinanceFeedConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(FinanceFeedConfirmHandler, POST, util.FinanceFeedConfirmEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		debug.PrintError(FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUserId].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := templates.PostFinanceFeedConfirm(expense)

	if err != nil {
		debug.PrintError(FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	debug.PrintSucc(FinanceFeedConfirmHandler, POST, util.FinanceFeedConfirmEndpoint)
}