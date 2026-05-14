package middleware

import "strings"

func cleanToken(dirtToken string) string {
	return strings.TrimPrefix(dirtToken, TokenPrefix)
}
