package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	_ "github.com/Xe/gopreload"
	"github.com/Xe/ln"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		spew()

		ln.Log(ln.F{"action": "gc_spew", "who": r.RemoteAddr})

		fmt.Fprintln(w, "done")
	})

	http.ListenAndServe(":9184", nil)
}

func makeBuffer() []byte {
	return make([]byte, rand.Intn(5000000)+5000000)
}

func spew() {
	pool := make([][]byte, 20)

	var m runtime.MemStats
	makes := 0
	for _ = range make([]struct{}, 50) {
		b := makeBuffer()
		makes += 1
		i := rand.Intn(len(pool))
		pool[i] = b

		time.Sleep(time.Millisecond * 250)

		bytes := 0

		for i := 0; i < len(pool); i++ {
			if pool[i] != nil {
				bytes += len(pool[i])
			}
		}

		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes)
	}
}
