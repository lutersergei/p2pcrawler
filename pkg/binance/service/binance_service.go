package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"p2p_crawler/pkg/binance"
	"p2p_crawler/pkg/price"
	"strconv"
)

type BinanceService struct {
	client *http.Client
}

func NewBinanceService(client *http.Client) *BinanceService {
	return &BinanceService{client: client}
}

func (svc *BinanceService) DoRequest() (*price.PriceModel, error) {
	resp, err := svc.client.Do(makeRequest())
	if err != nil {
		return nil, fmt.Errorf("binance request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	data := &binance.Response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	return svc.handleResponse(data)
}

func (svc *BinanceService) handleResponse(resp *binance.Response) (*price.PriceModel, error) {
	if len(resp.Data) == 0 {
		return nil, nil
	}

	bestResult := resp.Data[0]

	priceFloat, err := strconv.ParseFloat(bestResult.Adv.Price, 64)
	if err != nil {
		return nil, err
	}

	surplusAmount, err := strconv.ParseFloat(bestResult.Adv.SurplusAmount, 64)
	if err != nil {
		return nil, err
	}

	rawJSON, _ := json.Marshal(bestResult)

	model := &price.PriceModel{
		BestPrice:     priceFloat,
		Username:      bestResult.Advertiser.NickName,
		RawJSON:       string(rawJSON),
		SurplusAmount: surplusAmount,
		Exchange:      svc.GetName(),
	}

	return model, nil
}

func (svc *BinanceService) GetName() string {
	return "Binance"
}

func makeRequest() *http.Request {
	reqData, _ := json.Marshal(&binance.Request{
		ProMerchantAds: false,
		Page:           1,
		Rows:           10,
		PayTypes:       []string{"BAKAIBANK", "DEMIRBANK", "ELCART"},
		Asset:          "USDT",
		Fiat:           "KGS",
		TradeType:      "SELL",
	})

	reqBody := bytes.NewBuffer(reqData)

	req, _ := http.NewRequest("POST", "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", reqBody)
	addHeaders(req)

	return req
}

func addHeaders(req *http.Request) {
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("csrftoken", "d41d8cd98f00b204e9800998ecf8427e")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.167 YaBrowser/22.7.5.934 (beta) Yowser/2.5 Safari/537.36")
}
