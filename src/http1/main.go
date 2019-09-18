package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"encoding/base64"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func init() {

}

type Post struct {
	User    string
	Threads []string
}

func process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	log.Println(r.Form)
}

func writeExample(w http.ResponseWriter, r *http.Request) {
	str := `<html><head><title>go web programming </title>
	<body><h1>Hello wangjr!!!</h1></body>
	</html>`

	w.WriteHeader(200)
	w.Write([]byte(str))
}

func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
	str := `No such service, try next door.`

	w.WriteHeader(501)
	w.Write([]byte(str))
}

func headerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://127.0.0.1:8080/write")
	w.WriteHeader(302)

}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tye", "application/json")

	post := &Post{
		User:    "wangjr",
		Threads: []string{"first", "second", "third"},
	}

	json, _ := json.Marshal(post)
	w.Write(json)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "wangjr",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "no work...",
		HttpOnly: true,
	}

	/*//方法1
	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
	*/

	//方法2
	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	h := r.Header["Cookie"]
	fmt.Fprintln(w, h)
}

func setMessage(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, wangjr!")
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}

	http.SetCookie(w, &c)
}
func getMessage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "no message found.")
		}
	} else {
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
		}

		http.SetCookie(w, &rc)
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	}
}

func processTmpl(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())

	var t *template.Template
	if rand.Intn(10) > 5 {
		t, _ = template.ParseFiles("layout.html", "tmp1.html")
	} else {
		t, _ = template.ParseFiles("layout.html")
	}
	t.ExecuteTemplate(w, "layout", "")
}
func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	log.WithFields(log.Fields{
		"hostname": "127.0.0.1:8080",
	}).Info("server start......")

	http.HandleFunc("/process", process)
	http.HandleFunc("/write", writeExample)
	http.HandleFunc("/writeheader", writeHeaderExample)
	http.HandleFunc("/redirect", headerExample)
	http.HandleFunc("/json", jsonExample)
	http.HandleFunc("/set-cookie", setCookie)
	http.HandleFunc("/get-cookie", getCookie)
	http.HandleFunc("/set-message", setMessage)
	http.HandleFunc("/get-message", getMessage)
	http.HandleFunc("/template", processTmpl)

	err := server.ListenAndServe()
	if err != nil {
		log.WithFields(log.Fields{
			"hostname": "127.0.0.1:8080",
		}).Info("server stop.")
	}
}
