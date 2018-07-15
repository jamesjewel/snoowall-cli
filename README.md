# Snoowall

Snoowall is a utility which updates your desktop background from the Reddit subreddits you specify. It is written in *golang*, by a nifty young developer with a supposedly girly name.

## Usage

```bash
snoowall -sub earthporn -top -allow-nsfw
```

### Arguments

`-sub`: Name of the subreddit to fetch images from. If omitted, defaults to /r/wallpaper.

`-top`: Fetches the top image instead of a random one.

`-allow-nsfw`: Gives a pass to NSFW content that is blocked by default.
