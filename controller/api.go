package main

import (
	"fmt"
	"net/http"
)

func webHandlers() {

	err := http.ListenAndServe(apiPort, nil)
	if err != nil {
		panic(err)
	}

	//fileserver
	//fs := http.FileServer(http.Dir("src/"))
	//http.Handle("/src/", http.StripPrefix("/src/", fs))

	//handle pages
	http.HandleFunc("/auth", authWeb)
	http.HandleFunc("/fetch", fetchWeb)
	http.HandleFunc("/modify", modifyWeb)

	fmt.Println("Web - Started at " + apiPort)

}

func authWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Web - Error with request:", r)
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

}

func fetchWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Web - Error with request:", r)
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

}

func modifyWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Web - Error with request:", r)
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

}
