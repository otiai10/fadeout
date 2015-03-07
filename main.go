package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type conf struct {
	App struct {
		ConsumerKey    string `json:"consumer_key"`
		ConsumerSecret string `json:"consumer_secret"`
	} `json:"app"`
	User struct {
		Token       string `json:"token"`
		TokenSecret string `json:"token_secret"`
	} `json:"user"`
}

var (
	api      *anaconda.TwitterApi
	_conf    conf
	duration = 60 * 3
	min      = 200
)

const (
	nantonakuSleep = 3 * time.Second
)

func init() {
	f, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(f, &_conf); err != nil {
		log.Fatalln(err)
	}

	anaconda.SetConsumerKey(_conf.App.ConsumerKey)
	anaconda.SetConsumerSecret(_conf.App.ConsumerSecret)
	api = anaconda.NewTwitterApi(_conf.User.Token, _conf.User.TokenSecret)
}

func main() {
	for {
		work()
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func work() {
	id, err := choose()
	if err != nil {
		failure(err)
		return
	}
	user, err := leave(id)
	if err != nil {
		failure(err)
		return
	}
	success(user)
	return
}

func choose() (int64, error) {
	cursor, err := api.GetFriendsIds(nil)
	if err != nil {
		return 0, err
	}
	if len(cursor.Ids) < min {
		log.Println("[FINISHED] これ以上はleaveしない")
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())
	return cursor.Ids[rand.Intn(len(cursor.Ids))], nil
}

func leave(id int64) (user anaconda.User, err error) {

	user, err = api.BlockUserId(id, nil)
	if err != nil {
		return user, err
	}
	time.Sleep(nantonakuSleep)
	user, err = api.UnblockUserId(id, nil)
	if err != nil {
		log.Println("[WARN] blocked but failed to unblock, please check this URL", url(user.ScreenName))
		return user, err
	}

	return user, err
}

func success(user anaconda.User) {
	log.Println("[LEAVED]", url(user.ScreenName))
}

func failure(err error) {
	log.Println("[FAIED]", err.Error())
}

func url(username string) string {
	return "https://twitter.com/" + username
}
