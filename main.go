package main

import (
	"net/http"
	"log"
	"time"
	"os"
	"github.com/BurntSushi/toml"
	"flag"
	"math/rand"
)

const cookieName = "KEY"

var flags = flag.NewFlagSet("protected", flag.ExitOnError)
var addr = flags.String("addr", ":3000", "listen address")
var keysFile = flags.String("keys", "", "keys to load")
var serving = flags.String("serving", "", "serving directory")

var keys map[string]string

func main() {
	flags.Parse(os.Args[1:])

	fs := http.FileServer(http.Dir(*serving))
	http.Handle("/", withLogging(withAuthentication(fs.ServeHTTP)))

	keys = ReadKeys(*keysFile)

	log.Println("Listening on", *addr)
	http.ListenAndServe(*addr, nil)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Reads info from config file
func ReadKeys(configfile string) (keys map[string]string) {
	keys = make(map[string]string)
	if configfile == "" {
		key := RandStringRunes(10)
		keys[key] = "GENERATED"
		log.Println("Created random key", key)
		return
	}

	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &keys); err != nil {
		log.Fatal(err)
	}

	return keys
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ts=%s url=%s remoteAddr=%s", time.Now(), r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func checkKey(key string) bool {
	if val, ok := keys[key]; ok {
		log.Println("Authenticated with: ", val)
	}
	return true
}

func withAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		granted := false
		key := r.URL.Query().Get("key")
		if key != "" {
			if checkKey(key) {
				addKeyCookie(w, key)
				granted = true
			}
		} else {
			cookie, e := r.Cookie(cookieName)
			if e != nil && cookie != nil {
				ck := cookie.Value
				if checkKey(ck) {
					granted = true
				}
			}
		}
		if granted {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
		}
	}
}

func addKeyCookie(w http.ResponseWriter, k string) {
	c := http.Cookie{
		Name:     cookieName,
		Value:    k,
		MaxAge:   7 * 24 * 3600,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
}
