package main

import (
	"fmt"
	"net/http"
)

func main()  {
	http.HandleFunc("/",func(w http.ResponseWriter, r*http.Request){

		fmt.Fprintf(w,"<h1>hello %s </h1>",r.FormValue("name"))
	})

	http.ListenAndServe(":8888",nil)
}
