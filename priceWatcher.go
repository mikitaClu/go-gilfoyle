package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

const changeLabel = `
/***
 *    ______ ______  _____  _____  _____   _____  _   _   ___   _   _  _____  _____ 
 *    | ___ \| ___ \|_   _|/  __ \|  ___| /  __ \| | | | / _ \ | \ | ||  __ \|  ___|
 *    | |_/ /| |_/ /  | |  | /  \/| |__   | /  \/| |_| |/ /_\ \|  \| || |  \/| |__  
 *    |  __/ |    /   | |  | |    |  __|  | |    |  _  ||  _  || .  || | __ |  __|
 *    | |    | |\ \  _| |_ | \__/\| |___  | \__/\| | | || | | || |\  || |_\ \| |___
 *    \_|    \_| \_| \___/  \____/\____/   \____/\_| |_/\_| |_/\_| \_/ \____/\____/
 *
 *
 */
`


type priceResponse struct {
	AskPrice string `json:"askPrice"`
	AskQty string `json:"askQty"`
	BidPrice string `json:"bidPrice"`
	BidQty string `json:"bidQty"`
	CloseTime time.Time `json:"closeTime"`
	HighPrice string `json:"highPrice"`
	FirstId int `json:"firstId"`
	LastId int `json:"lastId"`
	LastPrice string `json:"lastPrice"`
	LastQty string `json:"lastQty"`
	LowPrice string `json:"lowPrice"`
	OpenPrice string `json:"openPrice"`
	OpenTime time.Time `json:"openTime"`
	PrevClosePrice string `json:"prevClosePrice"`
	PriceChange string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	QuoteVolume string `json:"quoteVolume"`
	Symbol string `json:"symbol"`
	Volume string `json:"volume"`
	WeightedAvgPrice string `json:"weightedAvgPrice"`
}


type PriceWatcher struct {
	API string
	PriceSnapshot priceResponse
	SelectedPairs string
	PriceOffsetNotification float64
	p *Player
}

func (pw* PriceWatcher) getPrice() {
	resp, err := http.Get(fmt.Sprintf("%s/ticker/24hr?symbol=%s", pw.API, pw.SelectedPairs))
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
	}
	defer resp.Body.Close()

	prevValue := pw.PriceSnapshot
	prevPrice, _ := strconv.ParseFloat(prevValue.OpenPrice, 64)
	err = json.NewDecoder(resp.Body).Decode(&pw.PriceSnapshot)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
	}
	currentPrice, _ := strconv.ParseFloat(pw.PriceSnapshot.OpenPrice, 64)
	curr := math.Abs(prevPrice - currentPrice)
	fmt.Printf("Previous Price is: %s, Current Price is: %s \n", prevValue.OpenPrice, pw.PriceSnapshot.OpenPrice)
	if (curr >= pw.PriceOffsetNotification) && prevPrice != 0 {
		fmt.Printf("%s\n", changeLabel)
		pw.p.PlayAudio()
	}
}



func NewPriceWatcher(pairs string, priceOffset float64) *PriceWatcher {
	return &PriceWatcher{
		API: "https://api.binance.com/api/v3",
		SelectedPairs: pairs,
		PriceOffsetNotification: priceOffset,
		PriceSnapshot: priceResponse{},
		p: NewPlayer(),
	}
}
