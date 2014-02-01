package main

import (
	"encoding/json"
	"github.com/ToQoz/go-formspec"
	"log"
	"net/http"
)

var (
	sampleFormSpec = formspec.New()
)

func init() {
	sampleFormSpec.Rule("name", formspec.RuleRequired())
	sampleFormSpec.Rule("age", formspec.RuleInt()).Message("must be integer. ok?").AllowBlank()
	sampleFormSpec.Rule("nick", formspec.RuleRequired()).FullMessage("Please enter your cool nickname.")
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) Validate() *formspec.Result {
	r := formspec.NewOkResult()

	if u.Id == 0 {
		r.Errors = append(r.Errors, formspec.NewError("id", "id is required"))
	}

	if u.Name == "" {
		r.Errors = append(r.Errors, formspec.NewError("name", "name is required"))
	}

	if len(r.Errors) > 0 {
		r.Ok = false
	}

	return r
}

func main() {
	// *** Validate form ***
	// This is basic usage.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vr := sampleFormSpec.Validate(r)

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

	// *** Validate model ***
	// Use *formspec.Result for model.Validate return value.
	// This way is good for unifying validation error expression in app.
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		u := &User{}
		vr := u.Validate()

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

		w.Write([]byte("ok"))
	})

	addr := ":8888"
	log.Print("Listen " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
