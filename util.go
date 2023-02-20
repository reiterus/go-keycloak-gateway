package main

import (
	"fmt"
	"net/url"
	"strings"
)

// SsoUrl Structure to parse url TODO
type SsoUrl struct {
	Scheme      string `json:"scheme"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Realm       string `json:"realm"`
	UrlOriginal string `json:"url_original"`
	UrlNoRealm  string `json:"url_no_realm"`
	UrlAdmin    string `json:"url_admin"`
}

// parseUrl Parse url with realm to structure TODO
func (su SsoUrl) parseUrl(domain string) SsoUrl {

	item, _ := url.Parse(domain)
	slice := strings.Split(domain, "/")
	realm := slice[len(slice)-1]
	format := "%s://%s/auth/admin/realms/%s/users"

	return SsoUrl{
		Scheme:      item.Scheme,
		Host:        item.Host,
		Path:        item.Path,
		Realm:       slice[len(slice)-1],
		UrlOriginal: domain,
		UrlNoRealm:  strings.Replace(domain, realm, "", 1),
		UrlAdmin:    fmt.Sprintf(format, item.Scheme, item.Host, realm),
	}
}
