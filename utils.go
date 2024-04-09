package cookie_header

import (
	"strings"
)

func fixHost(host string) string {
	last := strings.LastIndexByte(host, '.')
	if last == -1 {
		return host
	}
	last--
	for ; last >= 0; last-- {
		if host[last] == '.' {
			return host[last+1:]
		}
	}
	return host
}
