package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StockInfo struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

func FetchTWSEPrice(code string) (float64, error) {
	url := fmt.Sprintf("https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_%s.tw", code)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, err
	}

	if data, ok := parsed["msgArray"].([]interface{}); ok && len(data) > 0 {
		info := data[0].(map[string]interface{})
		priceStr := info["z"].(string)
		var price float64
		fmt.Sscanf(priceStr, "%f", &price)
		return price, nil
	}

	return 0, fmt.Errorf("找不到股價資料")
}
