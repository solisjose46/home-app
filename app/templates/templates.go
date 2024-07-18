package templates

import (
    "bytes"
    "html/template"
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

func GetLoginServerResponse(serverResponse models.ServerResponse) (string, error) {

    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        return "", err
    }

    loginData := models.Login{
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
    categories, err := dao.GetCategoriesForCurrentMonth() // todo implement this
    if err != nil {
        return "", err
    }

    // Prepare data for the template
    now := time.Now()
    month := now.Format("January 2006")
    finance := models.Finance{
        FinanceTrack: models.FinanceTrack{
            Month: month,
            Categories: categories,
        },
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

    financeTrack := {
        Categories: dao.GetCategoriesForCurrentMonth(),
        FinanceTrackConfirm: models.FinanceTrackConfirm{
            Expense: expense,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-track", financeTrack)
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

    financeTrack := {
        Categories: dao.GetCategoriesForCurrentMonth(),
        ServerResponse: serverResponse
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-track", financeTrack)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeed(userId string) (string, error) {
    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinance)
    if err != nil {
        return "", err
    }

    finance := models.Finance{
        FinanceFeed: models.FinanceFeed{
            Expenses: dao.GetExpensesForCurrentMonth(userId)
        }
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance", finance)
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

    financeFeed := models.FinanceFeed{
        Expenses: dao.GetExpensesForCurrentMonth(expense.UserId),
        FinanceFeedEdit: models.FinanceFeedEdit{
            Expense: expense,
        }
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed", financeFeed)
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

    financeFeed := models.FinanceFeed{
        Expenses: dao.GetExpensesForCurrentMonth(newExpense.UserId)
        FinanceFeedConfirm: models.FinanceFeedConfirm{
            OldExpense: oldExpense,
            NewExpense: newExpense,
        }
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed", financeFeed)
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

    financeFeed := models.FinanceFeed{
        Expenses: dao.GetExpensesForCurrentMonth(newExpense.UserId)
        ServerResponse: serverResponse
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "finance-feed", financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}
