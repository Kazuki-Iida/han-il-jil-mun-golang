package controllers

import (
	"boiler-plate/models"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

type UserController struct {
	Model *models.UserModel
}

func NewUserController(m *models.UserModel) *UserController {
	return &UserController{Model: m}
}

// func (h *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	// models/user.goのAll関数を使ってデータ取得
// 	users, err := h.Model.All()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(users)
// }

func (h *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/users/([0-9]+)/$")
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
	// models/user.goのGetUserById関数を使ってデータ取得
	user, err := h.Model.GetUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(*user)
}

func (h *UserController) SaveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// models/user.goのInsert関数を使ってデータ挿入
	id, err := h.Model.Insert(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id
	json.NewEncoder(w).Encode(user)
}

func (h *UserController) EditUser(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/users/([0-9]+)/$")
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
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// models/user.goのUpdate関数を使ってデータ挿入
	err = h.Model.Update(user.Name, user.Email, user.Password, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id
	json.NewEncoder(w).Encode(user)
}

func (h *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/users/([0-9]+)/$")
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
	// }

	// func (h UserController) GetLoginUser(w http.ResponseWriter, r *http.Request) {
	// 	token := auth0.GetJWT(r.Context())
	// 	fmt.Printf("jwt %+v\n", token)

	// 	// token.Claimsをjwt.MapClaimsへ変換
	// 	claims := token.Claims.(jwt.MapClaims)
	// 	// claimsの中にペイロードの情報が入っている
	// 	sub := claims["sub"].(string)

	// 	// userを取得する
	// 	user := getUser(sub)
	// 	if user == nil {
	// 		http.Error(w, "user not found", http.StatusNotFound)
	// 		return
	// 	}

	// 	// レスポンスを返す
	// 	res, err := json.Marshal(user)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Write(res)

	// // fmt.Println(sessions)
	// // // CookieからセッションIDを取得
	// // cookie, err := r.Cookie("session_id")
	// // if err != nil {
	// // 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// // 	return
	// // }
	// // sessionMutex.Lock()
	// // userId := sessions[cookie.Value].userId
	// // for id := range sessions {
	// // 	fmt.Println("loginHandler: Current session ID is ", id)
	// // }
	// // sessionMutex.Unlock()
	// // user, err := h.Model.GetUserById(userId)
	// // if err != nil {
	// // 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// // 	return
	// // }
	// // json.NewEncoder(w).Encode(*user)
}
