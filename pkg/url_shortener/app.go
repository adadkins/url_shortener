package url_shortener

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	hostname string
	db       *sql.DB
}

func NewApp(hostname string, db *sql.DB) App {
	a := App{
		hostname: hostname,
		db:       db,
	}
	return a
}

func (a *App) Start() error {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", a.CreateShortURLHandler).Methods("GET")
	rtr.HandleFunc("/{^.*[a-zA-Z0-9]+.*$}", a.RedirectHandler).Methods("GET")
	rtr.HandleFunc("/", a.SaveHandler).Methods("POST")

	err := http.ListenAndServe(":8080", rtr)
	return err
}

func (a *App) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../pkg/url_shortener/create_short_url.html")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t.Execute(w, nil)
}

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	link := a.queryHash(path[len(path)-1])
	if link != "" {
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("../pkg/url_shortener/link_doesnt_exist.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Resource Not Found"))
	}
	t.Execute(w, nil)
}

func (a *App) SaveHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("url")
	if body == "" || len(body) < 5 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//need some logic to test if has http/https and if not, add it to the stored url
	//TODO theres a bug here.
	if body[0:4] != "http" {
		body = fmt.Sprintf("http://%v", body)
	}
	hashedURL := Hashstr(body)
	//check for collision
	//TODO dont create a new hash if the link is the same
	i := 5
	for {
		existing := a.queryHash(hashedURL[:i])
		if existing == "" {
			err := a.insertHash(hashedURL[:i], body)
			if err != nil {
				log.Println(err)
			}
			break
		}
		i++
	}

	t, _ := template.ParseFiles("../pkg/url_shortener/generated_short_url.html")
	t.Execute(w, fmt.Sprintf("http://%v/%v", a.hostname, hashedURL[:i]))
}

func Hashstr(Txt string) string {
	h := sha1.New()
	h.Write([]byte(Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x\n", bs))
	return sh
}

func (a *App) queryHash(hash string) string {
	var link string
	queryStatement := "SELECT link from public.shorturl where hashcode = $1"
	if err := a.db.QueryRow(queryStatement, hash).Scan(&link); err != nil {
		if err != nil || err == sql.ErrNoRows {
			log.Println(err)
			return ""
		}
	}
	return link
}

func (a *App) insertHash(hash, link string) error {
	queryStatement := "INSERT into public.shorturl (hashcode, link) values($1, $2);"
	if _, err := a.db.Exec(queryStatement, hash, link); err != nil {
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
