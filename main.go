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
    rawCookie string
    quality   string
    res       string
    tmp       string
    isPremium bool
}

func main() {
    envCheck()

    handler := &bahamut{
        startTime: time.Now().Unix(),
    }

    flag.StringVar(&handler.sn, "sn", "", "set sn")
    flag.StringVar(&handler.sn, "s", "", "set sn(shorthand)")
    flag.StringVar(&handler.cookie, "cookie", "", "cookie") // raw cookie
    flag.StringVar(&handler.cookie, "c", "", "cookie(shorthand)")
    flag.StringVar(&handler.quality, "quality", "720p", "set resolution")
    flag.StringVar(&handler.quality, "q", "720p", "set resolution(shorthand)")
    flag.Parse()

    handler.askForSN()
    handler.readCookieFile()
    handler.refreshCookie()
    handler.getDeviceId()
    handler.gainAccess()
    handler.checkPremium()
    handler.unlock()
    handler.checkLock()
    handler.unlock()
    handler.unlock()
    if !handler.isPremium {
        handler.startAd()
        time.Sleep(8 * time.Second)
        handler.skipAd()
    }
    handler.videoStart()
    handler.checkNoAd()
    handler.getM3U8()
    handler.parseMasterList()
    handler.downloadM3U8()
    handler.parseM3U8()
    handler.start()
    handler.mergeChunk()
    handler.cleanUp()
}
