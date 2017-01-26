package Messages

import (
	"net/http"
	"github.com/gocql/gocql"
	"encoding/json"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Stream"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Cassandra"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Users"
)

func Get(w http.ResponseWriter, r *http.Request) {
	var messageList []Message
	var enrichedMessages []Message
	var userList []gocql.UUID
	var err error
	m := map[string]interface{}{}

	globalMessages, err := Stream.Client.FlatFeed("messages", "global")
	// fetch from Stream
	if err == nil {
		activities, err := globalMessages.Activities(nil)
		if err == nil {
			fmt.Println("Fetching activities from Stream")
			for _, activity := range activities.Activities {
				fmt.Println(activity)
				user_id, _ := gocql.ParseUUID(activity.Actor.Value())
				message_id, _ := gocql.ParseUUID(activity.Object.Value())
				messageList = append(messageList, Message{
					ID:      message_id,
					UserID:  user_id,
					Message: activity.MetaData["message"],
				})
				userList = append(userList, user_id)
			}
		}
	}
	// if Stream fails, pull from database instead
	if err != nil {
		fmt.Println("Fetching activities from Database")
		query := "SELECT id,user_id,message FROM messages"
		iterable := Cassandra.Session.Query(query).Iter()
		for iterable.MapScan(m) {
			user_id := m["user_id"].(gocql.UUID)
			messageList = append(messageList, Message{
				ID:      m["id"].(gocql.UUID),
				UserID:  user_id,
				Message: m["message"].(string),
			})
			userList = append(userList, user_id)
			m = map[string]interface{}{}
		}
	}

	names := Users.Enrich(userList)
	for _, message := range messageList {
		message.UserFullName = names[message.UserID.String()]
		enrichedMessages = append(enrichedMessages, message)
	}
	fmt.Println("message list after enrichment", enrichedMessages)

	json.NewEncoder(w).Encode(MessagesResponse{Messages: enrichedMessages})
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	var message Message
	var errs []string
	var found bool = false

	vars := mux.Vars(r)
	id := vars["message_uuid"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,user_id,message FROM messages WHERE id=? LIMIT 1"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			user_id := m["user_id"].(gocql.UUID)
			names := Users.Enrich([]gocql.UUID{user_id})
			fmt.Println("names", names)
			message = Message{
				ID: user_id,
				UserID: m["user_id"].(gocql.UUID),
				UserFullName: names[user_id.String()],
				Message: m["message"].(string),
			}
		}
		if !found {
			errs = append(errs, "Message not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetMessageResponse{Message: message})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
