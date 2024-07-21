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
	FinanceFeedConfirmPoint			= "/finance/feed/confirm"
	FinanceFeedEditEndpoint			= "/finance/feed/edit"
	SessionName						= "session-name"
	Username 						= "useranme"
	Password 						= "password"
	UserId 							= "user-id"
	Amount 							= "amount"
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

	loginAtttempt, err := templates.GetLoginAttempt(username, password)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if loginAtttempt != "" {
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

	htmlHome := util.GetFilePath(templates.TmplPath, templates.TmplHome, templates.HtmlExtension)

	http.ServeFile(w, r, htmlHome)
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

	http.Redirect(w, r, FinanceTrackEndpoint, http.StatusSeeOther)
	util.PrintMessage(Redirect, FinanceTrackEndpoint)
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

	response, err := templates.GetFinanceTrackValidateExpense(expense)

	if err != nil {
		util.PrintError(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if response != "" {
		w.Write([]byte(response))
		return
	}

	response, err := templates.GetFinanceTrackConfirm(expense)
	
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

	response, err := 

	if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
		response, err := templates.GetFinanceTrackServerResponse(
			models.ServerResponse{
				Message: InvalidExpenseInput,
				ReturnEndpoint: FinanceTrack,
			}
		)

		if err != nil {
			util.PrintError(err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return
	}

	// revisit this part
	// I dont dao ops here
	success, err := dao.AddExpense(expense)
	if err != nil || !success {
		response, _ := templates.GetFinanceTrackServerResponse()
		w.Write([]byte(response))
		return
	}

	response, _ := templates.GetFinanceTrackServerResponse(
		models.ServerResponse{
			Message: "Expense added successfully",
			ReturnEndpoint: FinanceTrackEndpoint,
	})
	w.Write([]byte(response))
}

func FinanceFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId, ok:= session.Values["userId"].(string)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := templates.GetFinanceFeed(userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
}

func FinanceFeedPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	expenseId := r.FormValue("expenseId")
	expense, err := dao.GetExpense(expenseId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := templates.GetFinanceFeedEdit(expense)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
}

func FinanceFeedEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// fix user id validation here
	expense := models.Expense{
		ExpenseId: r.FormValue("expenseId"),
		UserId:    session.Values["userId"].(string),
		Name:      r.FormValue("name"),
		Amount:    amount,
		Category:  r.FormValue("category"),
	}

	response, err := templates.GetFinanceFeedConfirm(expense)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
}

func FinanceFeedConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userId := session.Values["userId"].(string)
	expense := models.Expense{
		ExpenseId: r.FormValue("expenseId"),
		UserId:    userId,
		Name:      r.FormValue("name"),
		Amount:    amount,
		Category:  r.FormValue("category"),
	}

	// Validate input
	if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
		response, _ := templates.GetFinanceFeedServerResponse(
			userId,
			models.ServerResponse{
				Message: "Invalid input",
				ReturnEndpoint: FinanceFeedEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	success, err := dao.UpdateExpense(expense)
	if err != nil || !success {
		response, _ := templates.GetFinanceFeedServerResponse(
			userId,
			models.ServerResponse{
				Message: "Failed to update expense",
				ReturnEndpoint: FinanceFeedEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	response, _ := templates.GetFinanceFeedServerResponse(
		userId,
		models.ServerResponse{
			Message: "Expense updated successfully",
			ReturnEndpoint: FinanceFeedEndpoint,
	})
	w.Write([]byte(response))
}