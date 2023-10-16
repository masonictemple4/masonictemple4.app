package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/masonictemple4/masonictemple4.app/internal/repository"
	"github.com/masonictemple4/masonictemple4.app/utils"
)

// Each handler file should begin with a `TypeRoutes() *mux.Router` function
// Responsible for registering all the routes required for those handlers.
//
// Once complete you can call this in the official registerRoutes function
// defined inside of the server.go file.
//
// Routes:
//
//	"/blog": s.PostBlogHandler [http.MethodPost]
//	"/blog": s.QueryBlogsHandler [http.MethodGet]
//	"/blog/{slug}": s.BlogDetailHandler [http.MethodGet]
//	"/blog/{id}": s.UpdateBlogHandler [http.MethodPut]
//	"/blog/{id}": s.DeleteBlogHandler [http.MethodDelete]
//
// TODO: enable the previous routes once security is enabled.
func (s *Server) BlogRoutes() *mux.Router {
	router := mux.NewRouter().PathPrefix("/blog").Subrouter()
	// router.HandleFunc("", s.PostBlogHandler).Methods(http.MethodPost)
	router.HandleFunc("", s.BlogFeedHandler).Methods(http.MethodGet)
	router.HandleFunc("/{slug}", s.BlogDetailHandler).Methods(http.MethodGet)
	// router.HandleFunc("/{id}", s.UpdateBlogHandler).Methods(http.MethodPut)
	// router.HandleFunc("/{id}", s.DeleteBlogHandler).Methods(http.MethodDelete)
	return router
}

// For now just ignore the user associated with the media and the blog, however,
// we can make middleware that extends the http.Request object and sets a user value
// on it so we can later use it.
func (s *Server) PostBlogHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Once middleware is in use we can get the user from the request and add it
	// to the authors list.
	var reqBody dtos.PostInput
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid body try again.")
		return
	}

	var blog models.Post
	err = utils.Convert(reqBody, &blog)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem generating the blog from the post body.")
		return
	}

	blog.State = models.BlogStateDraft

	err = s.DB.Create(&blog).Error
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem creating the blog record.")
		return
	}

	// 	for tag := range reqBody.Tags {
	// 		var t models.Tag
	// 		s.DB.FirstOrCreate(&t, models.Tag{Name: reqBody.Tags[tag]})
	// 		err := s.DB.Model(&blog).Association("Tags").Append(&t).Error
	// 		if err != nil {
	// 			ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem updating the tags.")
	// 			return
	// 		}
	// 	}

	var returnObj dtos.PostReturn
	err = utils.Convert(blog, &returnObj)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem generating the return object.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnObj)
}

func (s *Server) BlogFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add query params.
	var queryRepo models.Post
	queryOpts := &repository.RepositoryOpts{
		Preloads: map[string][]string{
			"Authors":  make([]string, 0),
			"Comments": make([]string, 0),
			"Tags":     make([]string, 0),
		},
	}

	var results []models.Post

	err := queryRepo.All(s.DB, queryOpts, &results)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	feedReturn := make([]dtos.PostReturn, 0)
	err = utils.Convert(results, &feedReturn)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&feedReturn)

}

func (s *Server) BlogDetailHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	slug, ok := vars["slug"]
	if !ok {
		ErrorHandler(w, r, http.StatusNotFound, "")
		return
	}

	queryOpts := &repository.RepositoryOpts{
		Preloads: map[string][]string{
			"Authors":  make([]string, 0),
			"Comments": make([]string, 0),
			"Tags":     make([]string, 0),
		},
	}

	var result models.Post
	err := result.FindBySlug(s.DB, slug, queryOpts)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	var detailRet dtos.PostReturn
	err = utils.Convert(result, &detailRet)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("blogdetail: utils.convert results to detail return %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&detailRet)

}

func (s *Server) UpdateBlogHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody dtos.UpdatePostInput
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid body try again.")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Post
	// TODO: When auth is handled we'll check the user matches here.
	err = s.DB.First(&blog, id).Error
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem finding the blog record.")
		return
	}

	updateMap := map[string]any{
		"title":      reqBody.Title,
		"ContentUrl": reqBody.ContentUrl,
		"subtitle":   reqBody.Subtitle,
	}

	if reqBody.State != "" && models.ValidBlogState(reqBody.State) {
		updateMap["state"] = reqBody.State
	}

	err = s.DB.Model(&blog).Updates(updateMap).Error
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem updating the base blog record.")
		return
	}

	for tag := range reqBody.Tags {
		var t models.Tag
		s.DB.FirstOrCreate(&t, models.Tag{Name: reqBody.Tags[tag]})
		err := s.DB.Model(&blog).Association("Tags").Append(&t).Error
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem updating the tags.")
			return
		}
	}

	var returnObj dtos.PostReturn
	err = utils.Convert(blog, &returnObj)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem generating the return object.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnObj)

}

func (s *Server) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.DB.Delete(&models.Post{}, id).Error
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "There was a problem deleting the blog record.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos.SuccessResponse{Message: "Blog deleted successfully."})
}
