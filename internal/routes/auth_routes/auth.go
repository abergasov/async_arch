package auth_routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (ar *AuthAppRouter) generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (ar *AuthAppRouter) oauthGoogleLogin(c echo.Context) error {
	// generate cookie
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie(&http.Cookie{Name: "oauthstate", Value: state, Expires: time.Now().Add(365 * 24 * time.Hour)})

	u := ar.googleOAuth.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (ar *AuthAppRouter) oauthGoogleCallback(c echo.Context) error {
	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")
	if c.FormValue("state") != oauthState.Value {
		logger.Error("invalid oauth google state", errors.New("form and cookie mismatch"))
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	usr, err := ar.getUserDataFromGoogle(c.FormValue("code"))
	if err != nil {
		logger.Error("error get data from google", err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	jwt, err := ar.uService.Login(usr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]bool{"ok": false})
	}
	code, err := ar.exchanger.SetKey(jwt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]bool{"ok": false})
	}
	return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000?code="+code.String())
}

func (ar *AuthAppRouter) getUserDataFromGoogle(code string) (*entities.GoogleUser, error) {
	// Use code to get token and get user info from Google.
	token, err := ar.googleOAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	var uData entities.GoogleUser
	if err = json.Unmarshal(contents, &uData); err != nil {
		return nil, err
	}
	return &uData, nil
}

func (ar *AuthAppRouter) exchangeCode(c echo.Context) error {
	var u struct {
		Code uuid.UUID `json:"code"`
	}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]bool{"ok": false})
	}
	code, err := ar.exchanger.GetKey(u.Code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]bool{"ok": false})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true, "code": code})
}
