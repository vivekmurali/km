package main

import "net/http"

func (s *app) UnloggedInRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session")

		authenticated, ok := session.Values["authenticated"].(bool)

		if !authenticated || !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *app) LoggedInRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session")

		authenticated, ok := session.Values["authenticated"].(bool)

		if authenticated && ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		next.ServeHTTP(w, r)
	})
}
