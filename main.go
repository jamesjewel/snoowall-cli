package main

/*
	1. Imgur Integration (Done)
	2. Wallpaper Setting (Done)
	3. Debugging (Done)
	4. Optimization: Switch to indexed Lurker implementation
	5. Command Line Options
	6. Logging
	.
	.
	.
	NaN. Graphical User Interface
*/

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/reujab/wallpaper"
	"github.com/turnage/graw/reddit"
)

var loc = "/home/pro/Dropbox/Code/golang/snoowall/Wallpapers/"
var datafile = "data"
var name = "info.agent"
var path = fmt.Sprintf("%s%s", loc, name)
var index int
var subreddit = "gmbwallpapers"

func saveWall(b []byte) (file string, err error) {
	timestamp := time.Now()
	filename := fmt.Sprintf("%s%s_%s.jpg", loc, subreddit, timestamp.Format("2006-01-02_15-04-05"))
	err = ioutil.WriteFile(filename, b, 0600)
	if err == nil {
		fmt.Println("Wallpaper saved!")
	}
	return filename, err
}

func setWall(file string) error {
	background, err := wallpaper.Get()
	if err != nil {
		fmt.Println("[DEBUG] Can't find previous wallpaper:", err)
	}
	fmt.Println("Current wallpaper:", background)
	err = wallpaper.SetFromFile(file)
	if err == nil {
		fmt.Println("Wallpaper set!")
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

	// var after string
	// bin, _ := ioutil.ReadFile(datafile)
	// after = string(bin)
	// fmt.Println("After:", after)

	harvest, err := script.Listing(fmt.Sprintf("/r/%s", subreddit), "t3_8mo8o8")
	if err != nil {
		fmt.Printf("[DEBUG] Failed to fetch /r/%s: %s", subreddit, err)
		return
	}
	fmt.Println("[DEBUG] Post array length: ", len(harvest.Posts))
	post := harvest.Posts[84]
	// str := fmt.Sprintf("Harvest:\n %#v", harvest.Posts[1])
	// ioutil.WriteFile("harvest", []byte(str), 0600)
	ioutil.WriteFile(datafile, []byte(post.Name), 0600)
	fmt.Println("After:", post.Name)
	fmt.Printf("[Title]: %s\n[URL]: %s\n", post.Title, post.URL)
	// fmt.Printf("[Type]: %s - %s - %s\n", post.Media.OEmbed.Type, post.Media.OEmbed.ProviderName, post.Media.OEmbed.ProviderURL)
	// fmt.Printf("%+v", post)

	resp, err := http.Get(post.URL)
	if err != nil || post.IsRedditMediaDomain == false {
		fmt.Println("[DEBUG]: Couldn't fetch resource:", post.URL, err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	filename, _ := saveWall(body)
	err = setWall(filename)
	if err != nil {
		fmt.Println("[DEBUG] Wallpaper setting error:", err)
		return
	}

}
