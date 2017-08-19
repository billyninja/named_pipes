package main

import (
    "syscall"
    "encoding/gob"
    "bytes"
    "log"
    "fmt"
    "io/ioutil"
    "time"
    "github.com/billyninja/named_pipes/stats"
    "github.com/billyninja/named_pipes/client"
)

var kn []*stats.Node

func read() {
    var st2 stats.Status

    rChunk, err := ioutil.ReadFile("/tmp/stats_pipe")

    if err != nil {
        log.Fatal(err)
    }

    rbuff := bytes.NewBuffer(rChunk)

    dec := gob.NewDecoder(rbuff)
    err = dec.Decode(&st2)
    if err != nil {
        fmt.Printf("Err: %v", err)
    }

    fmt.Printf("x:%v\n\ny:%v\n\n", rChunk, rbuff.Bytes())
    fmt.Printf("2 - Read %d \n%v\n", len(rChunk), st2.X200)

    st2 = stats.Status{}
}

func main() {
    syscall.Mkfifo("/tmp/stats_pipe", 0666)
    cls := []*client.Client{}
    for i := 1; i <= 10; i++ {
        cls = append(cls, client.NewClient(i))
    }
    fmt.Printf("Spawned %d clients", len(cls))

    go func () {
        for {
            read()
            time.Sleep(50 * time.Millisecond)
        }
    }()

    go func () {
        for {
            for _, c := range cls {
                c.Report()
                time.Sleep(100 * time.Millisecond)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }()

    for {
        time.Sleep(100 * time.Millisecond)
    }
}
