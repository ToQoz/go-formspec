package main

import (
	"encoding/json"
	"fmt"
	"github.com/ToQoz/go-formspec"
	"log"
	"net/http"
)

type ErrorJson struct {
	Errors []*formspec.Error `json:"errors"`
}

func main() {
	s := formspec.New()

	s.Rule("name", formspec.RuleRequired())
	s.Rule("age", formspec.RuleInt()).Message("must be integer. ok?").AllowBlank()
	s.Rule("nick", formspec.RuleRequired()).FullMessage("Please enter your cool nickname.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vr := s.Validate(r)

		if !vr.Ok {
			w.WriteHeader(403)

			for _, verr := range vr.Errors {
				w.Write([]byte(fmt.Sprintf("Validation error in %s. %s\n", verr.Field, verr.Message)))
			}

			return
		}

		w.Write([]byte("ok"))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		vr := s.Validate(r)

		if !vr.Ok {
			j, err := json.Marshal(vr)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf8")
			w.WriteHeader(403)
			w.Write(j)

			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write([]byte(`{"message": "ok"}`))
	})

	addr := ":8888"
	log.Print("Listen " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
