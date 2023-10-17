package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func callShell(script string, w http.ResponseWriter) {
	cmd := exec.Command("./" + script)
	cmd.Stdout = w
	w.WriteHeader(200)
	err := cmd.Run()
	if err != nil {
		log.Println("Run cmd return error: ", err)
	} else {
		log.Println("Call \"" + script + "\" success")
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
	callShell(scriptName, w)
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Can not get current work directory.")
	}
	log.Println("Search script file in \"" + cwd + "\"")
	http.HandleFunc("/", handle)
	log.Println("Listen for :8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
