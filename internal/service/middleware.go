package service

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookie
		cookie, err := app.cookieStore.Get(r, app.Name)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if is_authed is true and authed_at is within the last 30 days
		isAuthed, ok := cookie.Values["is_authed"].(bool)
		fmt.Println(isAuthed)
		authedAt, okAt := cookie.Values["authed_at"].(int64)
		fmt.Println(authedAt)
		if !ok || !okAt || !isAuthed || !isAuthedWithinLast30Days(authedAt) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, okUserID := cookie.Values["user_id"].(string)
		fmt.Println(userID)
		if okUserID {
			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)
		}

		// If the checks pass, call the next handler
		next.ServeHTTP(w, r)
	})
}

func isAuthedWithinLast30Days(authedAt int64) bool {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Unix()
	return authedAt > thirtyDaysAgo
}
