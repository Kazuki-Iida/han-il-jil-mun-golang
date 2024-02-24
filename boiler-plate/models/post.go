package models

import (
	"database/sql"
	"fmt"
)

// Post構造体
type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	UserId       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type PostModel struct {
	DB *sql.DB
}

func NewPostModel(DB *sql.DB) *PostModel {
	return &PostModel{DB: DB}
}

// 全件取得する関数
func (m *PostModel) All() ([]Post, error) {
	rows, err := m.DB.Query("SELECT p.id, p.title, p.body, u.id, u.name, c.id, c.name FROM posts p INNER JOIN users u ON p.user_id = u.id INNER JOIN categories c ON p.category_id = c.id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.UserId, &p.UserName, &p.CategoryId, &p.CategoryName); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// 一件取得
func (m *PostModel) GetPostById(id int) (*Post, error) {
	row, err := m.DB.Query("SELECT p.id, p.title, p.body, u.id, u.name, c.id, c.name FROM posts p INNER JOIN users u ON p.user_id = u.id INNER JOIN categories c ON p.category_id = c.id WHERE p.id = ?;", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var p Post
	for row.Next() {
		err = row.Scan(&p.ID, &p.Title, &p.Body, &p.UserId, &p.UserName, &p.CategoryId, &p.CategoryName)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		} else if err != nil {
			return nil, err
		}
	}
	if p.ID == 0 {
		return nil, fmt.Errorf("post not found")
	} else {
		return &p, nil
	}
}

// 新規作成する関数
func (m *PostModel) Insert(title string, body string) (int, error) {
	result, err := m.DB.Exec("INSERT INTO posts (title, body) VALUES (?, ?)", title, body)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Update(title string, body string, id int) error {
	_, err := m.DB.Exec("UPDATE posts SET title = ?, body = ? WHERE id = ?", title, body, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
