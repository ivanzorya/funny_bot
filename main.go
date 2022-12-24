package main

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)


type JokeResponse struct {
	Error  bool `json:"error"`
	Joke   string `json:"joke"`
}

var buttons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "Get Joke"},
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func getJoke() string {
	c := http.Client{}
	resp, err := c.Get("https://v2.jokeapi.dev/joke/Programming?type=single")
	if err != nil {
		return "jokes API not responding"
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	
	joke := JokeResponse{}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		log.Fatal(err)
		return "Joke error"
	}

	return joke.Joke
}

func main() {

	bot, err := tgbotapi.NewBotAPI(goDotEnvVariable("APIKEY"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		var message tgbotapi.MessageConfig
		log.Println("received text: ", update.Message.Text)

		switch update.Message.Text {
		case "Get Joke":
			message = tgbotapi.NewMessage(update.Message.Chat.ID, getJoke())
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Press "Get Joke" to receive joke`)
		}

		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

		bot.Send(message)
	}
}
