package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, e.Routes(), " ")
	})

	e.POST("/token/get", tokenGet)
	e.POST("/token/verify", tokenVerify)
	e.POST("/users/create", usersCreate)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Get Keycloak token for my GitHub account
func tokenGet(c echo.Context) error {
	data := url.Values{}

	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("username", os.Getenv("USERNAME"))
	data.Set("password", os.Getenv("PASSWORD"))
	data.Set("grant_type", "password")

	req, _ := http.NewRequest("POST", endpoint("token"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.String(http.StatusOK, send(req))
}

// Verify Keycloak token
func tokenVerify(c echo.Context) error {
	bearer := "Bearer " + token(c)
	req, _ := http.NewRequest("POST", endpoint("userinfo"), nil)
	req.Header.Add("Authorization", bearer)

	return c.String(http.StatusOK, send(req))
}

// Create new user TODO
func usersCreate(c echo.Context) error {
	fullUrl := os.Getenv("DOMAIN_WITH_REALM")
	ssoUrl := SsoUrl{}
	sso := ssoUrl.parseUrl(fullUrl)
	bearer := "Bearer " + token(c)

	form := c.Request().Form
	data := url.Values{}

	data.Set("firstName", form.Get("firstName"))
	data.Set("lastName", form.Get("lastName"))
	data.Set("email", form.Get("email"))
	data.Set("enabled", form.Get("enabled"))
	data.Set("username", form.Get("username"))

	req, _ := http.NewRequest("POST", sso.UrlAdmin, strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.String(http.StatusOK, "")
}

// Get token from "Authorization" header
func token(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	slice := strings.Split(bearer, " ")

	return slice[len(slice)-1]
}

// Get Keycloak endpoint
func endpoint(action string) string {
	return fmt.Sprintf("%s/protocol/openid-connect/%s", os.Getenv("DOMAIN_WITH_REALM"), action)
}

// Send request to Keycloak
func send(req *http.Request) string {
	result := "..."
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		result = "Error on response: " + err.Error()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			result = "Something wrong: " + err.Error()
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		result = "Error while reading the response bytes: " + err.Error()
	}

	result = string(body)

	return result
}
