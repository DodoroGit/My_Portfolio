// utils/fetch_history.go
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DailyPrice struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

type StockInfo struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

func FetchTWSEHistory(symbol string) ([]DailyPrice, error) {
	type response struct {
		Data [][]string `json:"data"`
	}

	now := time.Now()
	currentMonth := now.Format("200601")
	prevMonth := now.AddDate(0, -1, 0).Format("200601")

	dates := []string{prevMonth + "01", currentMonth + "01"}
	var allData []DailyPrice

	for _, date := range dates {
		url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/STOCK_DAY?response=json&date=%s&stockNo=%s", date, symbol)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var res response
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}

		for _, row := range res.Data {
			if len(row) >= 7 {
				// 日期格式為 "113/05/01"，轉成 "05/01"
				parts := strings.Split(row[0], "/")
				if len(parts) == 3 {
					date := fmt.Sprintf("%s/%s", parts[1], parts[2])
					priceStr := strings.ReplaceAll(row[6], ",", "")
					price, err := strconv.ParseFloat(priceStr, 64)
					if err == nil {
						allData = append(allData, DailyPrice{Date: date, Price: price})
					}
				}
			}
		}
	}

	// 按時間排序並取最新30筆
	sort.Slice(allData, func(i, j int) bool {
		t1, _ := time.Parse("01/02", allData[i].Date)
		t2, _ := time.Parse("01/02", allData[j].Date)
		return t1.Before(t2)
	})

	if len(allData) > 30 {
		allData = allData[len(allData)-30:]
	}

	return allData, nil
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

func TaiwanFee(amount float64) float64 {
	raw := amount * 0.001425 * 0.35
	rounded := math.Floor(raw/5) * 5
	if rounded < 1 {
		return 1
	}
	return rounded
}
