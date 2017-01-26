package Messages

import (
	"net/http"
	"github.com/gocql/gocql"
	"encoding/json"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Stream"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Cassandra"
	getstream "github.com/GetStream/stream-go"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var errStr, userIdStr, message string

	if userIdStr, errStr = processFormField(r, "user_id"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}
	user_id, err := gocql.ParseUUID(userIdStr)
	if err != nil {
		errs = append(errs, "Parameter 'user_id' not an integer")
	}

	if message, errStr = processFormField(r, "message"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}

	gocqlUuid := gocql.TimeUUID()

	var created bool = false
	if len(errs) == 0 {
		if err := Cassandra.Session.Query(`
		INSERT INTO messages (id, user_id, message) VALUES (?, ?, ?)`,
			gocqlUuid, user_id, message).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	if created {
		// send message to Stream
		globalMessages, err := Stream.Client.FlatFeed("messages", "global")
		if err == nil {
			globalMessages.AddActivity(&getstream.Activity{
				Actor: getstream.FeedID(user_id.String()),
				Verb: "post",
				Object: getstream.FeedID(gocqlUuid.String()),
				MetaData: map[string]string{
					// add as many custom keys/values here as you like
					"message": message,
				},
			})
		}

		json.NewEncoder(w).Encode(NewMessageResponse{ID: gocqlUuid})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
