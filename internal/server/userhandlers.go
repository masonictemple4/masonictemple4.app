package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func (s *Server) UserRoutes() *mux.Router {
	router := mux.NewRouter().PathPrefix("/user").Subrouter()
	router.HandleFunc("/google", s.GoogleAuth).Methods(http.MethodPost)
	router.HandleFunc("/github", s.GithubAuth).Methods(http.MethodPost)
	router.HandleFunc("/logout", s.LogoutHandler).Methods(http.MethodPost)
	return router
}

func (s *Server) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var body dtos.GoogleAuthInput
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Setup a local user if one does not exist and auth them with a JWT.
	payload, err := idtoken.Validate(context.Background(), body.IdToken, os.Getenv("GOOGLE_ID"))
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (s *Server) GithubAuth(w http.ResponseWriter, r *http.Request) {

	var body dtos.GithubAuthInput
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Setup a local user if one does not exist and auth them with a JWT.
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: body.Token})
	tc := oauth2.NewClient(r.Context(), ts)

	client := github.NewClient(tc)

	user, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}
