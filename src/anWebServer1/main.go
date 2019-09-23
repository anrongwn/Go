package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func init() {

}

type X string

// String : X String
func (x X) String() string {
	return fmt.Sprintf("<%s>", string(x))
}

func change(s ...string) {
	s[0] = "Golang"
	s = append(s, "playground")
	fmt.Println(s)
}

func main() {
	welcome := []string{"hello", "wangjr"}
	change(welcome...)
	fmt.Println(welcome)

	a := [...]float64{67.7, 89.8, 21, 78}
	sum := float64(0)
	for i, v := range a { //range returns both the index and value
		fmt.Printf("%d the element of a is %.2f\n", i, v)
		sum += v
	}
	fmt.Println("\nsum of all elements of a", sum)

	var b []float64
	b = a[0:len(a)]
	for ii, vv := range b {
		fmt.Printf("%d the element of a is %.2f\n", ii, vv)
		vv++
		b[ii]++
	}

	for i, v := range a { //range returns both the index and value
		fmt.Printf("%d the element of a is %.2f\n", i, v)
		sum += v
	}

	salse := make(map[string]int)
	salse["wangjr"] = 100
	salse["wangjia"] = 90
	salse["wangyuayua"] = 100
	fmt.Println(salse)
	delete(salse, "wangjr")
	for k, v := range salse {
		fmt.Printf("key=%s, value=%d\n", k, v)
	}

	return

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	log.Println("anWebServer1 start...")
	http.HandleFunc("/post/", handleRequest)
	server.ListenAndServe()

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t\t")
	if err != nil {
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(output)
	return
}

func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var post Post
	json.Unmarshal(body, &post)
	err = post.create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return

}
func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	err = post.delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return

}
