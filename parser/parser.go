package parser

import (
	"cryptoGo/database"
	"cryptoGo/depth"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
)

type Parser struct {
	db       *sql.DB
	Symbols  []string
	Interval string
	stop     chan struct{}
}

func (p *Parser) StartParser() {
	for _, symbol := range p.Symbols {
		go p.Parse(symbol)
	}
	fmt.Printf("Start parsing. The number of Assets is: %v.\n", len(p.Symbols))
}

func (p *Parser) ParseDepth(symbol string, ch chan depth.DepthInfo) {
	bidsList := []common.PriceLevel{}
	asksList := []common.PriceLevel{}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	wsDepthHandler := func(event *binance.WsDepthEvent) {
		bidsList = append(bidsList, event.Bids...)
		asksList = append(asksList, event.Asks...)
	}

	doneC, stopC, err := binance.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		p.stop <- struct{}{}
		return
	}

	for {
		select {
		case <-p.stop:
			fmt.Printf("Stop parsing %s\n", symbol)
			stopC <- struct{}{}
			return
		case depthInfo := <-ch:
			depthInfo.SetDepth(bidsList, asksList)

			// Reset lists for the next depth update
			bidsList = nil
			asksList = nil

			// Send depth info to the channel
			ch <- depthInfo
		case <-doneC:
			return
		}
	}
}

func (p *Parser) Parse(symbol string) {
	errHandler := func(err error) {
		fmt.Println(err)
	}
	ch := make(chan depth.DepthInfo)

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		if event.Kline.IsFinal {
			ch <- depth.DepthInfo{}
			depthInfo := <-ch

			fmt.Println(event.Kline)
			database.AddNewRow(event.Kline, depthInfo, os.Getenv("DbTableName"), p.db)
		}
	}

	go p.ParseDepth(symbol, ch)
	doneC, stopC, err := binance.WsKlineServe(symbol, p.Interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		p.stop <- struct{}{}
		return
	}

	for {
		select {
		case <-p.stop:
			fmt.Printf("Stop parsing %s\n", symbol)
			stopC <- struct{}{}
			return
		case <-doneC:
			return
		}
	}
}

func (p *Parser) StopParser() {
	for i := 0; i < len(p.Symbols)*2; i++ {
		time.Sleep(time.Millisecond * 10)
		p.stop <- struct{}{}
	}
	close(p.stop)
}

func InitParser(db *sql.DB) *Parser {
	symbolsStr := os.Getenv("BINANCE_SYMBOLS")
	if symbolsStr == "" {
		log.Fatal("BINANCE_SYMBOLS is not set")
	}

	symbols := strings.Split(symbolsStr, ",")

	return &Parser{
		db:       db,
		Symbols:  symbols,
		Interval: "1m",
		stop:     make(chan struct{}),
	}
}
