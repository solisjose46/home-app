package handlers

import (
	"fmt"
	"net/http"
	"home-app/app/dao"
	"home-app/app/models"
	"home-app/app/templates"
	"home-app/app/util"
	"strconv"
	"github.com/gorilla/sessions"
)

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
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] != nil {
		http.Redirect(w, r, HomeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
        response, err := templates.GetLogin()
		fmt.Println("ERROR here")
        if err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        w.Write([]byte(response))
        return
    }

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		loginServerResponse, _ := templates.GetLoginServerResponse(models.ServerResponse{
			Message: "Invalid input",
			ReturnEndpoint: LoginEndpoint,
		})
		w.Write([]byte(loginServerResponse))
		return
	}

	valid, err := dao.ValidateUser(username, password)
	if err != nil || !valid {
		response, _ := templates.GetLoginServerResponse(
			models.ServerResponse{
				Message: "Bad login",
				ReturnEndpoint: LoginEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	userId, err := dao.GetUserId(username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["username"] = username
	session.Values["userId"] = userId
	session.Save(r, w)

	http.Redirect(w, r, HomeEndpoint, http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	htmlHome := util.GetFilePath(templates.TmplPath, templates.TmplHome, templates.HtmlExtension)

	http.ServeFile(w, r, htmlHome)
}

func FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, FinanceTrackEndpoint, http.StatusSeeOther)
}

func FinanceTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values["userId"].(string)

	if r.Method == http.MethodGet {
		response, err := templates.GetFinanceTrack(userId)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(response))
		return
	}

	// Handle POST method: gather expense data from form and return confirmation page

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		UserId:   userId,
		Name:     r.FormValue("name"),
		Amount:   amount,
		Category: r.FormValue("category"),
	}

	response, err := templates.GetFinanceTrackConfirm(expense)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
}

func FinanceTrackConfirmHandler(w http.ResponseWriter, r *http.Request) {
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

	expense := models.Expense{
		UserId:   session.Values["userId"].(string),
		Name:     r.FormValue("name"),
		Amount:   amount,
		Category: r.FormValue("category"),
	}

	// Validate input
	if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
		response, _ := templates.GetFinanceTrackServerResponse(models.ServerResponse{
			Message: "Invalid input",
			ReturnEndpoint: FinanceTrackEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	success, err := dao.AddExpense(expense)
	if err != nil || !success {
		response, _ := templates.GetFinanceTrackServerResponse(models.ServerResponse{
			Message: "Failed to add expense",
			ReturnEndpoint: FinanceTrackEndpoint,
		})
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