package Users

import (
	"net/http"
	"github.com/gocql/gocql"
	"encoding/json"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Cassandra"
	"fmt"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUuid gocql.UUID

	user, errs := FormToUser(r)

	var created bool = false
	if len(errs) == 0 {
		fmt.Println("creating a new user")
		gocqlUuid = gocql.TimeUUID()
		if err := Cassandra.Session.Query(`
		INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUuid, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	if created {
		fmt.Println("user_id", gocqlUuid)
		json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUuid})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
