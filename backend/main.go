package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Link struct {
	Code      string    `json:"code"`
	Original  string    `json:"original"`
	Short     string    `json:"short"`
	CreatedAt time.Time `json:"created_at"`
}

// GET /api/links — все ссылки из БД
func listLinks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT original_url, code, created_at FROM links ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	links := []Link{}
	for rows.Next() {
		var l Link
		rows.Scan(&l.Original, &l.Code, &l.CreatedAt)
		l.Short = os.Getenv("BASE_URL") + "/s/" + l.Code
		links = append(links, l)
	}
	json.NewEncoder(w).Encode(links)
}

// POST /api/shorten {"url": "..."} — получить код у микросервиса и сохранить в БД
func shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "url required", 400)
		return
	}
	body, _ := json.Marshal(req)
	resp, err := http.Post(os.Getenv("SHORTENER_URL")+"/shorten", "application/json", bytes.NewReader(body))
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	defer resp.Body.Close()
	var res struct {
		Code string `json:"code"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	var l Link
	err = db.QueryRow("INSERT INTO links (original_url, code) VALUES ($1, $2) RETURNING created_at",
		req.URL, res.Code).Scan(&l.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	l.Code = res.Code
	l.Original = req.URL
	l.Short = os.Getenv("BASE_URL") + "/s/" + res.Code
	json.NewEncoder(w).Encode(l)
}

// DELETE /api/links/{code} — удалить ссылку
func deleteLink(w http.ResponseWriter, r *http.Request) {
	if _, err := db.Exec("DELETE FROM links WHERE code = $1", r.PathValue("code")); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /s/{code} — редирект на оригинальную ссылку
func redirect(w http.ResponseWriter, r *http.Request) {
	var url string
	err := db.QueryRow("SELECT original_url FROM links WHERE code = $1", r.PathValue("code")).Scan(&url)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("GET /api/links", listLinks)
	http.HandleFunc("POST /api/shorten", shorten)
	http.HandleFunc("DELETE /api/links/{code}", deleteLink)
	http.HandleFunc("GET /s/{code}", redirect)
	log.Println("backend on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
