# snoowall-cli

`snoowall-cli` is a command-line utility which sets images from Reddit as your desktop background. All you have to do is pass it a subreddit of your choice. It is written in *golang*, by a nifty young developer with a supposedly girly name.

## Usage

Synopsis:
```bash
snoowall-cli [-R, --refresh][-s, --sort][-n, --allow-nsfw] subreddit
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

```



## Examples

```bash
snoowall-cli earthporn 
```
Sets a random image from 'earthporn' as the desktop background.
```bash
snoowall-cli NSFW_Wallpapers --allow-nsfw 
```
Sets the top image from 'NSFW_Wallpapers', even if it is NSFW (which in this case, it clearly is).
```bash
snoowall-cli -Rs hot gmbwallpapers
```
Refreshes cache with 'hot' posts from 'gmbwallpapers' and sets a random image as the desktop background.  

## Installation

Compile from source:
```bash
go get github.com/reujab/wallpaper
go get github.com/turnage/graw/reddit
go get github.com/spf13/pflag
go build snoowall-cli.go
```
For convinience you can add it to the system `PATH` or make a symlink to the `snoowall-cli` executable in `/usr/bin`

YAAAY! I hope Snoo doesn't mind.
