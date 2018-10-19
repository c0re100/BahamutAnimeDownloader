package main

import (
    "flag"
    "time"

    "gopkg.in/cheggaaa/pb.v1"
)

type bahamut struct {
    sn        string
    mUrl      string
    plName    string
    chuckList []string
    bar       *pb.ProgressBar
    startTime int64
    deviceId  string
    cookie    string
    quality   string
    res       QCChoirce
    tmp       string
}

type QCChoirce struct {
    s360  string
    s540  string
    s720  string
    s1080 string
}

func main() {
    envCheck()

    handler := &bahamut{
        sn:        askForSN(),
        startTime: time.Now().Unix(),
    }

    flag.StringVar(&handler.quality, "quality", "720p", "set resolution")
    flag.Parse()

    handler.getDeviceId()
    handler.gainAccess()
    handler.Unlock()
    handler.CheckLock()
    handler.Unlock()
    handler.Unlock()
    handler.startAd()
    handler.skipAd()
    handler.VideoStart()
    handler.checkNoAd()
    handler.getM3U8()
    handler.parseMasterList()
    handler.downloadM3U8()
    handler.parseM3U8()
    handler.start()
    handler.mergeChunk()
    handler.cleanUp()
}
