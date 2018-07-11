package main

/*
	1. Imgur Integration
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
var subreddit = "MapPorn"

func saveWall(b []byte) {
	timestamp := time.Now()
	filename := fmt.Sprintf("%s%s_%s.jpg", loc, subreddit, timestamp.Format("2006-01-02_15-04-05"))
	fmt.Println("Saving ...")
	ioutil.WriteFile(filename, b, 0600)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	rate := 5 * time.Second
	script, _ := reddit.NewScript("graw:doc_script:0.3.1 by /u/psychemerchant", rate)
	harvest, _ := script.Listing(fmt.Sprintf("/r/%s", subreddit), "")
	post := harvest.Posts[rand.Intn(30)]
	fmt.Printf("[Title]: %s\n[URL]: %s\n", post.Title, post.URL)
	resp, _ := http.Get(post.URL)
	body, _ := ioutil.ReadAll(resp.Body)
	saveWall(body)

}
