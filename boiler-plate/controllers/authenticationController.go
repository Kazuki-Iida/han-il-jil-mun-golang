package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

type sessionInfo struct {
	isAllowed bool
	userId    int
}

var (
	// セッション情報を保存するためのマップ
	sessions = make(map[string]*sessionInfo)
	// セッション情報へのアクセスを同期するためのミューテックス
	sessionMutex = &sync.Mutex{}
)

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h UserController) GetLoginUser(w http.ResponseWriter, r *http.Request) {
	session_id := r.Header.Get("session_id")
	sessionMutex.Lock()
	userId := sessions[session_id].userId
	for id := range sessions {
		fmt.Println("loginHandler: Current session ID is ", id)
	}
	sessionMutex.Unlock()
	user, err := h.Model.GetUserById(userId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(*user)
}

func (h UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 本来はIDとパスワードでのユーザー認証のロジックを実装する
	// ->してみた。
	var loginInfo LoginInfo
	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	user, err := h.Model.GetUserByEmail(loginInfo.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if loginInfo.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid password")
		return
	} else {
		seed := time.Now().UnixNano()
		rand.Seed(seed)
		sessionID := strconv.FormatUint(rand.Uint64(), 16)

		// サーバーサイドでセッションを保存
		sessionMutex.Lock()
		sessions[sessionID] = &sessionInfo{isAllowed: true, userId: user.ID}
		for id := range sessions {
			fmt.Println("loginHandler: Current session ID is ", id)
		}
		sessionMutex.Unlock()

		data := map[string]string{"session_id": sessionID}
		json.NewEncoder(w).Encode(data)
	}
}

func (h UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 本来はIDとパスワードでのユーザー認証のロジックを実装する
	// ->してみた。
	var registerInfo RegisterInfo
	err := json.NewDecoder(r.Body).Decode(&registerInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}

	id, err := h.Model.Insert(registerInfo.Name, registerInfo.Email, registerInfo.Password)
	if err != nil {
		dupField := parseDuplicateField(err.Error())
		switch err.(*mysql.MySQLError).Number {
		case 1062:
			// 重複キーエラー
			if dupField == "email" {
				http.Error(w, "duplicate email", http.StatusBadRequest)
				// } else if dupField == "name" {
				// 	http.Error(w, "duplicate name", http.StatusBadRequest)
			}
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	seed := time.Now().UnixNano()
	rand.Seed(seed)
	sessionID := strconv.FormatUint(rand.Uint64(), 16)

	// サーバーサイドでセッションを保存
	sessionMutex.Lock()
	sessions[sessionID] = &sessionInfo{isAllowed: true, userId: id}
	for id := range sessions {
		fmt.Println("loginHandler: Current session ID is ", id)
	}
	sessionMutex.Unlock()

	data := map[string]string{"session_id": sessionID}
	json.NewEncoder(w).Encode(data)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからセッションIDを取得
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionMutex.Lock()
	sessionInfo, ok := sessions[cookie.Value]
	if ok {
		for id := range sessions {
			fmt.Println("dashboardHandler: Current session ID is ", id)
		}
	}
	sessionMutex.Unlock()

	// セッションが有効かチェック
	if !ok || !sessionInfo.isAllowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Write([]byte("Welcome to your dashboard!"))
}

func (h UserController) IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	// CookieからセッションIDを取得
	session_id := r.Header.Get("session_id")
	if session_id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	sessionMutex.Lock()
	sessionInfo, ok := sessions[session_id]
	if ok {
		for id := range sessions {
			fmt.Println("IsAuthenticated: Current session ID is ", id)
		}
	}
	sessionMutex.Unlock()
	// セッションが有効かチェック
	if !ok || !sessionInfo.isAllowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}

	return true
}

func (h UserController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// CookieからセッションIDを取得し、サーバーサイドでセッションを削除
	cookie, err := r.Cookie("session_id")
	if err == nil {
		sessionMutex.Lock()
		delete(sessions, cookie.Value)
		if len(sessions) == 0 {
			fmt.Println("logoutHandler: Current session ID is nothing")
		}
		sessionMutex.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Cookieを削除
	})
	fmt.Println(sessions)
	w.Write([]byte("Logout successful!"))
}

func (h UserController) DelayTimeOut(w http.ResponseWriter, r *http.Request) {
	// CookieからセッションIDを取得し、サーバーサイドでセッションを延長

	cookie, err := r.Cookie("session_id")
	userId := sessions[cookie.Value].userId
	if err == nil {
		sessionMutex.Lock()
		// 一旦sessionIDを削除
		delete(sessions, cookie.Value)
		if len(sessions) == 0 {
			fmt.Println("logoutHandler: Current session ID is nothing")
		}
		sessionMutex.Unlock()
	}
	// 再作成
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	sessionID := strconv.FormatInt(rand.Int63(), 10)
	// サーバーサイドでセッションを保存
	sessionMutex.Lock()
	sessions[sessionID] = &sessionInfo{isAllowed: true, userId: userId}
	for id := range sessions {
		fmt.Println("loginHandler: Current session ID is ", id)
	}
	sessionMutex.Unlock()

	// クライアントにセッションIDをCookieとして送信
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  sessionID,
		Path:   "/",
		MaxAge: 60, // 有効期限を60秒（1分）に設定
	})

	w.Write([]byte("Delay Time Out successful!"))
}

func parseDuplicateField(msg string) string {

	// メッセージから重複キーを抽出する正規表現
	regex := `Duplicate entry '(?P<entry>.+)' for key 'users\.(?P<column>\w+)'`

	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(msg)

	return match[2]
}
