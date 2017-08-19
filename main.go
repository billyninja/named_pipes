package main

import (
    "encoding/gob"
    "bytes"
    "log"
    "os"
    "syscall"
    "fmt"
    "io/ioutil"
    "time"
)

type ThroughPut struct {
    Count   int64
    Len     int64
    AvgLat  time.Duration
    MinLat  time.Duration
    MaxLat  time.Duration
}

func (tp ThroughPut) String() string {
    return fmt.Sprintf("Count:\n%v\nLen:\n%v\nAvgLat\n%v\n", tp.Count, tp.Len, tp.AvgLat)
}

type Status struct {
    X200 *ThroughPut
    X400 *ThroughPut
    X500 *ThroughPut
}

func (s Status) String() string {
    return fmt.Sprintf("2XX:\n%v\n4XX:\n%v\n5XX\n%v\n", s.X200, s.X400, s.X500)
}


func write(st Status) {
    t1 := time.Now()
    var wbuff bytes.Buffer
    enc := gob.NewEncoder(&wbuff)
    enc.Encode(st)
    err := ioutil.WriteFile("/tmp/stats_pipe", wbuff.Bytes(), os.ModeNamedPipe)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("1 - written %d \n", len(wbuff.Bytes()))
    lat := time.Since(t1)

    st.X200.Count += 1
    st.X200.Len += int64(len(wbuff.Bytes()))
    st.X200.AvgLat += lat
    st.X200.AvgLat /= 2
}

func read() {
    var st2 Status

    rChunk, err := ioutil.ReadFile("/tmp/stats_pipe")
    if err != nil {
        log.Fatal(err)
    }

    rbuff := bytes.NewBuffer(rChunk)
    dec := gob.NewDecoder(rbuff)
    dec.Decode(&st2)

    fmt.Printf("2 - Read %v\n", st2)
}

func main() {
    syscall.Mkfifo("/tmp/stats_pipe", 0666)

    st := Status{
        X200: &ThroughPut{},
        X400: &ThroughPut{},
        X500: &ThroughPut{},
    }

    for {

        go func() {
            read()
        }()

        go func() {
            write(st)
        }()

        go func() {
            read()
        }()

        go func() {
            write(st)
        }()

        time.Sleep(1 * time.Second)
    }
}
