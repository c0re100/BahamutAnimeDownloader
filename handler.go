package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

func envCheck() {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    isErr("Path not found", err)
    os.Setenv("PATH", os.Getenv("PATH")+";"+dir)

    err = exec.Command("ffmpeg").Run()
    if err != nil {
        if strings.Contains(err.Error(), "executable file not found") {
            isErr("ffmpeg not found -", errors.New("please download from official website first"))
        }
    }
}

func isErr(msg string, err error) {
    if err != nil {
        f, e := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
        if e != nil {
            log.Fatal(err.Error())
        }

        log.SetOutput(f)
        msg := msg + " " + err.Error()
        fmt.Println(msg)
        log.Fatal(msg)
    }
}

func (h *bahamut) getQuality() (string, string) {
    return h.res, h.quality
}

func (h *bahamut) request(action, url string) *http.Response {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        isErr("Create "+action+" request failed - ", err)
    }

    req.Header.Add("cookie", h.rawCookie)
    req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
    req.Header.Add("referer", "https://ani.gamer.com.tw/animeVideo.php?sn="+h.sn)
    req.Header.Add("origin", "https://ani.gamer.com.tw")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        isErr("Create "+action+" request failed - ", err)
    }

    return resp
}

func (h *bahamut) mergeChunk() {
    fmt.Println("Merging chunk file...please wait a moment...")
    os.Mkdir("output", 0755)
    exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-y", "-i", h.tmp+"/"+h.plName, "-c", "copy", "output/"+h.sn+".mp4").Run()
    fmt.Println("File location: output/" + h.sn + ".mp4")
}

func (h *bahamut) cleanUp() {
    // Delete a temporary directory
    os.RemoveAll(h.tmp)

    fmt.Println("Cleaned up.")
    fmt.Println(fmt.Sprintf("Total time: %ds", time.Now().Unix()-h.startTime))
}
