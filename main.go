package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	selectedCryptoCurrency = "Null"
	selectedCryptoCurrencyIndex int
	selectedCurrency = "Null"
	balance float64 = 0 //0 = nil | your crypto amount
	previousPrice = 0.0
	currentText string
	newChange string
	console *tview.TextView
	app *tview.Application
	)


type requestData struct {
	Data []map[string]interface{} `json:"data"`
}

func consoleCleaner(console *tview.TextView, app *tview.Application){
	console.SetText("Starting...")
}

func consoleManager(text string, console *tview.TextView, app *tview.Application) {
	newText := fmt.Sprint(currentText+"[-:-:-]\n"+text)
	currentText = newText
	console.SetText(newText).SetTextAlign(tview.AlignLeft)
}

func refreshData(cc string, cci int,sc string, bal float64, console *tview.TextView){
	if cc == "Null" || sc == "Null" || bal == 0 {
		console.SetText("Please Select a cryptocurrency, a currency and set your balance")
	} else {
		consoleCleaner(console,app)
		currentText = ""
		previousPrice = 0.0

		go func() {
			for {
				if cc != selectedCryptoCurrency || bal != balance || sc != selectedCurrency{
					break
				}

				infoRequest(cci,sc,bal,console)
				time.Sleep(1 * time.Minute)
			}
		}()
	}
}

func infoRequest(cci int,sc string, bal float64, console *tview.TextView){
	var message = ""
	currentTime := time.Now().Format("[15:04:05]")
	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
		"Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.2; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0; Trident/5.0)",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729)"}

	randomIndex := rand.Intn(len(userAgentList))

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://www.coinbase.com/api/v2/assets/prices?base="+sc+"&filter=listed&resolution=latest", nil)
	req.Header.Add("User-Agent", userAgentList[randomIndex])
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "None")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")

	resp, err := client.Do(req)

	if err != nil {
		console.SetText("Error when sending request to the server")
	}

	defer resp.Body.Close()

	var bodyInBytes []byte
	if resp.Body != nil {
		bodyInBytes, _ = ioutil.ReadAll(resp.Body)
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyInBytes))

	jsonData := requestData{}
	err = json.Unmarshal(bodyInBytes, &jsonData)
	if err != nil {
		console.SetText(err.Error())
	}

	cryptoName := jsonData.Data[cci]["base"].(string)
	currency := jsonData.Data[cci]["currency"].(string)
	price := jsonData.Data[cci]["prices"].(map[string]interface{})["latest"].(string)
	change := jsonData.Data[cci]["prices"].(map[string]interface{})["latest_price"].(map[string]interface{})["percent_change"].(map[string]interface{})["hour"].(float64)
	floatPrice, _ := strconv.ParseFloat(price, 64)
	currentBalance := (floatPrice * bal) / 1

	//Market Status
	if change >= 0.0 {
		newChange = fmt.Sprintln(" [green:]Market is up by " + fmt.Sprintf("%g",change) + "%")
	} else if change < 0.0 {
		newChange = fmt.Sprintln(" [red:]Market down by " + fmt.Sprintf("%g",change) + "%")
	}

	//Price Status
	if previousPrice == floatPrice {
		message = fmt.Sprint("[gray:] "+currentTime+" [orange:] [*] same price "+fmt.Sprintf("%g", floatPrice))
	} else if previousPrice < floatPrice {
		upMessage := fmt.Sprint("[green:] [+] " + cryptoName + " price is: " + fmt.Sprintf("%g", floatPrice) + " " + currency)
		message = fmt.Sprint("[gray:] "+currentTime+upMessage+newChange+"[green:] since the last hour. "+"[yellow:]Balance: "+fmt.Sprintf("%.2f", currentBalance)+" "+currency)
		previousPrice = floatPrice
	} else if previousPrice > floatPrice {
		downMessage := fmt.Sprint("[red:] [-] " + cryptoName + " price is: " + fmt.Sprintf("%g", floatPrice) + " " + currency)
		message = fmt.Sprint("[gray:] "+currentTime+downMessage+newChange+"[red:] since the last hour. "+"[yellow:]Balance: "+fmt.Sprintf("%.2f", currentBalance)+" "+currency)
		previousPrice = floatPrice
	} else {}
	consoleManager(message,console,app)
}

