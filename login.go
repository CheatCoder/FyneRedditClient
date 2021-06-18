package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
)

func getLogedin(c chan string) error {
	srv := getCode(c)
	// Define scopes needed
	scopes := []string{"read", "identity", "history", "mysubreddits", "save", "subscribe"}
	// Get authorization URL for user to visit
	// TODO: generate random string
	rurl := o.AuthCodeURL("random string", scopes)

	loginurl, err := url.Parse(rurl)
	if err != nil {
		return err
	}

	a := fyne.CurrentApp()
	a.OpenURL(loginurl)

	code := <-c
	log.Println(code)
	// Create and set token using the authorization code
	err = o.CodeAuth(code)
	if err != nil {
		return err
	}
	c <- "ok"
	srv.Shutdown(context.TODO())

	return nil
}

func getCode(c chan string) *http.Server {
	srv := &http.Server{Addr: ":9987"}

	http.HandleFunc("/reddit", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		code := r.FormValue("code")
		//TODO: Check if random string is the same eg. line 18/19
		c <- code
		fmt.Fprintln(rw, "Try to Login ....")
		<-c
		fmt.Fprintln(rw, "Sucess, you can close this Tab now")

	})

	go func() {
		srv.ListenAndServe()
	}()

	return srv

}
