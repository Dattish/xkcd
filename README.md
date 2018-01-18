# xkcd
A wrapper for [xkcds api](https://xkcd.com/json.html) along with some functions for saving comics.  
Comes with a program that fetches xkcd comics.

### xkcd-downloader
If no flags are passed or if only the -d flag is passed the program will default to downloading all comics  
flags:
```bash
-a	Download all comics
  
-d string
    	Set the directory to download the comic[s] to
    	
-l	Download the latest comic
  
-n int
    	Download a specific comic
```