func cryptoCurrencyList(app *tview.Application) *tview.List{
	list := tview.NewList()
	list.AddItem("ETH", "Ethereum", 0, func() {
		selectedCryptoCurrency = "ETH"
		selectedCryptoCurrencyIndex = 3
		refreshData("ETH", 3, selectedCurrency, balance, console)
	})
	list.AddItem("XRP", "Ripple", 0, func() {
		selectedCryptoCurrency = "XRP"
		selectedCryptoCurrencyIndex = 18
		refreshData("XRP", 18, selectedCurrency, balance, console)
	})
	list.AddItem("LTC", "Litecoin", 0, func() {
		selectedCryptoCurrency = "LTC"
		selectedCryptoCurrencyIndex = 5
		refreshData("LTC", 5, selectedCurrency, balance, console)
	})
	list.AddItem("BCH", "Bitcoin Cash", 0, func() {
		selectedCryptoCurrency = "BCH"
		selectedCryptoCurrencyIndex = 1
		refreshData("BCH", 1, selectedCurrency, balance, console)
	})
	list.AddItem("EOS", "EOS", 0, func() {
		selectedCryptoCurrency = "EOS"
		selectedCryptoCurrencyIndex = 21
		refreshData("EOS", 21, selectedCurrency, balance, console)
	})
	list.AddItem("BSV", "Bitcoin SV", 0, func() {
		selectedCryptoCurrency = "BSV"
		selectedCryptoCurrencyIndex = 2
		refreshData("BSV", 2, selectedCurrency, balance, console)
	})
	list.AddItem("EXIT", "", 'q', func() {
		app.Stop()
	})
	return list
}

func currencyList() *tview.List{
	list := tview.NewList()
	list.AddItem("EUR","",0, func() {
		selectedCurrency = "EUR"
		refreshData(selectedCryptoCurrency, selectedCryptoCurrencyIndex, selectedCurrency, balance, console)
	})
	list.AddItem("USD", "", 0, func() {
		selectedCurrency = "USD"
		refreshData(selectedCryptoCurrency, selectedCryptoCurrencyIndex, selectedCurrency, balance, console)
	})
	list.AddItem("GBP", "", 0, func() {
		selectedCurrency = "GBP"
		refreshData(selectedCryptoCurrency, selectedCryptoCurrencyIndex, selectedCurrency, balance, console)
	})
	return list
}

func main() {
	app = tview.NewApplication()

	console = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Please Select a cryptocurrency, a currency and set your balance").SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})

	listGrid := tview.NewGrid().
		SetRows(2, 0).
		SetBorders(true).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Crypto Currency"), 0,0,1,1,0,0,false).
		AddItem(cryptoCurrencyList(app), 1,0,1,1,0, 0,false)


	balanceField := tview.NewInputField().
		SetLabel("Enter your balance: ").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldFloat).
		SetChangedFunc(func(text string) {
		balance, _ = strconv.ParseFloat(text, 64)
	}).SetDoneFunc(func(key tcell.Key) {
		refreshData(selectedCryptoCurrency, selectedCryptoCurrencyIndex, selectedCurrency, balance, console)
	})

	rightMenu := tview.NewGrid().
		SetRows(2,0,3).
		SetBorders(true).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Currency"), 0,0,1,1,0,0,false).
		AddItem(currencyList(), 1,0,1,1,0,0,false).
		AddItem(balanceField, 2,0,1,1,0,0,false)

	grid := tview.NewGrid().
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(listGrid, 0, 0, 1, 1, 0, 0, false).
		AddItem(console, 0,1,1,1,0,0,true).
		AddItem(rightMenu,0,2,1,1,0,0,false)

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
