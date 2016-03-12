package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type StockDetailsSent struct {
	Budget     float64
	StockNames []string
	Percentage []float64
}

type StockDetailsReceived struct {
	StockValue     []float64
	StockNames     []string
	UnvestedAmount float64
	TradeID        int
	StockQuantity  []int
}
type StockDetailsStore struct {
	StockValue     []float64
	StockNames     []string
	UnvestedAmount float64
	TradeID        int
	StockQuantity  []int
}

type StockDetailsView struct {
	StockValue         []float64
	OStockValue        []float64
	StockNames         []string
	StockQuantity      []int
	CurrentMarketValue float64
	UnvestedAmount     float64
	TradeID            int
}

type Ask struct {
	Names string `json:"Ask"`
}

type Quote struct {
	Ask Ask `json:"quote"`
}
type Results struct {
	Quote   Quote  `json:"results"`
	Count   int    `json:"count"`
	Created string `json:"created"`
	Lang    string `json:"lang"`
}
type Query struct {
	Results Results `json:"query"`
}

type Arith int

var i int

//var j int
var SS StockDetailsStore
var StockStore = make(map[int]StockDetailsStore)

func startServer() {
	arith := new(Arith)

	server := rpc.NewServer()
	server.Register(arith)

	l, e := net.Listen("tcp", ":8223")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))

	}
}
func (t *Arith) Mul(sds *StockDetailsSent, sdr *StockDetailsReceived) error {

	var UrlP1 string = "http://query.yahooapis.com/v1/public/yql?q=select%20Ask%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"
	var UrlP2 string = "%22)%0A%09%09&env=http%3A%2F%2Fdatatables.org%2Falltables.env&format=json"
	var Len int = len(sds.StockNames)

	sdr.StockNames = make([]string, Len)
	sdr.StockValue = make([]float64, Len)
	sdr.StockQuantity = make([]int, Len)
	sdr.UnvestedAmount = 0
	sdr.TradeID = rand.Intn(100000)

	for i = 0; i < Len; {

		sdr.StockNames[i] = sds.StockNames[i]
		URL := UrlP1 + sdr.StockNames[i] + UrlP2
		resp, err := http.Get(URL)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		var Queries Query

		body, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &Queries); err != nil {
			panic(err)
		}

		sdr.StockValue[i], _ = strconv.ParseFloat(Queries.Results.Quote.Ask.Names, 64)
		if sdr.StockValue[i] == 0 {
			sdr.StockQuantity[i] = 0
		}
		StockBudget := sds.Budget * sds.Percentage[i] / 100
		sdr.StockQuantity[i] = int(StockBudget / sdr.StockValue[i])
		sdr.UnvestedAmount = sdr.UnvestedAmount + StockBudget - (sdr.StockValue[i] * float64(sdr.StockQuantity[i]))

		i = i + 1

	}
	SS.StockValue = sdr.StockValue
	SS.StockNames = sdr.StockNames
	SS.UnvestedAmount = sdr.UnvestedAmount
	SS.StockQuantity = sdr.StockQuantity
	SS.TradeID = sdr.TradeID

	StockStore[SS.TradeID] = SS

	return nil

}

func (t *Arith) SQ(TradeID int, sdv *StockDetailsView) error {

	var UrlP1 string = "http://query.yahooapis.com/v1/public/yql?q=select%20Ask%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"
	var UrlP2 string = "%22)%0A%09%09&env=http%3A%2F%2Fdatatables.org%2Falltables.env&format=json"
	OS := StockStore[TradeID]
	var Len int = len(OS.StockNames)

	sdv.StockNames = make([]string, Len)
	sdv.StockValue = make([]float64, Len)
	sdv.OStockValue = make([]float64, Len)
	sdv.StockQuantity = make([]int, Len)
	sdv.CurrentMarketValue = 0
	sdv.UnvestedAmount = OS.UnvestedAmount
	sdv.TradeID = TradeID

	for i = 0; i < Len; {

		sdv.StockNames[i] = OS.StockNames[i]
		sdv.StockQuantity[i] = OS.StockQuantity[i]
		sdv.OStockValue[i] = OS.StockValue[i]
		URL := UrlP1 + sdv.StockNames[i] + UrlP2
		resp, err := http.Get(URL)
		if err != nil {
		}
		defer resp.Body.Close()
		var Queriesv Query

		body, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &Queriesv); err != nil {
			panic(err)
		}

		sdv.StockValue[i], _ = strconv.ParseFloat(Queriesv.Results.Quote.Ask.Names, 64)
		sdv.CurrentMarketValue = sdv.CurrentMarketValue + (sdv.StockValue[i] * float64(sdv.StockQuantity[i]))
		i = i + 1
	}
	return nil

}

func main() {

	// starting server in go routine (it ends on end
	// of main function
	startServer()
	var input string
	fmt.Scanln(&input)
}
