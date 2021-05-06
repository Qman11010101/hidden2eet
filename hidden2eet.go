package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

// ErrResponse : Twitter error response
type ErrResponse struct {
	Code    int16
	Message string
}

// SettingHidden2eet : Format of setting file
type SettingHidden2eet struct {
	Consumerkey    string `json:"consumerKey"`
	Consumersecret string `json:"consumerSecret"`
	Accesstoken    string `json:"accessToken"`
	Accesssecret   string `json:"accessSecret"`
}

// Version def
const version = 0.1

// Blank setting file
const settingBase = `{
	"consumerKey": "",
	"consumerSecret": "",
	"accessToken": "",
	"accessSecret": ""
}
`

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

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

var homeDirectory = os.Getenv("HOMEPATH")
var hidden2eetFile = "C:" + homeDirectory + "\\.hidden2eet"

func main() {

	flag.Parse()
	if flag.Arg(0) == "" || flag.Arg(0) == help {
		fmt.Println("hidden2eet v" + fmt.Sprint(version))
		fmt.Println("Usage: hidden2eet [command] [tweet-content]")
		fmt.Println("Commands:")
		fmt.Println("	register - Register account API keys")
		fmt.Println("	help     - Display this help")
	} else if flag.Arg(0) == register {
		if !exists(hidden2eetFile) {
			err := ioutil.WriteFile(hidden2eetFile, []byte(settingBase), 0666)
			if err != nil {
				panic(err)
			}
		}

		setting, err := ioutil.ReadFile(hidden2eetFile)
		if err != nil {
			panic(err)
		}

		var settingJson SettingHidden2eet
		err = json.Unmarshal(setting, &settingJson)
		if err != nil {
			panic(err)
		}

		for {
			consumerKey := settingJson.Consumerkey
			consumerSecret := settingJson.Consumersecret
			accessToken := settingJson.Accesstoken
			accessSecret := settingJson.Accesssecret

			fmt.Println("(1)API key:       " + isRegistered(consumerKey))
			fmt.Println("(2)API secret:    " + isRegistered(consumerSecret))
			fmt.Println("(3)Access token:  " + isRegistered(accessToken))
			fmt.Println("(4)Access secret: " + isRegistered(accessSecret))
			fmt.Print("Enter the number of the token you want to register or change (Press 'q' to quit): ")
			var regNum string
			fmt.Scan(&regNum)

			if regNum == "q" {
				if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
					fmt.Println("NOTICE: Unregistered API keys exist!")
				}
				break
			} else if regNum == "1" {
				fmt.Print("Input your API Key: ")
				var key string
				fmt.Scan(&key)
				settingJson.Consumerkey = key
			} else if regNum == "2" {
				fmt.Print("Input your API secret: ")
				var key string
				fmt.Scan(&key)
				settingJson.Consumersecret = key
			} else if regNum == "3" {
				fmt.Print("Input your Access token: ")
				var key string
				fmt.Scan(&key)
				settingJson.Accesstoken = key
			} else if regNum == "4" {
				fmt.Print("Input your Access secret: ")
				var key string
				fmt.Scan(&key)
				settingJson.Accesssecret = key
			}
		}

		settingWrite, err := json.Marshal(settingJson)
		if err != nil {
			panic(err)
		}
		outputSetting := []byte(settingWrite)
		err = ioutil.WriteFile(hidden2eetFile, outputSetting, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		setting, err := ioutil.ReadFile(hidden2eetFile)
		if err != nil {
			panic(err)
		}
		var settingJson SettingHidden2eet
		err = json.Unmarshal(setting, &settingJson)
		if err != nil {
			panic(err)
		}
		consumerKey := settingJson.Consumerkey
		consumerSecret := settingJson.Consumersecret
		accessToken := settingJson.Accesstoken
		accessSecret := settingJson.Accesssecret
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
