package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	result := TokenResponse{}
	_ = json.Unmarshal(getTokenResponse(), &result)

	return c.JSONPretty(http.StatusOK, result, " ")
}

// tokenVerify Verify Keycloak token
func tokenVerify(c echo.Context) error {
	bearer := "Bearer " + token(c)
	req, _ := http.NewRequest("POST", endpoint("userinfo"), nil)
	req.Header.Add("Authorization", bearer)

	result := VerifyResponse{}
	_ = json.Unmarshal(send(req), &result)

	return c.JSONPretty(http.StatusOK, result, " ")
}

// tokenOnly Get token only from response
func tokenOnly(c echo.Context) error {
	result := TokenResponse{}
	_ = json.Unmarshal(getTokenResponse(), &result)

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
func getTokenResponse() []byte {
	data := url.Values{}

	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("username", os.Getenv("USERNAME"))
	data.Set("password", os.Getenv("PASSWORD"))
	data.Set("grant_type", "password")

	req, _ := http.NewRequest("POST", endpoint("token"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return send(req)
}

// send request to Keycloak
func send(req *http.Request) []byte {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error on response: %v", err)
	}

	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			log.Fatalf("Something wrong: %v", e)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error while reading the response bytes: %v", err)
	}

	return body
}
