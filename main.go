package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Structure for the request data
type RequestData struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
}

// structure for the request parameters
type RequestParam struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// structure for the whole JSON Data
type Request struct {
	Template RequestData    `json:"template"`
	Users    RequestData    `json:"users"`
	Callback RequestData    `json:"callback"`
	Params   []RequestParam `json:"params"`
}

func main() {
	fmt.Println("Penlab Push Notification Worker")

	// Sample JSON data to be parsed
	jsonParam := []byte(`{
		"template": {
			"uri": "https://x5i7-qk19-xc7o.s2.xano.io/api:push-notification/template/like",
			"method":"POST"
		},
		"users": {
			"uri": "https://x5i7-qk19-xc7o.s2.xano.io/api:push-notification/users/like",
			"method":"POST"
		},
		"callback": {
			"uri": "https://x5i7-qk19-xc7o.s2.xano.io/api:push-notification/callback/like",
			"method":"POST"
		},
		"params": [
			{
				"field":"from",
				"value":"110325"
			},
			{
				"field":"to",
				"value":"110326"
			}
		]
	}`)
	// Deserializng JSON data to Object
	var r Request
	err := json.Unmarshal(jsonParam, &r)
	if err != nil {
		panic(err)
	}

	// Configuring Firebase Admin SDK
	opt := option.WithCredentialsFile("creds/penlab-duplicate-firebase-adminsdk-f1owv-19116da36d.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	// Connecting to FireStore
	client, err := app.Firestore(context.Background())
	ctx := context.Background()

	// Querying FCM Tokens
	fmt.Println("Get FCM Tokens:")
	iter := client.Collection("users").Where("uid", "in", []string{"0K0DRcLOb1QYJAleec982yKN4pj2"}).Documents(ctx)

	var tokens []string

	//Iterating FCM Tokens
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fcm_tokens := doc.Ref.Collection("fcm_tokens").Documents(ctx)
		for {
			fcm_token, fcm_err := fcm_tokens.Next()
			if fcm_err == iterator.Done {
				break
			}
			if fcm_err != nil {
				log.Fatalln(fcm_err)
			}
			tokens = append(tokens, fcm_token.Data()["fcm_token"].(string))
		}
	}

	if err != nil {
		log.Fatalln(err)
	}
	client.Close()

	//Initializing FCM clients
	fcm_client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	// Send Messages
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Sample",
			Body:  "Content of the sample",
		},
		Tokens: tokens,
	}

	br, err := fcm_client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
}
