# snoowall-cli

`snoowall-cli` is a command-line utility which sets images from Reddit as your desktop background. All you have to do is pass it a subreddit of your choice. It is written in *golang*, by a nifty young developer with a supposedly girly name.

## Usage

```bash
snoowall-cli [-sub subreddit] [-top] [-allow-nsfw] [-sync]
```

Options : 
```bash
-sub 
            Name of the subreddit to fetch images from. 
            If ommitted, defaults to 'wallpaper'.

-top 
            Fetches the top image instead of a random one.

-allow-nsfw
            Gives a pass to NSFW content that is blocked by default.

-sync
            Manually refresh the post cache.
```



## Examples

```bash
snoowall-cli -sub earthporn 
```
Sets a random image from 'earthporn' as the desktop background.
```bash
snoowall-cli -sub NSFW_Wallpapers -top -allow-nsfw 
```
Sets the top image from 'NSFW_Wallpapers', even if it is NSFW (which in this case, it clearly is).
```bash
snoowall-cli -sub gmbwallpapers -sync
```
Syncs new posts from 'gmbwallpapers' and sets a random image as the desktop background.  

## Installation
Download the compiled executable from here: [snoowall-cli_v.0.3.1](https://www.dropbox.com/s/s1897ki9hrc09c0/snoowall-cli?dl=0)

**OR**

Compile from source:
```bash
go get "github.com/reujab/wallpaper"
go get "github.com/turnage/graw/reddit"
go build snoowall-cli.go
```
For convinience you can add it to the system `PATH` or make a symlink to the `snoowall-cli` executable in `/usr/bin`

YAAAY! I hope Snoo doesn't mind.
