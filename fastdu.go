package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
    "sync/atomic"
    "syscall"

    "github.com/dustin/go-humanize"
    flag "github.com/spf13/pflag"
)

const VERSION = "1.0.0"

func formatSize(size int64, raw bool) string {
    if raw {
        return fmt.Sprintf("%d", size)
    }
    return strings.ReplaceAll(humanize.Bytes(uint64(size)), " ", "")
}

func main() {
    optSummary := flag.BoolP("summary", "s", false, "Show summary only")
    optRawUnit := flag.BoolP("raw", "r", false, "Show raw bytes")
    optWorker := flag.Int64P("worker", "w", 4, "Number of worker")
    optVersion := flag.BoolP("version", "v", false, "Show version")
    flag.Parse()

    if *optVersion {
        fmt.Printf("fastdu %s\n", VERSION)
        return
    }

    args := flag.Args()
    if len(args) < 1 {
        flag.PrintDefaults()
        os.Exit(1)
    }

    var wg sync.WaitGroup
    var totalSize int64
    fileChan := make(chan string)
    for i := 0; i < int(*optWorker); i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range fileChan {
                info, err := os.Stat(path)
                if err != nil {
                    continue
                } else {
                    statT := info.Sys().(*syscall.Stat_t)
                    size := statT.Blocks * 512
                    atomic.AddInt64(&totalSize, size)
                    if !*optSummary {
                        fmt.Printf("%s\t%s\n", formatSize(size, *optRawUnit), path)
                    }
                }
            }
        }()
    }

    file, err := os.Open(args[0])
    if err != nil {
        log.Fatalln(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fileChan <- scanner.Text()
    }

    close(fileChan)
    wg.Wait()
    if !*optRawUnit {
        readableSize := humanize.Bytes(uint64(totalSize))
        sanitizedSize := strings.Replace(readableSize, " ", "", -1)
        fmt.Printf("%s\ttotal\n", sanitizedSize)
    } else {
        fmt.Printf("%d\ttotal\n", totalSize)
    }
}
