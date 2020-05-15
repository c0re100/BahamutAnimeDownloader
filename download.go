package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "path"
    "strings"
    "time"

    "github.com/korovkin/limiter"
    "gopkg.in/cheggaaa/pb.v1"
)

func (h *bahamut) getM3U8() {
    req, err := http.NewRequest("GET", "https://ani.gamer.com.tw/ajax/m3u8.php?sn="+h.sn+"&device="+h.deviceId, nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", h.rawCookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    isErr("Get m3u8 playlist failed -", err)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    isErr("Read response failed -", err)

    var bahaData map[string]interface{}
    err = json.Unmarshal(body, &bahaData)
    isErr("Parse json failed -", err)

    if bahaData["src"].(string) != "" {
        h.mUrl = bahaData["src"].(string)
    } else {
        isErr("Please try again -", errors.New("src not found"))
    }
}

func (h *bahamut) downloadM3U8() {
    // Create a temporary directory for storing
    h.tmp = "tmp" + h.sn
    os.Mkdir(h.tmp, 0755)

    var choice string
    h.plName, choice = h.getQuality()
    fmt.Println("Your choice:", choice)

    out, err := os.Create(h.tmp + "/" + h.plName)
    isErr("Create m3u8 playlist failed -", err)

    defer out.Close()
    resp := h.request("downloadM3U8", strings.Replace(h.mUrl, "playlist.m3u8", h.plName, -1))

    defer resp.Body.Close()
    _, err = io.Copy(out, resp.Body)
    isErr("m3u8 playlist save failed -", err)

    fmt.Println("m3u8 playlist downloaded.")
}

func (h *bahamut) downloadKey(keyUrl string) string {
    filename := strings.Split(path.Base(keyUrl), "?")[0]

    out, err := os.Create(h.tmp + "/" + filename)
    isErr("Create key file failed -", err)

    defer out.Close()
    resp := h.request("downloadKey", keyUrl)

    defer resp.Body.Close()
    _, err = io.Copy(out, resp.Body)
    isErr("Key file save failed -", err)
    fmt.Println("m3u8 key file downloaded.")

    return strings.Split(path.Base(keyUrl), "?")[0]
}

func (h *bahamut) downloadChunk(chuckUrl string) bool {
    filename := strings.Split(path.Base(chuckUrl), "?")[0]

    // Check chunk exist or not
    fi, err := os.Stat(h.tmp + "/" + filename)
    if err == nil && fi.Size() != 0 {
        return true
    }

    // Create a chunk file
    out, err := os.Create(h.tmp + "/" + filename)
    isErr("Create "+filename+" failed -", err)

    req, err := http.NewRequest("GET", chuckUrl, nil)
    isErr("Create request failed - ", err)

    req.Header.Add("cookie", h.rawCookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Download "+filename+" file failed -", err)
        fmt.Println("Retrying -", filename)
        out.Close()
        os.Remove(h.tmp + "/" + filename)
        time.Sleep(500 * time.Millisecond)
        return false
    }

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        fmt.Println(filename+" save failed -", err)
        fmt.Println("Retrying -", filename)
        out.Close()
        os.Remove(h.tmp + "/" + filename)
        time.Sleep(500 * time.Millisecond)
        return false
    }
    out.Close()
    return true
}

func (h *bahamut) start() {
    h.bar = pb.StartNew(len(h.chuckList))
    limit := limiter.NewConcurrencyLimiter(64)

    for _, url := range h.chuckList {
        part := url
        limit.Execute(func() {
            for {
                if h.downloadChunk(part) {
                    h.bar.Increment()
                    break
                }
            }
        })
    }
    limit.Wait()

    h.bar.Finish()
}
