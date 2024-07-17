package templates

import (
    "bytes"
    "html/template"
    "path/filepath"
    "time"
    "home-app/app/models"
    "home-app/app/dao"
)

const (
	htmlFinanceFeedConfirm = "web/templates/finance-feed-confirm.html"
	htmlFinanceFeedEdit = "web/templates/finance-feed-edit.html"
	htmlFinanceFeed = "web/templates/finance-feed.html"
	htmlFinanceTrackConfirm = "web/templates/finance-track-confirm.html"
	htmlFinanceTrack = "web/templates/finance-track.html"
	htmlFinance = "web/templates/finance.html"
	htmlHome = "web/templates/home.html"
	htmlLogin = "web/templates/login.html"
	htmlServerResponse = "web/templates/server-response.html"
	// todo make template names consts

	// return endpoints const should go with handlers.go
)

func GetLoginServerResponse(serverResponse ServerResponse) (string, error) {

    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        return "", err
    }

    loginData := Login{
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "login", loginData)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceTrack retrieves current month expenses and embeds them in FinanceTrack template.
func GetFinanceTrack(userId string) (string, error) {
    // Get current month expenses
    expenses, err := dao.GetExpensesForCurrentMonth(userId)
    if err != nil {
        return "", err
    }

    // Prepare data for the template
    now := time.Now()
    month := now.Format("January 2006")
    financeTrack := models.FinanceTrack{
        Month: month,
        Expenses: expenses,
    }

    // Parse templates
    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlFinance)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-track", financeTrack)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceTrackConfirm embeds an expense in FinanceTrackConfirm template.
func GetFinanceTrackConfirm(expense models.Expense) (string, error) {
    tmpl, err := template.ParseFiles(htmlFinanceTrackConfirm, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    financeTrackConfirm := models.FinanceTrackConfirm{
        Expense: expense,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-track-confirm", financeTrackConfirm)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceTrackServerResponse embeds a server response in ServerResponse template and returns it in FinanceTrack template.
func GetFinanceTrackServerResponse(serverResponse models.ServerResponse) (string, error) {
    tmpl, err := template.ParseFiles(htmlServerResponse, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "server-response", serverResponse)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceFeed embeds FinanceFeed template in Finance template and returns it.
func GetFinanceFeed() (string, error) {
    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinance)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed", nil)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceFeedEdit embeds an expense in FinanceFeedEdit template and returns it in FinanceFeed template.
func GetFinanceFeedEdit(expense models.Expense) (string, error) {
    tmpl, err := template.ParseFiles(htmlFinanceFeedEdit, htmlFinanceFeed)
    if err != nil {
        return "", err
    }

    financeFeedEdit := models.FinanceFeedEdit{
        Expense: expense,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed-edit", financeFeedEdit)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedConfirm(newExpense models.Expense) (string, error) {
    tmpl, err := template.ParseFiles(htmlFinanceFeedConfirm, htmlFinanceFeed)
    if err != nil {
        return "", err
    }

	oldExpense := dao.GetExpense(newExpense.ExpenseId)

    financeFeedConfirm := models.FinanceFeedConfirm{
        OldExpense: oldExpense,
        NewExpense: newExpense,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed-confirm", financeFeedConfirm)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

// GetFinanceFeedServerResponse embeds a server response in ServerResponse template and returns it in FinanceFeed template.
func GetFinanceFeedServerResponse(serverResponse models.ServerResponse) (string, error) {
    tmpl, err := template.ParseFiles(htmlServerResponse, htmlFinanceFeed)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "server-response", serverResponse)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}
