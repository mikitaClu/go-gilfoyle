package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	pair := flag.String("pair", "XRPUSDT", "Crypto pair to watch")
	priceOffset := flag.Float64("offset", 0.0002, "Diff of the currency, that should be notified")
	flag.Parse()

	pw := NewPriceWatcher(*pair,  *priceOffset)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	setInfiniteLoop(pw, c)
}

func setInfiniteLoop(pw *PriceWatcher, quit <-chan os.Signal ) {
	pw.getPrice()
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <- ticker.C:
			pw.getPrice()
		case <- quit:
			ticker.Stop()
			os.Exit(1)
		}
	}
}


