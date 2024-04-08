package cookie_header

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	http "github.com/bogdanfinn/fhttp"
)

const (
	GLOBAL = "!#GLOBAL#!"
)

type cookieHeader struct {
	jar map[string]map[string]cookie
	wg  *sync.RWMutex
}

type cookie struct {
	name   string
	value  string
	domain string
}

// create new cookieHeader
func New() cookieHeader {
	newJar := cookieHeader{}
	newJar.jar = make(map[string]map[string]cookie)
	newJar.wg = &sync.RWMutex{}
	return newJar
}

func (ch *cookieHeader) LoadFromResponse(res *http.Response, overrideDomain ...string) {
	ch.wg.Lock()
	defer ch.wg.Unlock()

	for _, c := range res.Cookies() {

		if c.Domain == "" {
			c.Domain = "." + fixHost(res.Request.Host)
		}

		if c.MaxAge < 0 {
			ch.deletecookie(c.Name, c.Domain)
			continue
		}

		if len(overrideDomain) > 0 {
			c.Domain = overrideDomain[0]
		}

		ch.addcookie(c.Name, c.Value, c.Domain)
	}
}

func (ch *cookieHeader) CreateHeader(domains ...string) string {
	ch.wg.RLock()
	defer ch.wg.RUnlock()

	arr := []string{}

	for k, d := range domains {
		domains[k] = fixHost(d)
	}

	for domain, cks := range ch.jar {
		if len(domains) == 0 || domain == GLOBAL || slices.Contains(domains, domain) {
			for _, cookie := range cks {
				arr = append(arr, cookie.toString())
			}
		}
	}
	return strings.Join(arr, "; ")
}

func (ch *cookieHeader) GetAllCookies(domains ...string) []cookie {
	ch.wg.RLock()
	defer ch.wg.RUnlock()

	for k, d := range domains {
		domains[k] = fixHost(d)
	}

	allCks := []cookie{}

	for domain, cks := range ch.jar {
		if len(domains) == 0 || domain == GLOBAL || slices.Contains(domains, domain) {
			for _, cookie := range cks {
				allCks = append(allCks, cookie)
			}
		}
	}

	return allCks
}

func (ch *cookieHeader) GetCookieValue(name, domain string) string {
	ch.wg.RLock()
	defer ch.wg.RUnlock()

	domain = fixHost(domain)

	if cks, ok := ch.jar[domain]; ok {
		if cookie, ok := cks[name]; ok {
			return cookie.value
		}
	}
	return ""
}

func (ch *cookieHeader) AddCookie(name, value, domain string) {
	ch.wg.Lock()
	defer ch.wg.Unlock()

	ch.addcookie(name, value, domain)
}

func (ch *cookieHeader) DeleteCookie(name string, domains ...string) {
	ch.wg.Lock()
	defer ch.wg.Unlock()

	for k, d := range domains {
		domains[k] = fixHost(d)
	}

	for domain := range ch.jar {
		if len(domains) == 0 || domain == GLOBAL || slices.Contains(domains, domain) {
			ch.deletecookie(name, domain)
		}
	}
}

// private func: add single cookie
func (ch *cookieHeader) addcookie(name, value, domain string) {
	ogDomain := domain
	domain = fixHost(domain)
	if domain == "" {
		domain = GLOBAL
	}
	cks, ok := ch.jar[domain]
	if !ok {
		cks = make(map[string]cookie)
		ch.jar[domain] = cks
	}
	cks[name] = cookie{
		name:   name,
		value:  value,
		domain: ogDomain,
	}
}

// private func: delete single cookie
func (ch *cookieHeader) deletecookie(name, domain string) {
	domain = fixHost(domain)
	if cks, ok := ch.jar[domain]; ok {
		delete(cks, name)
	}
}

// private func: cookie to string format
func (cookie *cookie) toString() string {
	return fmt.Sprintf("%s=%s", cookie.name, cookie.value)
}

// print all cookies for debug stuff
func (ch *cookieHeader) DebugPrintAllCookies() {
	ch.wg.RLock()
	defer ch.wg.RUnlock()

	for domain, cks := range ch.jar {
		fmt.Println("[" + domain + "]")
		for _, cookie := range cks {
			fmt.Println(cookie.toString())
		}
	}
}
