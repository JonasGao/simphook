package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func write(w http.ResponseWriter, content string, title string) {
	_, err := w.Write([]byte(content))
	if err != nil {
		log.Println("Write response: ", title, " error: ", err)
	}
}

func run(script string, w http.ResponseWriter) {
	cmd := exec.Command("./" + script)
	cmd.Stdout = w
	w.WriteHeader(200)
	err := cmd.Run()
	write(w, "\n\n--------------------\n", "[split line]")
	if err != nil {
		log.Println("Run cmd return error: ", err)
		write(w, "Run cmd fail.", "[cmd result]")
	} else {
		log.Println("Call \"" + script + "\" success")
		write(w, "Run cmd success.", "[cmd result]")
	}
}

func handle(w http.ResponseWriter, req *http.Request) {
	scriptName := strings.TrimPrefix(req.URL.Path, "/") + ".sh"
	_, err := os.Stat(scriptName)
	if os.IsNotExist(err) {
		w.WriteHeader(404)
		log.Println("Not found: " + scriptName)
		return
	}
	run(scriptName, w)
}

func main() {
	// Get and print cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Can not get current work directory.")
	}
	log.Println("Search script file in \"" + cwd + "\"")
	// Setup handle
	http.HandleFunc("/", handle)
	// Get listen for http port
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		// Setup default port
		port = "8090"
	}
	// Get listen for http host
	host, set := os.LookupEnv("HTTP_HOST")
	if !set {
		host = "127.0.0.1"
	}
	addr := host + ":" + port
	log.Println("Listen for ", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
