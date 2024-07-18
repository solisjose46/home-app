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
)

func GetLogin() (string, error) {
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "login", nil)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

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

func GetFinanceTrack(userId string) (string, error) {
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
