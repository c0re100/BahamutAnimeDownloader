package main

import (
    "encoding/json"
    "errors"
    "strings"
    "fmt"
    "io/ioutil"
    "os"
    "math/rand"
    "net/http"
    "time"
    "github.com/MercuryEngineering/CookieMonster"
)

func randomString(num int) string {
    rand.Seed(time.Now().UTC().UnixNano())
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
    result := make([]byte, num)
    for i := 0; i < num; i++ {
        result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}

func (h *bahamut) getDeviceId() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/getdeviceid.php?id="+h.deviceId, nil)
    isErr("Create request failed - ", err)

    if h.cookie != "" {
        _, err := os.Stat(h.cookie)
        isErr("Cookie file error -", err)

        cookies, err := cookiemonster.ParseFile(h.cookie)
        isErr("Parsing Cookie file failed -", err)

        if (len(cookies) != 0) {
            for _, ck := range cookies {
                h.cookie += ck.Name + "=" + ck.Value + "; "
            }
        } else {
            data, err := ioutil.ReadFile(h.cookie)
            isErr("Read cookie file failed -", err)
            h.cookie = string(data)
            h.cookie = strings.TrimSuffix(h.cookie, "\n")
        }
        req.Header.Add("cookie", h.cookie)
    }

    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    isErr("Get Device ID failed -", err)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    isErr("Read response failed -", err)

    var bahaData map[string]interface{}
    err = json.Unmarshal(body, &bahaData)
    isErr("Parse json failed -", err)

    h.deviceId = bahaData["deviceid"].(string)

    for _, ck := range resp.Cookies() {
        if ck.Name == "nologinuser" {
            if h.cookie != "" {
                isErr("Can't using Cookie -", errors.New("Your Cookie may be incorrect. "))
            }
            h.cookie = "nologinuser=" + ck.Value
        }
    }
}

func (h *bahamut) gainAccess() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/token.php?adID=0&sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12), nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    isErr("Get Token failed -", err)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    isErr("Read response failed -", err)

    var bahaData map[string]interface{}
    err = json.Unmarshal(body, &bahaData)
    isErr("Parse json failed -", err)

    if bahaData["error"] != nil {
        isErr("Something happened -", errors.New("Where are you from? "))
    } else {
        fmt.Println("Gained access.")
    }
}

func (h *bahamut) checkNoAd() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/token.php?sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12), nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    isErr("Get Token failed -", err)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    isErr("Read response failed -", err)

    var bahaData map[string]interface{}
    err = json.Unmarshal(body, &bahaData)
    isErr("Parse json failed -", err)

    if bahaData["time"] != nil {
        if bahaData["time"].(float64) == 1 {
            fmt.Println("Adaway.")
        } else {
            isErr("Something happened -", errors.New("Ads not away? "))
        }
    } else {
        isErr("Something happened -", errors.New("Where are you from? "))
    }
}

func (h *bahamut) startAd() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=194699", nil)
    isErr("Create skipAd request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    _, err = http.DefaultClient.Do(req)
    isErr("Start ads failed -", err)
}

func (h *bahamut) skipAd() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=194699&ad=end", nil)
    isErr("Create skipAd request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    _, err = http.DefaultClient.Do(req)
    isErr("Skip ads failed -", err)
}

func (h *bahamut) Unlock() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/unlock.php?sn="+h.sn+"&ttl=0", nil)
    isErr("Create Unlock request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    _, err = http.DefaultClient.Do(req)
    isErr("Unlock failed -", err)
}

func (h *bahamut) CheckLock() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/checklock.php?device="+h.deviceId+"sn="+h.sn, nil)
    isErr("Create Check Lock request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    _, err = http.DefaultClient.Do(req)
    isErr("Check Lock failed -", err)
}

func (h *bahamut) VideoStart() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/videoStart.php?sn="+h.sn, nil)
    isErr("Create Check Lock request failed - ", err)

    req.Header.Add("cookie", h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    _, err = http.DefaultClient.Do(req)
    isErr("Video Start failed -", err)
}
