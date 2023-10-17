package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func callShell(script string) (string, string) {
	cmd := exec.Command("./" + script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	outString := out.String()
	if err != nil {
		log.Println("Run cmd return error: ", err)
		return outString, err.Error()
	}
	log.Println("Call \"" + script + "\" success")
	return outString, ""
}

func write(w http.ResponseWriter, content string) {
	_, err := w.Write([]byte(content))
	if err != nil {
		log.Println("Write response return error: ", err)
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
	successOutput, errorOutput := callShell(scriptName)
	if errorOutput != "" {
		w.WriteHeader(500)
		write(w, successOutput+"\n\n"+errorOutput+"\n")
		return
	}
	w.WriteHeader(200)
	write(w, successOutput)
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
