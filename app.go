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
	if err != nil {
		log.Println("has error: ", err)
		return "", err.Error()
	}
	outString := out.String()
	log.Println("call \"" + script + "\" success")
	return outString, ""
}

func write(w http.ResponseWriter, code int, content string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(content))
	if err != nil {
		log.Println("has error: ", err)
	}
}

func handle(w http.ResponseWriter, req *http.Request) {
	scriptName := strings.TrimPrefix(req.URL.Path, "/") + ".sh"
	_, err := os.Stat(scriptName)
	if os.IsNotExist(err) {
		w.WriteHeader(404)
		return
	}
	successOutput, errorOutput := callShell(scriptName)
	if errorOutput != "" {
		write(w, 500, errorOutput)
		return
	}
	write(w, 200, successOutput)
}

func main() {
	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
