package cookie_header

import (
	"fmt"
	"strings"
)

func fixHost(host string) string {
	host = strings.TrimLeft(host, ".")
	hParts := strings.Split(host, ".")

	if len(hParts) < 2 {
		return host
	}
	return fmt.Sprintf("%s.%s", hParts[len(hParts)-2], hParts[len(hParts)-1])
}
