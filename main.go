package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
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

type User struct {
	UID  string `json:"UID"`
	Name string `json:"name"`
}

type UserResponse struct {
	ItemsReceived int    `json:"itemsReceived"`
	CurrentPage   int    `json:"curPage"`
	NextPage      int    `json:"nextPage"`
	PageTotal     int    `json:"pageTotal"`
	Users         []User `json:"items"`
}

func main() {
	start := time.Now()
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

	/*** 				Load User UID Start 					**/

	var page_number, total_pages int
	page_number = 1
	total_pages = 1

	var page RequestParam
	var UIDs [][]string

	for next := true; next; next = page_number <= total_pages {
		fmt.Println("Page: ", page_number)
		var parameters []RequestParam
		page.Field = "page"
		str_page := strconv.Itoa(page_number)
		page.Value = str_page
		parameters = append(parameters, page)
		var response = httpRequestUser("https://x5i7-qk19-xc7o.s2.xano.io/api:vjzuYRWj:v1/test/users", parameters)
		page_number++
		total_pages = response.PageTotal
		fmt.Println("Current Page: ", response.CurrentPage)
		fmt.Println("Items Received: ", response.ItemsReceived)

		count := 0
		var _uid []string
		for _, u := range response.Users {
			count++
			_uid = append(_uid, u.UID)

			if count%10 == 0 {
				UIDs = append(UIDs, _uid)
				_uid = nil
			}
		}
	}

	/*** 				Load User UID End 					**/

	// var tokens [][]string

	for _, uidGroup := range UIDs {
		// Querying FCM Tokens
		fmt.Println("Get FCM Tokens:")
		iter := client.Collection("users").Where("uid", "in", uidGroup).Documents(ctx)

		var token []string

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
				token = append(token, fcm_token.Data()["fcm_token"].(string))
			}
		}

	}
	// // Querying FCM Tokens
	// fmt.Println("Get FCM Tokens:")
	// // iter := client.Collection("users").Where("uid", "in", UIDs).Documents(ctx)
	// iter := client.Collection("users").Where("uid", "in", []string{"0K0DRcLOb1QYJAleec982yKN4pj2", "klAWYeeXBKMm7EKwS9gfXWAIrRL2"}).Documents(ctx)

	// var token []string

	// //Iterating FCM Tokens
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	fcm_tokens := doc.Ref.Collection("fcm_tokens").Documents(ctx)
	// 	for {
	// 		fcm_token, fcm_err := fcm_tokens.Next()
	// 		if fcm_err == iterator.Done {
	// 			break
	// 		}
	// 		if fcm_err != nil {
	// 			log.Fatalln(fcm_err)
	// 		}
	// 		token = append(token, fcm_token.Data()["fcm_token"].(string))
	// 	}
	// }

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// client.Close()

	// //Initializing FCM clients
	// fcm_client, err := app.Messaging(ctx)
	// if err != nil {
	// 	log.Fatalf("error getting Messaging client: %v\n", err)
	// }
	// // Send Messages
	// message := &messaging.MulticastMessage{
	// 	Notification: &messaging.Notification{
	// 		Title: "Sample",
	// 		Body:  "Content of the sample",
	// 	},
	// 	Tokens: tokens,
	// }

	// br, err := fcm_client.SendMulticast(context.Background(), message)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)

	elapsed := time.Since(start)
	log.Printf("Loading took %s", elapsed)
}

/**
Function: Request User Data
Description: Gets User Display Name and UID from the database
Parameters:
	uri				URL Endpoint where to get the list of user data
	params			parameter sent to the Endpoint
Return:
	response		UserResponse

**/

func httpRequestUser(uri string, params []RequestParam) UserResponse {
	parameters := url.Values{}

	for i, p := range params {
		if i >= 0 {
			parameters.Add(p.Field, p.Value)
			fmt.Println("Field: ", p.Field)
			fmt.Println("Value: ", p.Value)
		}
	}

	http_request, err := http.PostForm(uri, parameters)
	if err != nil {
		log.Fatalln(err)
	}

	resBody, err := ioutil.ReadAll(http_request.Body)

	var response UserResponse
	json.Unmarshal(resBody, &response)

	return response

}
