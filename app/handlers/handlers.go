package handlers

import (
	"github.com/gorilla/sessions"
	"home-app/app/dao"
	"home-app/app/models"
	"home-app/app/templates"
	"home-app/app/util"
	"net/http"
	"strconv"
)

// TODO: import from ignored file
var store = sessions.NewCookieStore([]byte("something-very-secret"))

const (
	SessionName						= "session-name"
	Username 						= "username"
	Password 						= "password"
	UserId 							= "user-id"
	Amount 							= "amount"
	ExpenseId 						= "expense-id"
	Name 							= "name"
	Category 						= "category"
	InternalServerError 			= "Internal Server Error"
	Redirect 						= "Redirect 303"
	MethodNotAllowed 				= "Method Not Allowed"
    GET 							= "GET "
    POST 							= "POST "
	Expire	 						= -1
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] != nil {
		http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
		util.PrintMessage(GET, util.LoginEndpoint)

        response, err := templates.GetLogin()
        if err != nil {
			util.PrintError(err)
            http.Error(w, InternalServerError, http.StatusInternalServerError)
            return
        }
        w.Write([]byte(response))
		util.PrintSuccess(GET, util.LoginEndpoint)
        return
    }

	util.PrintMessage(POST, util.LoginEndpoint)

	username := r.FormValue(Username)
	password := r.FormValue(Password)

	util.PrintMessage("username\n", username)
	util.PrintMessage("password\n", password)


	response, err := templates.PostLogin(username, password)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if response != "" {
		w.Write([]byte(response))
		return
	}

	userId, err := dao.GetUserId(username)
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Values[Username] = username
	session.Values[UserId] = userId
	session.Save(r, w)

	http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
	util.PrintSuccess(POST, util.LoginEndpoint)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, util.LogoutEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = Expire
	session.Save(r, w)

	http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
	util.PrintSuccess(POST, util.LogoutEndpoint)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(GET, util.HomeEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	response, err := templates.GetHome()
	
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	util.PrintSuccess(GET, util.HomeEndpoint)
}

func FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(GET, util.FinanceEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[UserId].(string)
	response, err := templates.GetFinance(userId)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))

	util.PrintSuccess(GET, util.FinanceTrackEndpoint)
}

func FinanceTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[UserId].(string)

	if r.Method == http.MethodGet {
		util.PrintMessage(GET, util.FinanceTrackEndpoint)

		response, err := templates.GetFinanceTrack(userId)
		if err != nil {
			util.PrintError(err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		util.PrintSuccess(GET, util.FinanceTrackEndpoint)
		return
	}

	util.PrintMessage(POST, util.FinanceTrackEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(Amount), 64)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		UserId:   userId,
		Name:     r.FormValue(Name),
		Amount:   amount,
		Category: r.FormValue(Category),
	}

	response, err := templates.PostFinanceTrack(expense)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte(response))
	util.PrintSuccess(POST, util.FinanceTrackEndpoint)
}

func FinanceTrackConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, util.FinanceTrackConfirmEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(Amount), 64)
	
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		UserId:   session.Values[UserId].(string),
		Name:     r.FormValue(Name),
		Amount:   amount,
		Category: r.FormValue(Category),
	}

	response, err := templates.PostFinanceTrackConfirm(expense)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	util.PrintSuccess(POST, util.FinanceTrackConfirmEndpoint)
}

func FinanceFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId, ok:= session.Values[UserId].(string)
	if !ok {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		util.PrintMessage(GET, util.FinanceFeedEndpoint)

		response, err := templates.GetFinanceFeed(userId)
		if err != nil {
			util.PrintError(err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		util.PrintSuccess(GET, util.FinanceFeedEndpoint)
		w.Write([]byte(response))
		return
	}

	util.PrintMessage(POST, util.FinanceFeedEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(Amount), 64)
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(ExpenseId),
		UserId:    session.Values[UserId].(string),
		Name:      r.FormValue(Name),
		Amount:    amount,
		Category:  r.FormValue(Category),
	}

	response, err := templates.PostFinanceFeed(expense)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	util.PrintSuccess(POST, util.FinanceFeedEndpoint)
}

func FinanceFeedEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, util.FinanceFeedEditEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(Amount), 64)
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(ExpenseId),
		UserId:    session.Values[UserId].(string),
		Name:      r.FormValue(Name),
		Amount:    amount,
		Category:  r.FormValue(Category),
	}

	response, err := templates.PostFinanceFeedEdit(expense)
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
	util.PrintSuccess(POST, util.FinanceFeedEditEndpoint)
}

func FinanceFeedConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, util.FinanceFeedConfirmEndpoint)

	session, _ := store.Get(r, SessionName)
	if session.Values[Username] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(Amount), 64)
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue(ExpenseId),
		UserId:    session.Values[UserId].(string),
		Name:      r.FormValue(Name),
		Amount:    amount,
		Category:  r.FormValue(Category),
	}

	response, err := templates.PostFinanceFeedConfirm(expense)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	util.PrintSuccess(POST, util.FinanceFeedConfirmEndpoint)
}