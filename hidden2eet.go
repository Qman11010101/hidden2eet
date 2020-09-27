package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

// ErrResponse : Twitter error response
type ErrResponse struct {
	Code    int16
	Message string
}

// Version def
const version = 0.1

// API keys name def
const envCkey = "hidden2eet_consumer_key"
const envCsec = "hidden2eet_consumer_secret"
const envAtok = "hidden2eet_access_token"
const envAsec = "hidden2eet_access_secret"

// Arguments name def
const help = "help"
const register = "register"

func isRegistered(token string) string {
	var r string
	if token == "" {
		r = "Not registered"
	} else if len(token) > 50 {
		r = "Maybe invalid token; length is over 50"
	} else if len(token) < 20 {
		r = "Maybe invalid token; length is under 20"
	} else {
		r = "Registered"
	}
	return r
}

func registerToken(keyNum string) string {
	var n string
	var t string
	switch keyNum {
	case "1":
		n = "API key"
	case "2":
		n = "API secret"
	case "3":
		n = "access token"
	case "4":
		n = "access secret"
	}
	fmt.Println("Enter " + n)
	fmt.Scan(&t)

	return "OK"
}

func main() {

	flag.Parse()
	if flag.Arg(0) == "" || flag.Arg(0) == help {
		fmt.Println("hidden2eet v" + fmt.Sprint(version))
		fmt.Println("Usage: hidden2eet [command] [tweet-content]")
		fmt.Println("Commands:")
		fmt.Println("	register - Register account API keys")
		fmt.Println("	help     - Display this help")
	} else if flag.Arg(0) == register {
		for {
			consumerKey := os.Getenv(envCkey)
			consumerSecret := os.Getenv(envCsec)
			accessToken := os.Getenv(envAtok)
			accessSecret := os.Getenv(envAsec)

			fmt.Println("(1)API key:       " + isRegistered(consumerKey))
			fmt.Println("(2)API secret:    " + isRegistered(consumerSecret))
			fmt.Println("(3)Access token:  " + isRegistered(accessToken))
			fmt.Println("(4)Access Secret: " + isRegistered(accessSecret))
			fmt.Print("Enter the number of the token you want to register (Press 'q' to quit): ")
			var regNum string
			fmt.Scan(&regNum)

			if regNum == "q" {
				break
			}
			// switch regNum {
			// case "1":
			// 	fmt.Println("Enter API Key")
			// 	var
			// }
		}
	} else {
		consumerKey := os.Getenv(envCkey)
		consumerSecret := os.Getenv(envCsec)
		accessToken := os.Getenv(envAtok)
		accessSecret := os.Getenv(envAsec)
		if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
			fmt.Println("You haven't registered API keys!")
			fmt.Println("Execute hidden2eet with the argument 'register'.")
		} else {
			//TODO:Wrong API Key>Error
			anaconda.SetConsumerKey(consumerKey)
			anaconda.SetConsumerSecret(consumerSecret)
			api := anaconda.NewTwitterApi(accessToken, accessSecret)

			var s string
			for i := 0; i < len(flag.Args()); i++ {
				s = s + flag.Arg(i) + " "
			}
			postContent := s
			_, err := api.PostTweet(postContent, nil)
			if err == nil {
				fmt.Println("Your Tweet has been successfully sent!")
			} else {
				errStr := strings.Split(err.Error(), ", ")
				errDesc := errStr[0]
				errJSONStr := errStr[1][11 : len(errStr[1])-2]
				var errRes ErrResponse
				json.Unmarshal([]byte(errJSONStr), &errRes)
				fmt.Println("Error occured!")
				fmt.Println("Description: " + errDesc)
				fmt.Println("           : " + errRes.Message)
				fmt.Println(errStr[1])
			}
		}

	}
}
