# Bahamut Anime Downloader

Download anime from [動畫瘋](https://ani.gamer.com.tw/)

## Installation
This program uses ffmpeg to merge the chunk file and output the video with the original codecs.

Therefore, ffmpeg need to be available in `PATH` or same directory.

## Compilation
Using Makefile
```
git clone https://github.com/c0re100/BahamutAnimeDownloader.git
cd BahamutAnimeDownloader
make deps  # Install the dependences
make
```
Or Manual compile
```
git clone https://github.com/c0re100/BahamutAnimeDownloader.git
cd BahamutAnimeDownloader
go build

# You need to make sure the dependences have been installed
```

## Usage
For example, you want to download [驚爆危機 Invisible Victory 第12集](https://ani.gamer.com.tw/animeVideo.php?sn=10434)

Open the downloader and type the sn id (animeVideo.php?sn=`10434`)

Then you can downloading the videos as fast as possible (Depending on your network).

![example](https://i.imgur.com/BpuQckG.png)

## Command line arguments

### Select sn

`AniDownloader -s 10434`

### Video Quality Option
You can start a program with command line argument to select quality, otherwise the default quality is 720p.

Quality
* 360p
* 540p
* 720p
* 1080p # Only for bahamut premium member, and needed to be provided Cookies file

`AniDownloader -quality="720p"` or `Anidownloader -q 720p`

### Cookies
If you are the *premium member* and want to download 1080p video

Using -c option to provide your Cookies file

`AniDownloader -c cookies.txt -q 1080p`

Cookies file can be Raw cookie format or Nestscape cookie format
```
1. Raw cookies:

    name = value; name = value; ...

2. Nestscape format:

    # Netscape HTTP Cookie File

    .example.org	TRUE	/	FALSE	1552060831	remember_me	true
    .example.org	TRUE	/	FALSE	1552060831	APISID	DijdSAOAjgwijnhFMndsjiejfdSDNSgfsikasASIfgijsowITITeoknsd
    .example.org	TRUE	/	FALSE	1552060831	static_files	iy1aBf1JhQR

```

Then example of watching 1080p video with Cookies

![example](https://i.imgur.com/aoVMUVP.png)

### Output

Using -o option to set output path

`AniDownloader -s 10434 -q 720p -o .\Anime`