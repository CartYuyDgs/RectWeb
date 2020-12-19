package main

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("token")
	msg := r.PostFormValue("msg")
	base64 := r.PostFormValue("photobase64")

	fmt.Fprintln(w, "hello world")
	fmt.Fprintln(w, token)
	fmt.Fprintln(w, msg)
	fmt.Fprintln(w, base64)
}

func main() {
	http.HandleFunc("/",    IndexHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)
}



