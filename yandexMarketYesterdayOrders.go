package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)


type RequestBody struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",attr"`
}

type YandexMarketClient struct {
	http.Client
	url			string
	authHeader	string
}

func newYandexMarketClient(campaignId, authHeader string) *YandexMarketClient {
	url := fmt.Sprintf("https://api.partner.market.yandex.ru/v2/campaigns/%s/stats/orders.xml", campaignId)
	return &YandexMarketClient{http.Client{}, url, authHeader}
}

func createRequestBody () (data []byte, err error) {
	updateTo := time.Now()
	updateFrom := updateTo.Add(-24 * time.Hour)

	requestBody := RequestBody{
		XMLName: xml.Name{Local: "order-report-request"},
		Attributes: []xml.Attr{
			{Name: xml.Name{Local: "update-from"}, Value: updateFrom.Format("2006-01-02")},
			{Name: xml.Name{Local: "update-to"}, Value: updateTo.Format("2006-01-02")},
		},
	}
	return xml.Marshal(requestBody)
}

func (y *YandexMarketClient) DownloadYesterdayOrders () error {
	requestBody, err := createRequestBody()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", y.url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", y.authHeader)
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	resp, err := y.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create("orders.xml")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func main () {
	y := newYandexMarketClient("", ``)
	if err := y.DownloadYesterdayOrders(); err != nil {
		log.Fatal(err)
	}
}