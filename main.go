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
	e.POST("/token/only", tokenOnly)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// tokenGet Get Keycloak token for my GitHub account
func tokenGet(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, getTokenResponse(), " ")
}

// tokenVerify Verify Keycloak token
func tokenVerify(c echo.Context) error {
	bearer := "Bearer " + token(c)
	req, _ := http.NewRequest("POST", endpoint("userinfo"), nil)
	req.Header.Add("Authorization", bearer)

	return c.JSONPretty(http.StatusOK, send(req), " ")
}

// tokenOnly Get token only from response
func tokenOnly(c echo.Context) error {
	jsonData := getTokenResponse()
	response := TokenResponse{}
	result := response.parseTokenResponse(jsonData)

	return c.String(http.StatusOK, result.AccessToken)
}

// token Get token from "Authorization" header
func token(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	slice := strings.Split(bearer, " ")

	return slice[len(slice)-1]
}

// endpoint Get Keycloak endpoint
func endpoint(action string) string {
	return fmt.Sprintf("%s/protocol/openid-connect/%s", os.Getenv("DOMAIN_WITH_REALM"), action)
}

// getTokenResponse Get token from token-endpoint response
func getTokenResponse() string {
	data := url.Values{}

	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("username", os.Getenv("USERNAME"))
	data.Set("password", os.Getenv("PASSWORD"))
	data.Set("grant_type", "password")

	req, _ := http.NewRequest("POST", endpoint("token"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return send(req)
}

// send Send request to Keycloak
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
