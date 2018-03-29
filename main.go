package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/evalphobia/google-home-client-go/googlehome"
)

var (
	assistantAddr = flag.String("assistant", "192.168.0.6", "IP of Google Home")
	listenPort    = flag.Int("port", 8080, "Port number")
	home          *googlehome.Client
)

const (
	bodyMaxbyteLen = 1024 * 1024
)

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodyMaxbyteLen))
	if err != nil {
		log.Printf("%s %s %s header=%+v\n", r.RemoteAddr, r.Method, r.URL, r.Header)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	txt, err := getMessageFrom(body)
	if err != nil {
		log.Printf("parse error=%s : %s\n", err.Error(), string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Speak text on Google Home.
	log.Println("say", txt)
	home.Notify(txt)
	w.WriteHeader(http.StatusOK)
}

func newLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s header=%+v\n", r.RemoteAddr, r.Method, r.URL, r.Header)
		handler.ServeHTTP(w, r)
	})
}

func getMessageFrom(js []byte) (string, error) {
	obj := map[string]interface{}{}
	err := json.Unmarshal(js, &obj)
	if err != nil {
		return "", err
	}
	e, ok := obj["events"]
	if !ok || reflect.TypeOf(e).Kind() != reflect.Slice {
		return "", fmt.Errorf("unexpected events type")
	}
	head := e.([]interface{})[0]
	if reflect.TypeOf(head).Kind() != reflect.Map {
		return "", fmt.Errorf("unexpected events type")
	}
	m, ok := head.(map[string]interface{})["message"]
	if !ok || reflect.TypeOf(m).Kind() != reflect.Map {
		return "", fmt.Errorf("unexpected events type")
	}
	txt, ok := m.(map[string]interface{})["text"]
	if !ok || reflect.TypeOf(txt).Kind() != reflect.String {
		return "", fmt.Errorf("unexpected events type")
	}
	return txt.(string), nil
}

func main() {
	var err error
	flag.Parse()
	home, err = googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: *assistantAddr,
		Lang:     "ja",
		Accent:   "JP",
	})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), newLogHandler(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}
