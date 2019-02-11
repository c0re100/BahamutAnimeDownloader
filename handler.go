package main

import (
    "errors"
    "fmt"
    "log"
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
    f, e := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
    if e != nil {
        fmt.Println(err.Error())
    }

    log.SetOutput(f)
    if err != nil {
        msg := msg + " " + err.Error()
        fmt.Println(msg)
        log.Println(msg)
        f.Close()
        os.Exit(1)
    }
}

func (h *bahamut) getQuality() (string, string) {
    switch h.quality {
    case "360p":
        return h.res.s360, "360p"
    case "540p":
        return h.res.s540, "540p"
    case "720p":
        return h.res.s720, "720p"
    case "1080p":
        return h.res.s1080, "1080p"
    default:
        return h.res.s720, "Default(720p)"
    }
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
