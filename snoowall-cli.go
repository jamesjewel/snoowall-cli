package main

/*
	v.0.4.0
	Features:
	- Sorting options
	- Uses reddit's random listing, disables local randomizer
	- Auto refeshing based on system time
*/

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/reujab/wallpaper"
	flag "github.com/spf13/pflag"
	"github.com/turnage/graw/reddit"
)

var subreddit, sort string
var nsfw, refresh bool

var rcount = 0

func saveWall(filename string, b []byte) error {
	err := ioutil.WriteFile(filename, b, 0600)
	return err
}

func setWall(file string) error {
	// background, err := wallpaper.Get()
	// if err != nil {
	// 	fmt.Println("[INFO] Can't find previous wallpaper:", err)
	// }

	err := wallpaper.SetFromFile(file)
	if err == nil {
		fmt.Println("Updated Wallpaper:", file)
		log.Println("[INFO] Updated wallpaper:", file)
	}
	return err
}

type saveData struct {
	Time      time.Time
	Subreddit string
	Info      []postMeta
}
type postMeta struct {
	Title     string
	ID        string
	Permalink string
	NSFW      bool
	URL       string
}

/*
	The main code
*/
func main() {
	// collect flags
	flags := flag.NewFlagSet("snoowall-cli", flag.ExitOnError)
	flags.StringVarP(&sort, "sort", "s", "hot", "Specify the sorting method.")
	flags.BoolVarP(&nsfw, "allow-nsfw", "n", false, "Gives a pass to NSFW content that is blocked by default.")
	flags.BoolVarP(&refresh, "refresh", "R", false, "Refreshes the local post cache from Reddit.")

	flags.Parse(os.Args[1:])
	// flags.MarkDeprecated("refersh"), deprecate when auto-refresh is implemented.
	subreddit = flags.Args()[0]
	if sort != "hot" && sort != "top" && sort != "new" && sort != "controversial" {
		fmt.Println("Invalid sort option.")
		return
	}

	// setup logging
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("[DEBUG] subreddit:%s;\nflags:sort=%s; refresh=%t; allow-nsfw=%t; tail:%v\n", subreddit, sort, refresh, nsfw, flags.Args())

	// initialize script to API
	rate := 5 * time.Second
	script, err := reddit.NewScript("graw:snoowall:0.4.0 by /u/psychemerchant", rate)
	if err != nil {
		fmt.Println("Fatal error! Check log file for more info.")
		log.Fatalln("[FATAL] Failed to create script handle: ", err)
	}

	// if cache does not exist, sync
	syncLoc := fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".cache/snoowall-cli")
	cachefile := fmt.Sprintf("%s/%s_%s", syncLoc, subreddit, sort)
	if _, err := os.Stat(cachefile); os.IsNotExist(err) {
		refresh = true
	}

	// generating cache
	if refresh == true {
		if _, err := os.Stat(syncLoc); os.IsNotExist(err) {
			os.MkdirAll(syncLoc, os.ModePerm)
		}
		fmt.Printf("Saving post cache... /r/%s to %s\n", subreddit, cachefile)

		harvest, err := script.Listing(fmt.Sprintf("/r/%s/%s", subreddit, sort), "")
		if err != nil {
			log.Printf("[FATAL] Failed to fetch /r/%s/%s: %s", subreddit, sort, err)
			fmt.Println("Subreddit does not exist.")
			return
		}
		var subdata saveData
		subdata.Time = time.Now()
		subdata.Subreddit = subreddit
		subdata.Info = make([]postMeta, 0)

		length := len(harvest.Posts)
		for i := 0; i < length; i++ {
			post := harvest.Posts[i]
			subdata.Info = append(subdata.Info, postMeta{post.Title, post.ID, post.Permalink, post.NSFW, post.URL})
		}

		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)
		err = enc.Encode(subdata)
		if err != nil {
			fmt.Println("Fatal error! Check log file for more info.")
			log.Fatalln("[FATAL]: Encoding error", err)
		}
		err = ioutil.WriteFile(cachefile, buff.Bytes(), 0600)
		if err != nil {
			fmt.Println("Fatal error! Check log file for more info.")
			log.Fatalln("[FATAL]: Cache saving error", err)
		}
		log.Printf("[INFO] Synced /r/%s to %s", subreddit, cachefile)

	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s", cachefile))
	if err != nil {
		fmt.Println("Fatal error! Check log file for more info.")
		log.Fatalln("[ERROR]: Cache file read error.")
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	var cursubdata saveData
	cursubdata.Info = make([]postMeta, 0)
	err = dec.Decode(&cursubdata)
	var post postMeta
retry:

	rand.Seed(time.Now().UTC().UnixNano())
	post = cursubdata.Info[rand.Intn(len(cursubdata.Info))]

	// if allow-nsfw - false, nsfw check, retry
	if nsfw == false {
		if post.NSFW == true {
			if rcount == 3 {
				fmt.Println("Post is NSFW. Try an SFW subreddit.")
				return
			}
			// if top == false {
			fmt.Println("Post is NSFW. Retrying...")
			rcount++
			goto retry
			// } else {
			// fmt.Println("[DEBUG] Top post is NSFW.")
			// return
		}
	}

	fmt.Printf("Title: %s\nURL: %s\n", post.Title, post.URL)
	resp, err := http.Get(post.URL)
	filetype := post.URL[len(post.URL)-4:]
	if filetype != ".jpg" && filetype != ".png" {
		log.Println("[ERROR] Not an image.")
		fmt.Println("Post is not an image. Retrying...")
		goto retry
	}
	if err != nil {
		log.Println("[ERROR]: Couldn't fetch resource:", post.URL, err)
		fmt.Println("Couldn't fetch resource:", post.URL, err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	// if save=true, save in user's home directory, else save in cache
	var loc string
	loc = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".cache/snoowall-cli/")
	filename := fmt.Sprintf("%s%s_%s%s", loc, subreddit, post.ID, filetype)
	err = saveWall(filename, body)
	if err != nil {
		fmt.Println("Fatal error! Check log file for more info.")
		log.Fatalln("[ERROR] Wallpaper saving error:", err)
	}

	err = setWall(filename)
	if err != nil {
		fmt.Println("Fatal error! Check log file for more info.")
		log.Println("[ERROR] Wallpaper setting error:", err)
		return
	}
}
