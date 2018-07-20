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
        req.Header.Add("cookie", "nologinuser="+h.cookie)
    }
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
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
            h.cookie = ck.Value
        }
    }
}

func (h *bahamut) gainAccess() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/token.php?adID=0&sn="+h.sn+"&device="+h.deviceId+"&hash="+randomString(12), nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", "nologinuser="+h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
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

    req.Header.Add("cookie", "nologinuser="+h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
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
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=83666", nil)
    isErr("Create skipAd request failed - ", err)

    req.Header.Add("cookie", "nologinuser="+h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    _, err = http.DefaultClient.Do(req)
    isErr("Start ads failed -", err)
}

func (h *bahamut) skipAd() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/videoCastcishu.php?sn="+h.sn+"&s=83666&ad=end", nil)
    isErr("Create skipAd request failed - ", err)

    req.Header.Add("cookie", "nologinuser="+h.cookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    _, err = http.DefaultClient.Do(req)
    isErr("Skip ads failed -", err)
}
