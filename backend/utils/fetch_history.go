// utils/fetch_history.go
package utils

import (
	"encoding/json"
	"fmt"
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
