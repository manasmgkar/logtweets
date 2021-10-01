package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

//Struct to parse webhook load
type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

//Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

//Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

func CreateClient() *http.Client {
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

func registerWebhook(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Registering webhook...")
	httpClient := CreateClient()

	//Set parameters
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/webhooks.json"
	values := url.Values{}
	values.Set("url", os.Getenv("APP_URL")+"/webhook/twitter")

	//Make Oauth Post with parameters
	resp, _ := httpClient.PostForm(path, values)
	defer resp.Body.Close()
	//Parse response and check response
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp)
	//mt.Println(body)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		fmt.Fprintf(writer, "Webhook id of "+data["id"].(string)+" has been registered")
		log.Println("Webhook id of " + data["id"].(string) + " has been registered")
	} else if resp.StatusCode != 200 {
		log.Println("Could not register webhook response below:")
		log.Println(string(body))
		fmt.Fprintf(writer, string(body)) // convert body from json to string
	}
	subscribeWebhook()
}

func subscribeWebhook() {
	fmt.Println("Subscribing webapp......")
	client := CreateClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
	resp, err := client.PostForm(path, nil)
	if err != nil {
		log.Println(err)
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	//If response code is 204 it was succesful
	if resp.StatusCode == 204 {
		log.Println("Subscribed succesfully")
	} else if resp.StatusCode != 204 {
		log.Println("Could not subscribe to webhook response below:") // convert body from json to string
		log.Println(resp.StatusCode)
	}
}

func CrcCheck(writer http.ResponseWriter, req *http.Request) {
	// setresponse header to json type
	writer.Header().Set("Content-Type", "application/json")
	// get crc token in parameter
	token := req.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(writer, "No crc_token given")
		log.Println("No crc_token given")
		return
	}

	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//Geberate resoibnse strubg nao
	response := make(map[string]string)
	response["response_token"] = "sha256=" + encoded
	//turn response map to json and sent it to writer
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(writer, string(responseJson))
	log.Println(string(responseJson))
}

func twitterwebhookdump(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Handler called")
	webhookData := make(map[string]interface{})
	err := json.NewDecoder(req.Body).Decode(&webhookData)
	if err != nil {
		fmt.Fprintln(writer, "An error occured in decoding the webhookload: "+err.Error())
		log.Println("An error occured in decoding the webhookload: " + err.Error())

	}
	fmt.Println("got webhook payload: ")
	jsonizedata, err := json.Marshal(webhookData)
	if err != nil {
		log.Panic(err)
		return
	}
	if err != nil {
		log.Println("An error occured: " + err.Error())
	}
	f, err := os.OpenFile("./data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//log.Println(err)
		log.Panic(err)
		return
	}
	_, err = fmt.Fprintln(f, string(jsonizedata))
	if err != nil {
		//log.Println(err)
		log.Panic(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("Latest webhook data dumped sucessfully")
}
