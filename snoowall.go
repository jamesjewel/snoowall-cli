package main

/*
	1. Imgur Integration (Done)
	2. Wallpaper Setting (Done)
	3. Debugging (Done)
	4. Optimization: Switch to indexed Lurker implementation (Done)
						  Change file name to Name (Done)
						  Bring binary data files (Done)
	5. Command Line Options (Done)
	6. Logging (Done) -> Improvements
	7. Configuration File
	8. PNG Problem (Done)
	9. NSFW (Done)
	10. Develop Syncing (Done)
	11. Auto Syncing
	12. Sync Randomizer

	NaN. Graphical User Interface
*/

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/reujab/wallpaper"
	"github.com/turnage/graw/reddit"
)

var loc string
var index int
var subreddit string
var top, nsfw, sync bool
var logfile = "LOG.log"

func saveWall(filename string, b []byte) error {
	err := ioutil.WriteFile(filename, b, 0600)
	return err
}

func setWall(file string) error {
	background, err := wallpaper.Get()
	if err != nil {
		fmt.Println("[DEBUG] Can't find previous wallpaper:", err)
	}
	fmt.Println("Current wallpaper:", background)

	err = wallpaper.SetFromFile(file)
	if err == nil {
		fmt.Println("Updated Wallpaper:", file)
		log.Println("[INFO] Updated wallpaper:", file)
	}
	return err
}

type saveData struct {
	Time       time.Time
	Subreddit  string
	Permalinks []string
}

/*
	The main code
*/
func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// collect flags
	flag.StringVar(&subreddit, "sub", "wallpaper", "Specify the subreddit to fetch wallpapers from.")
	flag.BoolVar(&top, "top", false, "Select the top wallpaper instead of a random one.")
	flag.BoolVar(&nsfw, "allow-nsfw", false, "Gives a pass to NSFW content that is blocked by default.")
	flag.BoolVar(&sync, "sync", false, "Syncs the local database with reddit.")
	flag.IntVar(&index, "index", 1, "Post index (0-99)")
	flag.Parse()
	fmt.Printf("[DEBUG] Arguments: sub:%s;  top:%t;  index:%d;  sync:%t;  allow-nsfw:%t;  tail:%v\n", subreddit, top, index, sync, nsfw, flag.Args())

	// setup logging
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// initialize script to API
	rate := 5 * time.Second
	script, err := reddit.NewScript("graw:snoowall:0.3.1 by /u/psychemerchant", rate)
	if err != nil {
		log.Fatalln("[FATAL] Failed to create script handle: ", err)
		return
	}

	// generating cache
	if sync == true {
		if _, err := os.Stat("cache"); os.IsNotExist(err) {
			os.Mkdir("cache", os.ModePerm)
		}
		harvest, err := script.Listing(fmt.Sprintf("/r/%s", subreddit), "")
		if err != nil {
			log.Fatalf("[FATAL] Failed to fetch /r/%s: %s", subreddit, err)
			return
		}
		var subdata saveData
		subdata.Time = time.Now()
		subdata.Subreddit = subreddit
		var postPermalinks []string
		length := len(harvest.Posts)
		if length == 0 {
			log.Println("[ERROR]: No posts! Subreddit might not exist.")
			return
		}
		for i := 0; i < length; i++ {
			post := harvest.Posts[i]
			postPermalinks = append(postPermalinks, post.Permalink)
		}
		subdata.Permalinks = postPermalinks

		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)
		err = enc.Encode(subdata)
		if err != nil {
			log.Println("[ERROR]: Encoding error", err)
		}
		ioutil.WriteFile(fmt.Sprintf("%s/%s", "cache", subreddit), buff.Bytes(), 0600)
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", "cache", subreddit))
	if err != nil {
		log.Fatalln("[ERROR]: Cache file reatimed error.")
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	var cursubdata saveData
	dec.Decode(&cursubdata)
	var postPermalink string
retry:
	if top == true {
		postPermalink = cursubdata.Permalinks[index]
	} else if top == false {
		postPermalink = cursubdata.Permalinks[rand.Intn(len(cursubdata.Permalinks))]
	}
	// thread fetching
	var post *reddit.Post
	post, err = script.Thread(postPermalink)
	if err != nil {
		log.Fatalf("[FATAL] Failed to fetch post: %s err:%s", postPermalink, err)
	}

	fmt.Printf("[Title]: %s\n[URL]: %s\n", post.Title, post.URL)
	resp, err := http.Get(post.URL)
	filetype := post.URL[len(post.URL)-4:]
	if filetype != ".jpg" && filetype != ".png" {
		log.Println("[ERROR] Not an image.")
		goto retry
	}
	fmt.Println("[DEBUG] Image Type:", filetype)
	if err != nil {
		log.Println("[ERROR]: Couldn't fetch resource:", post.URL, err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	loc := fmt.Sprintf("%s/%s", os.Getenv("HOME"), "Pictures/Wallpapers/")
	filename := fmt.Sprintf("%s%s_%s%s", loc, subreddit, post.ID, filetype)
	err = saveWall(filename, body)
	if err != nil {
		log.Println("[ERROR] Wallpaper saving error:", err)
	}
	err = setWall(filename)
	if err != nil {
		log.Println("[ERROR] Wallpaper setting error:", err)
		return
	}
}
