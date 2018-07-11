package main

/*
	1. Imgur Integration (Done)
	2. Wallpaper Setting
	3. Debugging
*/

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/turnage/graw/reddit"
)

var loc = "/home/pro/Dropbox/Code/golang/snoowall/Wallpapers/"
var name = "info.agent"
var path = fmt.Sprintf("%s%s", loc, name)
var subreddit = "sexy"

func saveWall(b []byte) error {
	timestamp := time.Now()
	filename := fmt.Sprintf("%s%s_%s.jpg", loc, subreddit, timestamp.Format("2006-01-02_15-04-05"))
	err := ioutil.WriteFile(filename, b, 0600)
	if err == nil {
		fmt.Println("Saved")
	}
	return err
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	rate := 5 * time.Second
	script, err := reddit.NewScript("graw:snoowall:0.3.1 by /u/psychemerchant", rate)
	if err != nil {
		fmt.Println("[DEBUG] Failed to create script handle: ", err)
		return
	}
	m := make(map[string]string, 0)
	m["Content-type"] = "image/jpeg"
	harvest, err := script.ListingWithParams(fmt.Sprintf("/r/%s", subreddit), m)
	if err != nil {
		fmt.Println("[DEBUG] Failed to fetch /r/:", subreddit, err)
		return
	}
	post := harvest.Posts[rand.Intn(30)]
	fmt.Printf("[Title]: %s\n[URL]: %s\n", post.Title, post.URL)
	resp, err := http.Get(post.URL)
	if err != nil {
		fmt.Println("[DEBUG]: Couldn't fetch resource:", post.URL, err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	saveWall(body)
}
