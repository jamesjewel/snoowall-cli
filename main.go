package main

import "fmt"
import "github.com/turnage/graw/reddit"

var loc = "/home/pro/Dropbox/Code/golang/snoowall/"
var name = "info.agent"
var path = fmt.Sprintf("%s%s", loc, name)
var subreddit = "/r/EarthPorn"

func main() {
	fmt.Println("Hello Snoowall!")
	fmt.Println(path)
	bot, err := reddit.NewBotFromAgentFile(path, 0)
	fmt.Println(bot, err)
	harvest, err := bot.Listing(subreddit, "")
	fmt.Println(harvest, err)
	if err != nil {
		fmt.Println("Failed to fetch /r/golang: ", err)
	}
	for _, post := range harvest.Posts[:5] {
		fmt.Printf("[%s] posted [%s]\n", post.Author, post.Title)
	}
}
