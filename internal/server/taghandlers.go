package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"github.com/masonictemple4/masonictemple4.app/utils"
)

func (s *Server) TagRoutes() *mux.Router {
	router := mux.NewRouter().PathPrefix("/tags").Subrouter()
	router.HandleFunc("", s.GetTags)
	// router.HandleFunc("", s.PostTagHandler)
	// router.HandleFunc("/{id}", s.TagDetailHandler)
	// router.HandleFunc("/{id}", s.UpdateTagHandler)
	// router.HandleFunc("/{id}", s.DeleteTagHandler)
	return router
}

func (s *Server) GetTags(w http.ResponseWriter, r *http.Request) {

	var tags []models.Tag

	query := r.URL.Query().Get("name")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var err error
	limit := -1
	offset := -1

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = -1
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			offset = -1
		}
	}

	if query != "" {
		err = s.DB.Where("name LIKE ?", "%"+query+"%").Limit(limit).Offset(offset).Find(&tags).Error
	} else {
		err = s.DB.Limit(limit).Offset(offset).Find(&tags).Error
	}

	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Adds about ~30ms of response time.
	var response []dtos.TagReturn
	err = utils.Convert(tags, &response)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
