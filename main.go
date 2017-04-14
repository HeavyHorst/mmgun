package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/42wim/matterbridge/matterhook"
)

func mailHandler(mclient *matterhook.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println("[ERROR] couldn't parse form: " + err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		subject := r.Form.Get("subject")
		text := r.Form.Get("stripped-text")

		msg := matterhook.OMessage{
			UserName: r.Form.Get("from"),
			Text:     fmt.Sprintf("## %s\n%s", subject, text),
		}

		err = mclient.Send(msg)
		if err != nil {
			log.Println("[ERROR] couldn't send message: " + err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("[INFO] forwarded email from %s\n", msg.UserName)
	}
}

func main() {

	var hostaddr = flag.String("addr", "0.0.0.0", "specify the address to listen on.")
	var port = flag.String("port", "8003", "specify the port to listen on.")
	var url = flag.String("url", "", "the Mattermost Webhook url")
	var help = flag.Bool("help", false, "output usage informations")

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	mconfig := matterhook.Config{
		DisableServer: true,
	}

	client := matterhook.New(*url, mconfig)

	mux := http.NewServeMux()
	mux.HandleFunc("/mail", mailHandler(client))
	http.ListenAndServe(fmt.Sprintf("%s:%s", *hostaddr, *port), mux)
}
