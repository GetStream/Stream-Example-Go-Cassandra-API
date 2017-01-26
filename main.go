package main

/*
ccm create -v 2.1.1 streamdemoapi
ccm populate -n 1
echo "CREATE KEYSPACE demoapi WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};" | cqlsh

echo "use streamdemoapi; drop table messages; create table messages (id UUID, user_id UUID, message text, PRIMARY KEY(id));" | cqlsh
echo "use streamdemoapi; create index on messages (user_id);" | cqlsh
echo "use streamdemoapi; drop table users; CREATE TABLE users ( id UUID, firstname text, lastname text, age int, email text, city text, PRIMARY KEY (id));" | cqlsh
*/


import (
	"net/http"
	"github.com/gocql/gocql"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Cassandra"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Stream"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Users"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"github.com/GetStream/Stream-Example-Go-Cassandra-API/Messages"
)

type HeartbeatResponse struct {
	Status string `json:"status"`
	Code   int `json:"code"`
}

type Message struct {
	ID      gocql.UUID `json:"id"`
	UserID  gocql.UUID `json:"user_id"`
	Message string `json:"message"`
}

func main() {
	err := Stream.Connect(
		"ax3bm9tjcb35",
		"ec9tydddc78zmmc8r682j43z5y6exkaa98myamzk23pztm5yb8s4c66g5737eyey",
		"us-east")
	if err != nil {
		log.Fatal("Could not connect to Stream, abort")
	}

	CassandraSession := Cassandra.Session
	defer CassandraSession.Close()

	//var id gocql.UUID
	//var text string
	//if err := session.Query(`SELECT id, message FROM messages WHERE timeline = ? LIMIT 1`,
	//	"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Message:", id, text)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Heartbeat)

	router.HandleFunc("/users", Users.Get)
	router.HandleFunc("/users/new", Users.Post)
	router.HandleFunc("/users/{user_uuid}", Users.GetOne)

	router.HandleFunc("/messages", Messages.Get)
	router.HandleFunc("/messages/new", Messages.Post)
	router.HandleFunc("/messages/{message_uuid}", Messages.GetOne)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(HeartbeatResponse{Status: "OK", Code: 200})
}

