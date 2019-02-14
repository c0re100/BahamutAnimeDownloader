package main

import (
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path"
    "strconv"
    "strings"

    "github.com/grafov/m3u8"
)

func (h *bahamut) askForSN() {
    var sn string

    sn = h.sn
    for sn == "" {
        fmt.Println("Please type anime sn id: ")
        fmt.Scanln(&sn)
    }

    // extract URL input
    _, err := strconv.Atoi(sn)
    if err != nil {
        if strings.Contains(sn, "animeVideo") {
            u, err := url.Parse(sn)
            isErr("Extract url failed -", err)

            qStr, err := url.ParseQuery(u.RawQuery)
            isErr("Parse query string failed -", err)
            if qStr.Get("sn") == "" {
                isErr("Extract sn id failed -", errors.New("sn id not found"))
            }
            sn = qStr.Get("sn")
        } else {
            isErr("Please try again -", errors.New("format not corret, Allow: numeric / URL"))
        }
    }

    h.sn = sn
}

func (h *bahamut) parseMasterList() {
    req, err := http.NewRequest("GET", h.mUrl, nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    isErr("Get m3u8 playlist failed -", err)

    defer resp.Body.Close()
    p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
    isErr("Parse m3u8 playlist failed -", err)

    switch listType {
    case m3u8.MASTER:
        masterpl := p.(*m3u8.MasterPlaylist)
        for _, list := range masterpl.Variants {
            switch list.Resolution {
            case "640x360":
                h.res.s360 = strings.Split(list.URI, "?")[0]
            case "960x540":
                h.res.s540 = strings.Split(list.URI, "?")[0]
            case "1280x720":
                h.res.s720 = strings.Split(list.URI, "?")[0]
            case "1920x1080":
                h.res.s1080 = strings.Split(list.URI, "?")[0]
            }
        }
    }
    fmt.Println("Get m3u8 playlist.")
}

func (h *bahamut) parseM3U8() {
    f, err := os.Open(h.tmp + "/" + h.plName)
    isErr("Failed to read m3u8 playlist -", err)

    pl, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
    isErr("Parse m3u8 playlist failed -", err)

    switch listType {
    case m3u8.MEDIA:
        mediapl := pl.(*m3u8.MediaPlaylist)
        newPl, e := m3u8.NewMediaPlaylist(mediapl.WinSize(), mediapl.Count())
        isErr("Creating new m3u8 playlist failed - ", e)

        prefix := strings.Split(h.mUrl, "playlist.m3u8")[0]
        newPl.SetDefaultKey(mediapl.Key.Method, h.downloadKey(mediapl.Key.URI), "", "", "")
        for _, chuck := range mediapl.Segments {
            if chuck != nil {
                h.chuckList = append(h.chuckList, prefix+chuck.URI)
                newPl.Append(strings.Split(path.Base(chuck.URI), "?")[0], chuck.Duration, "")
            }
        }
        newPl.Close()

        ioutil.WriteFile(h.tmp+"/"+h.plName, newPl.Encode().Bytes(), 0755)
        fmt.Println("All segments parsed! Download is starting...")
    }
}
