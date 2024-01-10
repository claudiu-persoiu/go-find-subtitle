# go-find-subtitle

# Go Find Subtitle
Movie subtitle finder

The finder can work in two modes, as a command line or integrated with [Transmission Torrent](https://transmissionbt.com/) client.

## Installation

1. Create a profile on https://www.opensubtitles.com/
2. Login into the new account
3. Click on profile and then "Api Consumers"
4. NEW CONSUMER
5. Add a name, uncheck "Allow anonymous downloads" and check "Under dev" and SAVE
6. Copy the API key
7. Create a fine in ```/etc/go-find-subtitle.json``` with the structure:
```
   {
        "opensubititles": {
            "user": "<<username>>",
            "pass": "<<pass>>",
            "key": "<<api key>>",
            "languages": "<<lang1,lang2,lang3>>"
        }
   }
```

## Running

### As a command line tool
```
$ ./gofindsub-linux-amd64 path_to_movie
```

### Integrate with Transmission Torrent

After completing the steps above, stop transmission process:
```
$ sudo service transmission-daemon stop
```

Edit the config file located at: **/etc/transmission-daemon/settings.json**

NOTE: Path to config may be different for your setup, for Ubuntu desktop the path is **~/.config/transmission/settings.json**


Find the lines with:
```
"script-torrent-done-enabled": false,
"script-torrent-done-filename": "",
```
and replace with:
```
"script-torrent-done-enabled": true,
"script-torrent-done-filename": "/path/to/installation/gofindsub",
```

Start the server back:
```
$ sudo service transmission-daemon start
```

NOTE: In case you need to debug Transmission installation the output is redirected to default syslog file.

### Known issues

If more files finish in the same time some of them may not be able to get translations because of multiple login tokens in the same time.