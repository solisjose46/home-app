package templates

import (
    "bytes"
    "errors"
    "html/template"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
    "github.com/solisjose46/pretty-print/debug"
)

type TemplateParser struct {
    dao *dao.Dao
}

func NewTemplateParser(dao *dao.Dao) *TemplateParser {
    return &TemplateParser{dao: dao}
}

func (parser *TemplateParser) GetLogin(serverResponse *models.ServerResponse) (*string, error) {
    debug.PrintInfo(parser.GetLogin, "Getting login template")

    htmlLogin := util.GetTmplPath(TmplLogin)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)

    if err != nil {
        debug.PrintError(parser.GetLogin, err)
        return nil, errors.New("error parsing login template")
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, models.Login{
        ServerResponse: serverResponse,
    })

    if err != nil {
        debug.PrintError(parser.GetLogin, err)
        return nil, errors.New("error executing login template")
    }

    tmplString := buf.String()

    debug.PrintSucc(parser.GetLogin, "returning home login template")
    return &tmplString, nil
}

func (parser *TemplateParser) PostLogin(username, password string) (*string, error) {
    debug.PrintInfo(parser.PostLogin, "Getting post login", "validate input")

	if username == "" || password == "" {
		return parser.GetLogin(
            &models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
                ReturnTarget: util.ReturnTargetBody,
            },
        )
	}

    debug.PrintInfo(parser.PostLogin, "auth attempt")

    valid, err := parser.dao.ValidateUser(username, password)

    if err != nil {
        debug.PrintError(parser.PostLogin, err)
        return nil, errors.New("error validating user")
    }

    if !valid {
        debug.PrintInfo(parser.PostLogin, "user not auth")
        return parser.GetLogin(
            &models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
                ReturnTarget: util.ReturnTargetBody,
            },
        )
    }

    debug.PrintSucc(parser.PostLogin, "authenticated!")
    return nil, nil
}

func (parser *TemplateParser) GetHome() (string, error) {
    debug.PrintInfo(parser.GetHome, "Getting home template")

    htmlHome := util.GetTmplPath(TmplHome)
    
    tmpl, err := template.ParseFiles(htmlHome)

    if err != nil {
        debug.PrintError(parser.GetHome, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplHome, nil)
    
    if err != nil {
        debug.PrintError(parser.GetHome, err)
        return "", err
    }

    debug.PrintSucc(parser.GetHome, "returning home template")
    return buf.String(), nil
}