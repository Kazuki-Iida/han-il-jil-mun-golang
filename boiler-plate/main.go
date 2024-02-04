package main

import (
	"boiler-plate/controllers"
	"boiler-plate/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, session_id")

		// preflightリクエストへの応答
		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func main() {
	db, err := sql.Open("mysql", "dbuser:password@tcp(go_blog_mysql:3306)/go_blog_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userModel := models.NewUserModel(db)
	userHandler := controllers.NewUserController(userModel)
	postModel := models.NewPostModel(db)
	postHandler := controllers.NewPostController(postModel)

	http.HandleFunc("/login", enableCORS(userHandler.LoginHandler))
	http.HandleFunc("/register", enableCORS(userHandler.RegisterHandler))
	http.HandleFunc("/dashboard", enableCORS(controllers.DashboardHandler))
	http.HandleFunc("/logout", enableCORS(userHandler.LogoutHandler))

	http.HandleFunc("/isauth", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/json")
		if userHandler.IsAuthenticated(w, r) {
			isauth := map[string]bool{"isauth": true}
			json.NewEncoder(w).Encode(isauth)
		} else {
			isauth := map[string]bool{"isauth": false}
			json.NewEncoder(w).Encode(isauth)
		}
	}))

	http.HandleFunc("/users/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if match, _ := regexp.MatchString("/users/([0-9]+)/", r.URL.Path); match {
			switch r.Method {
			case http.MethodGet:
				userHandler.GetUser(w, r)
			}
			return
		}
		if userHandler.IsAuthenticated(w, r) {
			if match, _ := regexp.MatchString("/users/current/", r.URL.Path); match {
				switch r.Method {
				case http.MethodGet:
					userHandler.GetLoginUser(w, r)
				}
				return
			}
			// リクエストのパスが"/posts/"の時にリクエストのメソッドによって発火する関数を変える
			switch r.Method {
			case http.MethodPost:
				userHandler.SaveUser(w, r)
			case http.MethodPut:
				userHandler.EditUser(w, r)
			case http.MethodDelete:
				userHandler.DeleteUser(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}))
	// http.HandleFunc("/delay", enableCORS(userHandler.DelayTimeOut))

	http.HandleFunc("/posts/", enableCORS(func(w http.ResponseWriter, r *http.Request) {

		// リクエストのパスが"/post/数字/"の時にリクエストのメソッドによって発火する関数を変える
		if match, _ := regexp.MatchString("/posts/([0-9]+)/", r.URL.Path); match {
			switch r.Method {
			case http.MethodGet:
				postHandler.GetPost(w, r)
			case http.MethodPut:
				postHandler.EditPost(w, r)
			case http.MethodDelete:
				postHandler.DeletePost(w, r)
			}
			return
		}
		// リクエストのパスが"/posts/"の時にリクエストのメソッドによって発火する関数を変える
		switch r.Method {
		case http.MethodGet:
			postHandler.GetPosts(w, r)
		case http.MethodPost:
			postHandler.SavePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
