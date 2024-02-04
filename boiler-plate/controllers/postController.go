package controllers

import (
	"boiler-plate/models"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

type PostController struct {
	Model *models.PostModel
}

func NewPostController(m *models.PostModel) *PostController {
	return &PostController{Model: m}
}

func (h *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// models/post.goのAll関数を使ってデータ取得
	posts, err := h.Model.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (h *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/posts/([0-9]+)/$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	// models/post.goのGetPostById関数を使ってデータ取得
	post, err := h.Model.GetPostById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(*post)
}

func (h *PostController) SavePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// models/post.goのInsert関数を使ってデータ挿入
	id, err := h.Model.Insert(post.Title, post.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = id
	json.NewEncoder(w).Encode(post)
}

func (h *PostController) EditPost(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/posts/([0-9]+)/$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// models/post.goのUpdate関数を使ってデータ挿入
	err = h.Model.Update(post.Title, post.Body, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = id
	json.NewEncoder(w).Encode(post)
}

func (h *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/posts/([0-9]+)/$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.GetPosts(w, r)
}
