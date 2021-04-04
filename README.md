# snoowall-cli

`snoowall-cli` is a command-line utility which fetches images from Reddit and set them as your desktop background. All you have to do is pass it a subreddit of your choice. It is written in Go, by a nifty young developer.

## Usage

Synopsis:
```bash
snoowall-cli [-R, --refresh][-s, --sort] sort subreddit [-n, --allow-nsfw]
```

Description: 
```bash
subreddit - Name of the subreddit to fetch images from. If ommitted, defaults to 'wallpaper'.
Flags:            
-R, --refresh
            Manually refresh the post cache.
-s, --sort 
            Grabs posts from a list sorted by this mode: hot, top, controversial, new, best.
-n, --allow-nsfw
            Gives a pass to NSFW content that is blocked by default.
-r, --allow-all-sizes
            Bypass checking for wallpaper-suitable aspect ratios in the post title. (Some wallpaper-worthy subreddits don't post the resolution in the title.)
```

## Examples

```bash
$ snoowall-cli earthporn 
```
Sets an image from 'earthporn' as the desktop background.
```bash
$ snoowall-cli animewallpaper --allow-nsfw 
```
Sets an image from 'animewallpaper', even if it is NSFW.

```bash
$ snoowall-cli -s hot skyporn
```
Refreshes cache with 'hot' posts from 'skyporn' and sets a random image as the desktop background. 

## Installation

Compile from source:
```bash
$ git clone https://github.com/flakyhermit/snoowall-cli.git
$ cd snoowall-cli
$ go build .
```
For convenience you can add it to the system `PATH` or make a symlink to the `snoowall-cli` executable in `/usr/bin`

YAAAY! I hope Snoo doesn't mind.
