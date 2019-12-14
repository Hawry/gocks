package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hawry/gomocka/settings"
)

var availableHandlers []http.HandlerFunc

func createHandler(code int, body string, headers map[string]string, m settings.Mock, s settings.Settings) http.HandlerFunc {
	log.Printf("debug: mock data %+v", m)
	return func(w http.ResponseWriter, r *http.Request) {
		if s.RequireAuthentication() && !m.DisableAuth {
			u, p, _ := r.BasicAuth()
			if !s.VerifyBasicAuth(u, p) && !s.VerifyHeaderAuth(r.Header) && !s.VerifyOpenIDAuth(r.Header) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		b, wc := m.Wildcard()
		if b {
			vars := mux.Vars(r)
			for _, s := range wc {
				if val, ok := vars[s]; !ok {
					continue
				} else {
					dc := fmt.Sprintf("{%s}", s)
					body = strings.ReplaceAll(body, dc, val)
				}
			}
		}

		log.Printf("debug: %s %s", r.Method, r.RequestURI)
		for k, v := range headers {
			log.Printf("debug: adding %v, %v", k, v)
			w.Header().Add(k, v)
		}
		w.WriteHeader(code)
		w.Write([]byte(body))
	}
}
