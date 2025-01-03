package user

import (
	"errors"
	"gluttony/internal/security"
	"net/http"
)

func LogoutHandler(s *Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.Logout(r.Context()); err != nil {
			if errors.Is(err, security.ErrSessionNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, security.NewInvalidateCookie())
	})
}
