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
		ok, errs := s.Validate(r)

		if !ok {
			w.WriteHeader(403)

			for _, err := range errs {
				verr := err.(*formspec.Error)
				w.Write([]byte(fmt.Sprintf("Validation error in %s. %s\n", verr.Field, verr.Message)))
			}

			return
		}

		w.Write([]byte("ok"))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		ok, errs := s.Validate(r)

		if !ok {
			var verrs []*formspec.Error

			for _, err := range errs {
				verrs = append(verrs, err.(*formspec.Error))
			}

			j, err := json.Marshal(&ErrorJson{Errors: verrs})

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
