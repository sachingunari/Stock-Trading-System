package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
)

type StockDetailsSent struct {
	Budget     float64
	StockNames []string
	Percentage []float64
}

type StockDetailsReceived struct {
	StockValue     []float64
	StockNames     []string
	StockQuantity  []int
	UnvestedAmount float64
	TradeID        int
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

type Arith int

var Input string
var i int
var TradeID int
var j string
var loop bool
var total float64

func buyStockDetails(sds *StockDetailsSent) {

	fmt.Println("Enter The Budget\n")
	fmt.Scanf("%f\n", &sds.Budget)
	fmt.Println("Enter The Stock Name and Percentage in Format Stock1Name,Percentage,Stock2Name,Percentage,Stock3Name,Percentage \n")
	fmt.Scanf("%s\n", &Input)
	Lent := strings.Split(Input, ",")
	fmt.Println(Lent)
	sds.StockNames = make([]string, (len(Lent) / 2))
	sds.Percentage = make([]float64, (len(Lent) / 2))
	for i = 0; i < len(Lent); i = i + 2 {
		sds.StockNames[(i / 2)] = string(Lent[i])
		sds.Percentage[(i / 2)], _ = strconv.ParseFloat(Lent[i+1], 64)

	}
}

func displayStocks(sdr StockDetailsReceived) {
	fmt.Println("The UnvestedAmount is", sdr.UnvestedAmount)
	fmt.Println("The TradeID is: ", sdr.TradeID)
	fmt.Println("Name \t No of Stocks \t Stock Value \n")
	for i := 0; i < len(sdr.StockNames); i++ {
		fmt.Println(sdr.StockNames[i], "\t", sdr.StockQuantity[i], "\t\t", sdr.StockValue[i], "\t\t\n")
	}
}
func displayStocksv(sdv StockDetailsView) {
	fmt.Println("The UnvestedAmount is", sdv.UnvestedAmount)
	fmt.Println("The TradeID is: ", sdv.TradeID)
	fmt.Println("The CurrentMarket Value is: ", sdv.CurrentMarketValue)
	fmt.Println("Name \t No of Stocks \t Present Stock Value \n")
	for i := 0; i < len(sdv.StockNames); i++ {

		if sdv.StockValue[i] == sdv.OStockValue[i] {
			fmt.Println(sdv.StockNames[i], "\t", sdv.StockQuantity[i], "\t\t", "=", sdv.StockValue[i], "\t\t\n")
		} else if sdv.StockValue[i] < sdv.OStockValue[i] {
			fmt.Println(sdv.StockNames[i], "\t", sdv.StockQuantity[i], "\t\t", "-", sdv.StockValue[i], "\t\t\n")
		} else if sdv.StockValue[i] > sdv.OStockValue[i] {
			fmt.Println(sdv.StockNames[i], "\t", sdv.StockQuantity[i], "\t\t", "+", sdv.StockValue[i], "\t\t\n")
		}
	}
}
func startClient(sds StockDetailsSent, sdr *StockDetailsReceived) {

	conn, err := net.Dial("tcp", "localhost:8223")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := jsonrpc.NewClient(conn)

	err = c.Call("Arith.Mul", sds, &sdr)

	if err != nil {
		log.Fatal("arith error:", err)
	}

}
func QueryStocks(TradeID int, sdv *StockDetailsView) {

	fmt.Println("Enter The Trade ID to Enquire Status \n")
	fmt.Scanf("%d\n", &TradeID)

	conn, err := net.Dial("tcp", "localhost:8223")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := jsonrpc.NewClient(conn)

	err = c.Call("Arith.SQ", TradeID, &sdv)

	if err != nil {
		log.Fatal("arith error:", err)
	}

}

func main() {
	var sds StockDetailsSent
	var sdr StockDetailsReceived
	var sdv StockDetailsView

	for !loop {
		fmt.Println("Enter 'b' to buy stocks,'c' to check portfolio ,any other key to exit \n")
		fmt.Scanf("%s\n", &j)

		if j == "b" {
			buyStockDetails(&sds)
			for _, value := range sds.Percentage {
				total += value
			}
			if total == 100 {
				startClient(sds, &sdr)
				displayStocks(sdr)
			} else {
				fmt.Println("The sum of Percentages of bugdet must be 100")
			}
		} else if j == "c" {
			QueryStocks(TradeID, &sdv)
			displayStocksv(sdv)
		} else {
			loop = true
		}

	}
	var input string
	fmt.Scanln(&input)

}
