package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func fetchOffer() {
    _, err := http.Get("http://localhost:8090/offers")
    if err != nil {
        fmt.Println("Error:", err)
    }
}

func BenchmarkConcurrentFetchOffer(b *testing.B) {
    for {
        time.Sleep(time.Millisecond * 100)
        fetchOffer()
    }
}

func Single(b *testing.B) {
    time.Sleep(time.Second)
    fetchOffer()
    b.StopTimer()
}

func TestMain(b *testing.M) {
    benc := int64(0)
    lim := int64(100)
    for i := int64(0); i < lim; i += 1 {
        go BenchmarkConcurrentFetchOffer(&testing.B{})
    }

    for i := int64(0); i < lim; i += 1 {
        bs := &testing.B{}
        Single(bs)
        benc += bs.Elapsed().Microseconds() - time.Second.Microseconds()
    }
    benc /= lim
    fmt.Printf("benc: %v\n", benc)
}
