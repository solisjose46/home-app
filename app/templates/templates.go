package templates

import (
    "bytes"
    "html/template"
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
)

const (
    tmplPath = "web/templates/"
    tmplLogin = "login"
    tmplHome = "home"
    tmplFinance = "finance"
    tmplFinanceTrack = "finance-track"
    tmplFinanceTrackConfirm = "finance-track-confirm"
    tmplFinanceFeed = "finance-feed"
    tmplFinanceFeedEdit = "finance-feed-edit"
    tmplFinanceFeedConfirm = "finance-feed-confirm"
    tmplServerResponse = "server-response"
    htmlExtension = ".html"
)

func GetLogin() (string, error) {
    htmlLogin := GetFilePath(tmplPath, tmplLogin, htmlExtension)

    tmpl, err := template.ParseFiles(htmlLogin)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplLogin, nil)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetLoginServerResponse(serverResponse models.ServerResponse) (string, error) {
    htmlLogin := GetFilePath(tmplPath, tmplLogin, htmlExtension)
    htmlServerResponse := GetFilePath(tmplPath, tmplServerResponse, htmlExtension)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        return "", err
    }

    loginData := models.Login{
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplLogin, loginData)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceTrack(userId string) (string, error) {
    categories, err := dao.GetCategoriesForCurrentMonth()
    if err != nil {
        return "", err
    }

    month := now.Format("January 2006")

    finance := models.Finance{
        FinanceTrack: models.FinanceTrack{
            Month: time.Now(),
            Categories: categories,
        },
    }

    htmlFinance := GetFilePath(tmplPath, tmplFinance, htmlExtension)
    htmlFinanceTrack := GetFilePath(tmplPath, tmplFinanceTrack, htmlExtension)

    tmpl, err := template.ParseFiles(htmlFinance, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinance, finance)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceTrackConfirm(expense models.Expense) (string, error) {
    categories, err := dao.GetCategoriesForCurrentMonth()
    if err != nil {
        return "", err
    }

    htmlFinanceTrack := GetFilePath(tmplPath, tmplFinanceTrack, htmlExtension)
    htmlFinanceTrackConfirm := GetFilePath(tmplPath, tmplFinanceTrackConfirm, htmlExtension)

    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlFinanceTrackConfirm)
    if err != nil {
        return "", err
    }

    financeTrack := models.FinanceTrack{
        Categories: categories,
        FinanceTrackConfirm: models.FinanceTrackConfirm{
            Expense: expense,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinanceTrack, financeTrack)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceTrackServerResponse(serverResponse models.ServerResponse) (string, error) {
    categories, err := dao.GetCategoriesForCurrentMonth()
    if err != nil {
        return "", err
    }

    htmlFinanceTrack := GetFilePath(tmplPath, tmplFinanceTrack, htmlExtension)
    htmlServerResponse := GetFilePath(tmplPath, tmplServerResponse, htmlExtension)

    tmpl, err := template.ParseFiles(htmlServerResponse, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    financeTrack := models.FinanceTrack{
        Categories: categories,
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinanceTrack, financeTrack)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeed(userId string) (string, error) {
    expenses := dao.GetExpensesForCurrentMonth(userId)

    if err != nil {
        return "", err
    }

    htmlFinance := GetFilePath(tmplPath, tmplFinance, htmlExtension)
    htmlFinanceFeed := GetFilePath(tmplPath, tmplFinanceFeed, htmlExtension)

    tmpl, err := template.ParseFiles(htmlFinance, htmlFinanceFeed)
    if err != nil {
        return "", err
    }

    finance := models.Finance{
        FinanceFeed: models.FinanceFeed{
            Expenses: expenses,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinance, finance)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedEdit(expense models.Expense) (string, error) {
    expenses := dao.GetExpensesForCurrentMonth(userId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := GetFilePath(tmplPath, tmplFinanceFeed, htmlExtension)
    htmlFinanceEdit := GetFilePath(tmplPath, tmplFinanceFeedEdit, htmlExtension)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceEdit)
    if err != nil {
        return "", err
    }

    financeFeed := models.FinanceFeed{
        Expenses: expenses,
        FinanceFeedEdit: models.FinanceFeedEdit{
            Expense: expense,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedConfirm(newExpense models.Expense) (string, error) {
    expenses := dao.GetExpensesForCurrentMonth(newExpense.UserId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := GetFilePath(tmplPath, tmplFinanceFeed, htmlExtension)
    htmlFinanceFeedConfirm := GetFilePath(tmplPath, tmplFinanceFeedConfirm, htmlExtension)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedConfirm)
    if err != nil {
        return "", err
    }

    oldExpense := dao.GetExpense(newExpense.UserId)

    financeFeed := models.FinanceFeed{
        Expenses: expenses,
        FinanceFeedConfirm: models.FinanceFeedConfirm{
            OldExpense: oldExpense,
            NewExpense: newExpense,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedServerResponse(serverResponse models.ServerResponse) (string, error) {
    expenses := dao.GetExpensesForCurrentMonth(newExpense.UserId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := GetFilePath(tmplPath, tmplFinanceFeed, htmlExtension)
    htmlServerResponse := GetFilePath(tmplPath, tmplServerResponse, htmlExtension)

    financeFeed := models.FinanceFeed{
        Expenses: expenses,
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, tmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}
