package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CheatCoder/geddit"
)

var o *geddit.OAuthSession

func main() {
	startfyne()
	//For github
	id := flag.String("id", "", "")
	secret := flag.String("secret", "", "")
	if *id != "" && *secret != "" {
		client = *id
		clientsecret = *secret
	}
	c := make(chan string)
	fmt.Println(client, clientsecret)
	tmp, err := geddit.NewOAuthSession(
		client,                                  //Client String
		clientsecret,                            //Client Secret
		"/u/CheatGo write simple reddit client", // User-Agent
		"http://localhost:9987/reddit",          //Redirect URL - When change then Change it in login.go
	)
	o = tmp
	if err != nil {
		log.Fatal(err)
	}

	err = getLogedin(c)
	if err != nil {
		os.Exit(-99)
	}

	close(c)

	mainwinstart()
}
