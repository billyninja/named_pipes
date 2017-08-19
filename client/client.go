package client

import (
    "os"
    "bytes"
    "encoding/gob"
    "io/ioutil"
    "log"
    "time"
    "github.com/billyninja/named_pipes/stats"
)

type Client struct {
    id      int
    enc     *gob.Encoder
    Status  *stats.Status
    wbuff   *bytes.Buffer
}

func NewClient(id int) *Client {
    return &Client{
        id:     id,
        Status: &stats.Status{
            X200: &stats.ThroughPut{},
            X400: &stats.ThroughPut{},
            X500: &stats.ThroughPut{},
        },
    }
}

func (c *Client) Report() {
    t1 := time.Now()

    wbuff := &bytes.Buffer{}
    enc := gob.NewEncoder(wbuff)
    enc.Encode(c.Status)

    err := ioutil.WriteFile("/tmp/stats_pipe", wbuff.Bytes(), os.ModeNamedPipe)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("%d - written %d \n", c.id, len(wbuff.Bytes()))
    lat := time.Since(t1)
    c.Status.X200.Len     += int64(len(wbuff.Bytes()))
    c.Status.X200.Count   += 1
    c.Status.X200.AvgLat  += lat
    c.Status.X200.AvgLat  /= 2
}
