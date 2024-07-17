package templates

import (
    "bytes"
    "html/template"
    "path/filepath"
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