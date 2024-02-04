package models

import (
	"database/sql"
	"fmt"
)

// User構造体
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(DB *sql.DB) *UserModel {
	return &UserModel{DB: DB}
}

// // 全件取得する関数
// func (m *PostModel) All() ([]Post, error) {
// 	rows, err := m.DB.Query("SELECT id, title, body FROM posts")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []Post
// 	for rows.Next() {
// 		var post Post
// 		if err := rows.Scan(&post.ID, &post.Title, &post.Body); err != nil {
// 			return nil, err
// 		}
// 		posts = append(posts, post)
// 	}

//		return posts, nil
//	}
//
// idを使って一件取得
func (m UserModel) GetUserById(id int) (*User, error) {
	row, err := m.DB.Query("SELECT id, name, email, password FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var user User
	for row.Next() {
		err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		} else if err != nil {
			return nil, err
		}
	}
	// なぜかレコードが見つからなかった時にuser.IDに0がセットされるから一旦これでエラーハンドリング
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	} else {
		return &user, nil
	}
}

// メアドを使って一件取得
func (m UserModel) GetUserByEmail(email string) (*User, error) {
	row, err := m.DB.Query("SELECT id, name, email, password FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var user User
	for row.Next() {
		err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		} else if err != nil {
			return nil, err
		}
	}
	// なぜかレコードが見つからなかった時にuser.IDに0がセットされるから一旦これでエラーハンドリング
	if user.ID != 0 {
		return &user, nil
	} else {
		return nil, fmt.Errorf("user not found")
	}
}

// 新規作成する関数
func (m *UserModel) Insert(name string, email string, password string) (int, error) {
	result, err := m.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
	fmt.Println(err)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Update(name string, email string, password string, id int) error {
	_, err := m.DB.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", name, email, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
