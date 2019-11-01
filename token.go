package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "math/rand"
    "net/http"
    "time"
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
        req.Header.Add("cookie", h.rawCookie)
    }

    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
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
            h.rawCookie = "nologinuser=" + ck.Value
        }
    }
}

func (h *bahamut) gainAccess() {
    resp := h.request("gainAccess", "https://ani.gamer.com.tw/ajax/token.php?adID=0&sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12))

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

func (h *bahamut) checkPremium() {
    resp := h.request("checkPremium", "https://ani.gamer.com.tw/ajax/token.php?sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12))

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    isErr("Read response failed -", err)

    var bahaData map[string]interface{}
    err = json.Unmarshal(body, &bahaData)
    isErr("Parse json failed -", err)

    if bahaData["vip"] != nil && bahaData["vip"].(bool) == true {
        h.isPremium = true
        fmt.Println("You're Premium user")
    } else {
        h.isPremium = false
        fmt.Println("You're not Premium user")
    }
}

func (h *bahamut) checkNoAd() {
    resp := h.request("checkNoAd", "https://ani.gamer.com.tw/ajax/token.php?sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12))

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
    h.request("startAd", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=194699").Body.Close()
}

func (h *bahamut) skipAd() {
    h.request("skipAd", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=194699&ad=end").Body.Close()
}

func (h *bahamut) unlock() {
    h.request("unlock", "https://ani.gamer.com.tw/ajax/unlock.php?sn="+h.sn+"&ttl=0").Body.Close()
}

func (h *bahamut) checkLock() {
    h.request("checkLock", "https://ani.gamer.com.tw/ajax/checklock.php?device="+h.deviceId+"sn="+h.sn).Body.Close()
}

func (h *bahamut) videoStart() {
    h.request("videoStart", "https://ani.gamer.com.tw/ajax/videoStart.php?sn="+h.sn).Body.Close()
}
