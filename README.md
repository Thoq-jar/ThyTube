# ThyTube

Download YouTube videos locally and watch them without trackers or ads.

## Setup

Prerequisites:

- Go(lang) v1.25.0
- yt-dlp

1. Install the latest version of [yt-dlp](https://github.com/yt-dlp/yt-dlp) globally as `yt-dlp` or create a symbolic
   link.
   > Make sure to keep yt-dlp up to date as YouTube may patch it.\
   Do not complain to them if ThyTube doesn't work!
2. Clone this repository and enter the directory
3. Run `go run .` to start the server
4. Visit `http://localhost:9595`

## Navigating the UI

### Downloading

Enter a YouTube watch url into the search bar and either press `Get` or enter

### Watching

From the homepage, click a listed video
To return home, click the `ThyTube: Watch` title.

### Errors

Errors may occur, they will be shown inplace of expected content.