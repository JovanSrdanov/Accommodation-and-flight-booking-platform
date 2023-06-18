package utils

import "strings"

func ParseWebBrowser(userAgent string) string {
	// Split User-Agent string by space
	parts := strings.Split(userAgent, " ")

	if len(parts) > 0 {
		// Extract the first part as the web browser
		return parts[0]
	}

	return "Unknown"
}
