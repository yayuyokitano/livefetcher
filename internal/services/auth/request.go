package auth

import (
	"context"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) (user AuthUser) {
	ctx := context.Background()
	c, err := r.Cookie("authToken")
	if err != nil {
		return
	}

	user, err = verifyAuthToken(c.Value)
	if err == nil {
		return
	}

	c, err = r.Cookie("refreshToken")
	if err != nil {
		return
	}
	authToken, refreshToken, err := RefreshSession(ctx, c.Value)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    authToken,
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	user, _ = verifyAuthToken(authToken)
	return
}
