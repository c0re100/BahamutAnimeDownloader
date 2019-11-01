package main

import (
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/MercuryEngineering/CookieMonster"
)

func (h *bahamut) readCookieFile() {
    if h.cookie != "" {
        _, err := os.Stat(h.cookie)
        isErr("Cookie file error -", err)

        cookies, err := cookiemonster.ParseFile(h.cookie)
        isErr("Parsing Cookie file failed -", err)

        if len(cookies) != 0 {
            for _, ck := range cookies {
                h.rawCookie += ck.Name + "=" + ck.Value + "; "
            }
        } else {
            data, err := ioutil.ReadFile(h.cookie)
            isErr("Read cookie file failed -", err)
            h.rawCookie = string(data)
            h.rawCookie = strings.TrimRight(h.rawCookie, "\n\r")
        }
    }
}

func (h *bahamut) refreshCookie() {
    if h.rawCookie != "" {
        req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/", nil)
        isErr("Create request failed - ", err)

        req.Header.Add("cookie", h.rawCookie)
        req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
        req.Header.Add("referer", "https://ani.gamer.com.tw/")
        resp, err := http.DefaultClient.Do(req)
        isErr("Refresh Cookie failed -", err)

        for _, ck := range resp.Cookies() {
            if ck.Value == "deleted" {
                isErr("Set Cookie failed -", errors.New("Seems your cookie is expired."))
            } else {
                h.setCookie(ck.Name, ck.Value)
            }
        }

        if len(resp.Cookies()) > 0 {
            ioutil.WriteFile(h.cookie, []byte(h.rawCookie), 0755)
            fmt.Println("Cookie refreshed.")
        } else {
            fmt.Println("No need refresh cookie.")
        }
    }
}

func (h *bahamut) setCookie(name, value string) {
    header := http.Header{}
    header.Add("Cookie", h.rawCookie)
    request := http.Request{Header: header}

    // clear current cookie first
    h.rawCookie = ""

    for _, f := range request.Cookies() {
        if f.Name == name {
            f.Value = value
        }
        h.rawCookie += f.Name + "=" + f.Value + "; "
    }
}
