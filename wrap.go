package cookie_header

import (
	http "github.com/bogdanfinn/fhttp"
)

type CookieHeader interface {
	Load(res *http.Response, overrideDomain ...string)
	LoadOverrideEmpty(res *http.Response, overrideDomain string)
	LoadOverrideAll(res *http.Response, overrideDomain string)
	CreateHeader(domains ...string) string
	GetAllCookies(domains ...string) []cookie
	GetCookieValue(name, domain string) string
	AddCookie(name, value, domain, path string)
	DeleteCookie(name string, domains ...string)
}

// interface check
var _ CookieHeader = (*cookieHeader)(nil)
