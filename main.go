package main

import (
	"encoding/json"
	"fmt"
)

type RequestData struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
}

type RequestParam struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Request struct {
	Template RequestData    `json:"template"`
	Users    RequestData    `json:"users"`
	Callback RequestData    `json:"callback"`
	Params   []RequestParam `json:"params"`
}

func main() {
	fmt.Println("Penlab Push Notification Worker")
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
	var r Request
	err := json.Unmarshal(jsonParam, &r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Template URI: ", r.Template.Uri)
	fmt.Println("Users URI: ", r.Users.Uri)
	fmt.Println("Callback URI: ", r.Callback.Uri)
	fmt.Println("Parameters: ", r.Params[len(r.Params)-1].Field)
}
