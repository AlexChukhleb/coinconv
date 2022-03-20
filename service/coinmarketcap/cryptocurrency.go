package coinmarketcap

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type cryptoCurrency struct {
	ID int `json:"id"`
	//Name     string `json:"name"`
	Symbol string `json:"symbol"`
	//Slug     string `json:"slug"`
	//IsActive int    `json:"is_active"`
	//Status   string `json:"status"`
	//Rank     int    `json:"rank,omitempty"`
}

type coinMarketCap struct {
	Data struct {
		//ExchangeMap []struct {
		//	ID       int    `json:"id"`
		//	Name     string `json:"name"`
		//	Slug     string `json:"slug"`
		//	IsActive int    `json:"is_active"`
		//	Status   string `json:"status"`
		//} `json:"exchangeMap"`
		CryptoCurrencyMap []*cryptoCurrency `json:"cryptoCurrencyMap"`
	} `json:"data"`
	Status struct {
		//Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		//Elapsed      string    `json:"elapsed"`
		//CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}

type priceConversionV3 struct {
	Data struct {
		//Symbol      string  `json:"symbol"`
		//ID          string  `json:"id"`
		//Name        string  `json:"name"`
		//Amount      float64 `json:"amount"`
		//LastUpdated int64   `json:"last_updated"`
		Quote []struct {
			//CryptoID    int     `json:"cryptoId"`
			//Symbol      string  `json:"symbol"`
			Price float64 `json:"price"`
			//LastUpdated int64   `json:"lastUpdated"`
		} `json:"quote"`
	} `json:"data"`
	Status struct {
		//Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		//Elapsed      string    `json:"elapsed"`
		//CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}

func GetCoinMarketCap() (map[string]int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coinmarketcap.com/data-api/v3/map/all?start=1", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("coinMarketCap invalid status")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cmc := coinMarketCap{}
	err = json.Unmarshal(respBody, &cmc)
	if err != nil {
		return nil, err
	}

	m := map[string]int{}
	for key, val := range fiatMap {
		m[key] = val
	}

	for _, cur := range cmc.Data.CryptoCurrencyMap {
		m[cur.Symbol] = cur.ID
	}

	return m, nil
}

func PriceConversion(val string, curr1 int, curr2 int) (float64, error) {
	client := &http.Client{}

	//req, err := http.NewRequest("GET", "https://sandbox-api.coinmarketcap.com/v1/tools/price-conversion", nil) // server response does not match description
	//req, err := http.NewRequest("GET", "https://sandbox-api.coinmarketcap.com/v2/tools/price-conversion", nil) // server response does not match description
	req, err := http.NewRequest("GET", "https://api.coinmarketcap.com/data-api/v3/tools/price-conversion", nil)
	if err != nil {
		return 0, err
	}

	q := url.Values{}
	q.Add("amount", val)
	q.Add("convert_id", strconv.Itoa(curr2))
	q.Add("id", strconv.Itoa(curr1))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("coinMarketCap invalid status")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	pc := priceConversionV3{}
	err = json.Unmarshal(respBody, &pc)
	if err != nil {
		return 0, err
	}

	for _, q := range pc.Data.Quote {
		return q.Price, nil
	}

	return 0, errors.New("invalid price")
}
