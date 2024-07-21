package handlers

import (
	"fmt"
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
	LogoutEndpoint					= "/logout"
	LoginEndpoint					= "/login"
	HomeEndpoint					= "/home"
	FinanceEndpoint					= "/finance"
	FinanceTrackEndpoint			= "/finance/track"
	FinanceTrackConfirmEndpoint		= "/finance/track/confirm"
	FinanceFeedEndpoint				= "/finance/feed"
	FinanceFeedConfirmEndpoint		= "/finance/feed/confirm"
	FinanceFeedEditEndpoint			= "/finance/feed/edit"
	SessionName						= "session-name"
	Username 						= "useranme"
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
	InvalidInput 					= "Bad Username and/or Passoword"
	InvalidExpenseInput 			= "Name, Amount and Category cannot be empty!"
	FailedToAddExpense 				= "Sorry! Failed to add expense"
	FailedToUpdateExpense 			= "Sorry! Failed to update expense"
	SuccAddExpense 					= "Expense successfully added!"
	SuccUpdateExpense				= "Expense successfully updated!"
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
		http.Redirect(w, r, HomeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
		util.PrintMessage(GET, LoginEndpoint)

        response, err := templates.GetLogin()
        if err != nil {
			util.PrintError(err)
            http.Error(w, InternalServerError, http.StatusInternalServerError)
            return
        }
        w.Write([]byte(response))
		util.PrintSuccess(GET, LoginEndpoint)
        return
    }

	util.PrintMessage(POST, LoginEndpoint)

	username := r.FormValue(Username)
	password := r.FormValue(Password)

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

	http.Redirect(w, r, HomeEndpoint, http.StatusSeeOther)
	util.PrintSuccess(POST, LoginEndpoint)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, LogoutEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = Expire
	session.Save(r, w)

	http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
	util.PrintSuccess(POST, LogoutEndpoint)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(GET, HomeEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	response, err := templates.GetHome()
	
	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	util.PrintSuccess(GET, HomeEndpoint)
}

func FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(GET, FinanceEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[UserId].(string)
	response, err := GetFinance(userId)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))

	util.PrintSuccess(GET, FinanceTrackEndpoint)
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
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[UserId].(string)

	if r.Method == http.MethodGet {
		util.PrintMessage(GET, FinanceTrackEndpoint)

		response, err := templates.GetFinanceTrack(userId)
		if err != nil {
			util.PrintError(err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		util.PrintSuccess(GET, FinanceTrackEndpoint)
		return
	}

	util.PrintMessage(POST, FinanceTrackEndpoint)

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
	util.PrintSuccess(POST, FinanceTrackEndpoint)
}

func FinanceTrackConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, FinanceTrackConfirmEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
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
	util.PrintSuccess(POST, FinanceTrackConfirmEndpoint)
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
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId, ok:= session.Values[UserId].(string)
	if !ok {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		util.PrintMessage(GET, FinanceFeedEndpoint)

		response, err := templates.GetFinanceFeed(userId)
		if !err {
			util.PrintError(err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		util.PrintSuccess(GET, FinanceFeedEndpoint)
		w.Write([]byte(response))
		return
	}

	util.PrintMessage(POST, FinanceFeedEndpoint)

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
	util.PrintSuccess(POST, FinanceFeedEndpoint)
}

func FinanceFeedEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, FinanceFeedEditEndpoint)

	session, err := store.Get(r, SessionName)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[Username] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(amount), 64)
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
	util.PrintSuccess(POST, FinanceFeedEditEndpoint)
}

func FinanceFeedConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	util.PrintMessage(POST, FinanceFeedConfirmEndpoint)

	session, _ := store.Get(r, SessionName)
	if session.Values[Username] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(amount), 64)
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
	util.PrintSuccess(POST, FinanceFeedConfirmEndpoint)
}