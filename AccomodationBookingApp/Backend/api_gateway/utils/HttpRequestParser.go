package utils

import "github.com/mssola/user_agent"

func GetBrowserName(userAgent string) string {
	ua := user_agent.New(userAgent)
	browserName, _ := ua.Browser()

	return browserName
}
