package url_shortener

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type App struct {
	mappings map[string]string
	hostname string
}

func NewApp(hostname string) App {
	a := App{
		mappings: make(map[string]string),
		hostname: hostname,
	}

	return a
}

func (a *App) Start() error {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{.}", a.RedirectHandler)
	rtr.HandleFunc("/", a.CreateShortURLHandler)
	rtr.HandleFunc("/create/", a.SaveHandler)

	err := http.ListenAndServe(":8080", rtr)
	return err
}

func (a *App) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../pkg/url_shortener/create_short_url.html")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	t.Execute(w, nil)
}

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")

	if val, ok := a.mappings[path[len(path)-1]]; ok {
		http.Redirect(w, r, val, http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("../pkg/url_shortener/link_doesnt_exist.html")
	if err != nil {
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
	if body[0:4] != "http" {
		body = fmt.Sprintf("http://%v", body)
	}

	hashedURL := Hashstr(body)

	//check for collision
	i := 5
	for {
		existing := a.mappings[hashedURL[:i]]
		if existing == "" {
			a.mappings[hashedURL[:i]] = body
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
