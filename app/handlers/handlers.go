package handlers

import (
	"encoding/json"
	"net/http"
	"home-app/app/dao"
	"home-app/app/models"
	"home-app/app/templates"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

const (
	loginEndpoint					= "/login"
	homeEndpoint					= "/home"
	financeEndpoint					= "/finance"
	financeTrackEndpoint			= "/finance/track"
	financeTrackConfirmEndpoint		= "/finance/track/confirm"
	financeFeedEndpoint				= "/finance/feed"
	financeFeedConfirmPoint			= "/finance/feed/confirm"
	financeFeedEditEndpoint			= "/finance/feed/edit"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] != nil {
		http.Redirect(w, r, homeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
        response, err := templates.GetLogin()
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
			ReturnEndpoint: loginEndpoint,
		})
		w.Write([]byte(loginServerResponse))
		return
	}

	valid, err := dao.ValidateUser(username, password)
	if err != nil || !valid {
		response, _ := templates.GetLoginServerResponse(models.ServerResponse{
			Message: "Bad login",
			ReturnEndpoint: loginEndpoint
		})
		w.Write([]byte(response))
		return
	}

	session.Values["username"] = username
	session.Values["userId"] = dao.GetUserId(username)
	session.Save(r, w)

	http.Redirect(w, r, homeEndpoint, http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, htmlHome)
}

func FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, financeTrackEndpoint, http.StatusSeeOther)
}

func FinanceTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	if session.Values["username"] == nil {
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
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
	expense := models.Expense{
		UserId:   userId,
		Name:     r.FormValue("name"),
		Amount:   r.FormValue("amount"),
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
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	expense := models.Expense{
		UserId:   session.Values["userId"].(string),
		Name:     r.FormValue("name"),
		Amount:   r.FormValue("amount"),
		Category: r.FormValue("category"),
	}

	// Validate input
	if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
		response, _ := templates.GetFinanceTrackServerResponse(models.ServerResponse{
			Message: "Invalid input",
			ReturnEndpoint: financeTrackEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	success, err := dao.AddExpense(expense)
	if err != nil || !success {
		response, _ := templates.GetFinanceTrackServerResponse(models.ServerResponse{
			Message: "Failed to add expense",
			ReturnEndpoint: financeTrackEndpoint,
		})
		w.Write([]byte(response))
		return
	}

	response, _ := templates.GetFinanceTrackServerResponse(models.ServerResponse{
		Message: "Expense added successfully",
		ReturnEndpoint: financeTrackEndpoint
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
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	response, err := templates.GetFinanceFeed()
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
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
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
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue("expenseId"),
		UserId:    session.Values["userId"].(string),
		Name:      r.FormValue("name"),
		Amount:    r.FormValue("amount"),
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
		http.Redirect(w, r, loginEndpoint, http.StatusSeeOther)
		return
	}

	expense := models.Expense{
		ExpenseId: r.FormValue("expenseId"),
		UserId:    session.Values["userId"].(string),
		Name:      r.FormValue("name"),
		Amount:    r.FormValue("amount"),
		Category:  r.FormValue("category"),
	}

	// Validate input
	if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
		response, _ := templates.GetFinanceFeedServerResponse(models.ServerResponse{
			Message: "Invalid input",
			ReturnEndpoint: financeFeedEndpoint
		})
		w.Write([]byte(response))
		return
	}

	success, err := dao.UpdateExpense(expense)
	if err != nil || !success {
		response, _ := templates.GetFinanceFeedServerResponse(models.ServerResponse{
			Message: "Failed to update expense",
			ReturnEndpoint: financeFeedEndpoint
		})
		w.Write([]byte(response))
		return
	}

	response, _ := templates.GetFinanceFeedServerResponse(models.ServerResponse{
		Message: "Expense updated successfully",
		ReturnEndpoint: financeFeedEndpoint,
	})
	w.Write([]byte(response))
}