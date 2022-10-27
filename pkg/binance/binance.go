package binance

type Request struct {
	ProMerchantAds bool     `json:"proMerchantAds"`
	Page           int      `json:"page"`
	Rows           int      `json:"rows"`
	PayTypes       []string `json:"payTypes"`
	Asset          string   `json:"asset"`
	Fiat           string   `json:"fiat"`
	TradeType      string   `json:"tradeType"`
}

type Response struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []struct {
		Adv struct {
			AdvNo                string      `json:"advNo"`
			Classify             string      `json:"classify"`
			TradeType            string      `json:"tradeType"`
			Asset                string      `json:"asset"`
			FiatUnit             string      `json:"fiatUnit"`
			AdvStatus            interface{} `json:"advStatus"`
			PriceType            interface{} `json:"priceType"`
			PriceFloatingRatio   interface{} `json:"priceFloatingRatio"`
			RateFloatingRatio    interface{} `json:"rateFloatingRatio"`
			CurrencyRate         interface{} `json:"currencyRate"`
			Price                string      `json:"price"`
			InitAmount           interface{} `json:"initAmount"`
			SurplusAmount        string      `json:"surplusAmount"`
			AmountAfterEditing   interface{} `json:"amountAfterEditing"`
			MaxSingleTransAmount string      `json:"maxSingleTransAmount"`
			MinSingleTransAmount string      `json:"minSingleTransAmount"`
		} `json:"adv"`
		Advertiser struct {
			UserNo   string `json:"userNo"`
			NickName string `json:"nickName"`
		} `json:"advertiser"`
	} `json:"data"`
	Total   int  `json:"total"`
	Success bool `json:"success"`
}
