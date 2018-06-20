package tag_api

import (
	"encoding/json"

	"github.com/nats-io/go-nats"
)

func (data *ApiData) ConnectNATS() (err error) {
	Log.Info.Printf("Connecting to %s\n", data.NHost)
	data.NConn, err = nats.Connect(data.NHost)
	return
}

func (data *ApiData) ListenNATSSub() {
	var qMsg QueueMessage
	ch := make(chan *nats.Msg, 64)

	Log.Info.Printf("Subscribing to nats channel %q\n", NSub)
	sub, err := data.NConn.ChanSubscribe(NSub, ch)
	if err != nil {
		Log.Error.Println(err)
		return
	}
	defer sub.Unsubscribe()

	for {
		msg := <-ch
		err = json.Unmarshal(msg.Data, &qMsg)
		if err != nil {
			Log.Error.Println(err)
		}

		switch qMsg.Command {
		case "adduser":
			err = data.AddUser(msg.Data)
			if err != nil {
				Log.Error.Printf("Add User: %v\n", err)
			}
		default:
			Log.Info.Printf("Unrecognized command: %s\n", qMsg.Command)
		}
	}
}

func (data *ApiData) MessageAddUser(u User) (err error) {
	var b []byte

	uMsg := UserMessage{
		Command:   "adduser",
		Id:        u.Id,
		GroupId:   u.GroupId,
		Guid:      u.Guid,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
	b, err = json.Marshal(uMsg)
	if err != nil {
		return
	}

	// Send message to content server
	err = d.NConn.Publish(NSub, b)
	Log.Info.Printf("Authenticate: %s %s [id=%d]\n", u.FirstName, u.LastName, u.Id)
	return
}
