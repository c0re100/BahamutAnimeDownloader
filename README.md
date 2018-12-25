# Bahamut Anime Downloader

Download anime from [動畫瘋](https://ani.gamer.com.tw/)

## Installation
This program uses ffmpeg to merge the chunk file and output the video with the original codecs.

Therefore, ffmpeg need to be available in `PATH` or same directory.

### Compilation
```
git clone https://github.com/c0re100/BahamutAnimeDownloader.git
cd BahamutAnimeDownloader
go build
```

#### Usage
For example, you want to download [驚爆危機 Invisible Victory 第12集](https://ani.gamer.com.tw/animeVideo.php?sn=10434)

Open the downloader and type the sn id (animeVideo.php?sn=`10434`)

Then you can downloading the videos as fast as possible (Depending on your network).

![example](https://i.imgur.com/BpuQckG.png)

##### Video Quality Option
You can start a program with command line argument to select quality, otherwise the default quality is 720p.

`AniDownloader -quality="720p"`

Quality
* 360p
* 540p
* 720p
