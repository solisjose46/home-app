package templates

import (
    "fmt"
    "bytes"
    "html/template"
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
)

const (
    TmplPath = "web/templates/"
    TmplLogin = "login"
    TmplHome = "home"
    TmplFinance = "finance"
    TmplFinanceTrack = "finance-track"
    TmplFinanceTrackConfirm = "finance-track-confirm"
    TmplFinanceFeed = "finance-feed"
    TmplFinanceFeedEdit = "finance-feed-edit"
    TmplFinanceFeedConfirm = "finance-feed-confirm"
    TmplServerResponse = "server-response"
    HtmlExtension = ".html"
)

func GetLogin() (string, error) {
    htmlLogin := util.GetFilePath(TmplPath, TmplLogin, HtmlExtension)
    fmt.Println(htmlLogin)
    tmpl, err := template.ParseFiles(htmlLogin)
    if err != nil {
        fmt.Println("error parsing template", htmlLogin)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, models.Login{
        ServerResponse: models.ServerResponse{},
    })
    if err != nil {
        fmt.Println("template parse error")
        return "", err
    }

    return buf.String(), nil
}

func GetLoginServerResponse(serverResponse models.ServerResponse) (string, error) {
    htmlLogin := util.GetFilePath(TmplPath, TmplLogin, HtmlExtension)
    htmlServerResponse := util.GetFilePath(TmplPath, TmplServerResponse, HtmlExtension)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        return "", err
    }

    loginData := models.Login{
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, loginData)
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

    month := time.Now().Month().String()

    finance := models.Finance{
        FinanceTrack: models.FinanceTrack{
            Month: month,
            Categories: categories,
        },
    }

    htmlFinance := util.GetFilePath(TmplPath, TmplFinance, HtmlExtension)
    htmlFinanceTrack := util.GetFilePath(TmplPath, TmplFinanceTrack, HtmlExtension)

    tmpl, err := template.ParseFiles(htmlFinance, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, finance)
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

    htmlFinanceTrack := util.GetFilePath(TmplPath, TmplFinanceTrack, HtmlExtension)
    htmlFinanceTrackConfirm := util.GetFilePath(TmplPath, TmplFinanceTrackConfirm, HtmlExtension)

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
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
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

    htmlFinanceTrack := util.GetFilePath(TmplPath, TmplFinanceTrack, HtmlExtension)
    htmlServerResponse := util.GetFilePath(TmplPath, TmplServerResponse, HtmlExtension)

    tmpl, err := template.ParseFiles(htmlServerResponse, htmlFinanceTrack)
    if err != nil {
        return "", err
    }

    financeTrack := models.FinanceTrack{
        Categories: categories,
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeed(userId string) (string, error) {
    expenses, err := dao.GetExpensesForCurrentMonth(userId)

    if err != nil {
        return "", err
    }

    htmlFinance := util.GetFilePath(TmplPath, TmplFinance, HtmlExtension)
    htmlFinanceFeed := util.GetFilePath(TmplPath, TmplFinanceFeed, HtmlExtension)

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
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, finance)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedEdit(expense models.Expense) (string, error) {
    expenses, err := dao.GetExpensesForCurrentMonth(expense.UserId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := util.GetFilePath(TmplPath, TmplFinanceFeed, HtmlExtension)
    htmlFinanceEdit := util.GetFilePath(TmplPath, TmplFinanceFeedEdit, HtmlExtension)

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
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedConfirm(newExpense models.Expense) (string, error) {
    expenses, err := dao.GetExpensesForCurrentMonth(newExpense.UserId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := util.GetFilePath(TmplPath, TmplFinanceFeed, HtmlExtension)
    htmlFinanceFeedConfirm := util.GetFilePath(TmplPath, TmplFinanceFeedConfirm, HtmlExtension)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedConfirm)
    if err != nil {
        return "", err
    }

    oldExpense, err := dao.GetExpense(newExpense.UserId)
    if err != nil {
        return "", err
    }

    financeFeed := models.FinanceFeed{
        Expenses: expenses,
        FinanceFeedConfirm: models.FinanceFeedConfirm{
            OldExpense: oldExpense,
            NewExpense: newExpense,
        },
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}

func GetFinanceFeedServerResponse(userId string, serverResponse models.ServerResponse) (string, error) {
    expenses, err := dao.GetExpensesForCurrentMonth(userId)

    if err != nil {
        return "", err
    }

    htmlFinanceFeed := util.GetFilePath(TmplPath, TmplFinanceFeed, HtmlExtension)
    htmlServerResponse := util.GetFilePath(TmplPath, TmplServerResponse, HtmlExtension)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlServerResponse)
    if err != nil {
        return "", err
    }

    financeFeed := models.FinanceFeed{
        Expenses: expenses,
        ServerResponse: serverResponse,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}