# Snoowall

`snoowall` is a command-line utility which sets images from Reddit as your desktop background. All you have to do is pass it a subreddit of your choice. It is written in *golang*, by a nifty young developer with a supposedly girly name.

## Usage

```bash\
snoowall [OPTIONS]
```

Arguments

`-sub` : Name of the subreddit to fetch images from. If omitted, defaults to 'wallpaper'.

`-top` : Fetches the top image instead of a random one.

`-allow-nsfw` : Gives a pass to NSFW content that is blocked by default.

`-sync` :  Manually refresh the post cache.

## Examples

```bash
snoowall -sub earthporn 
```
Sets a random image from 'earthporn' as the desktop background.
```bash
snoowall -sub NSFW_Wallpapers -top -allow-nsfw 
```
Sets the top image from 'NSFW_Wallpapers', even if it is NSFW (which in this case, it clearly is).
```bash
snoowall -sub gmbwallpapers -sync
```
Syncs new posts from 'gmbwallpapers' and sets a random image as the desktop background.  

YAAAY! 