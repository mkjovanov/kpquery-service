package service

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"kpquery-service/config"
	"kpquery-service/model"
	"kpquery-service/util"

	"github.com/PuerkitoBio/goquery"
)

func Start() {
	fmt.Print("Enter search query: ")
	//fmt.Scanln(&_query)
	_query = "zelda"

	file := createResultsFile()
	defer file.Close()

	buildRequestUrl(0)
	fetchNumberOfPages()
	for i := 1; i <= int(_numOfPages); i++ {
		buildRequestUrl(i)
		fetchAds()
	}
	for _, ad := range _ads {
		fmt.Fprintln(file, "----------------")
		fmt.Fprintln(file, ad.Name)
		fmt.Fprintln(file, ad.Price)
		fmt.Fprintln(file, ad.Url)
		fmt.Println(file, "----------------")
		fmt.Println("----------------")
		fmt.Println(ad.Name)
		fmt.Println(ad.Price)
		fmt.Println(ad.Url)
		fmt.Println("----------------")
	}
	fmt.Printf("Press any key to quit...")
	fmt.Scanln(&_query)
}

var (
	_query      string
	_ads        []model.Advertisement
	_names      []string
	_urls       []string
	_prices     []string
	_numOfPages int64
	_requestUrl string
	_mainConf   *config.Configuration
)

func formatSearchUrl(pageNum int) string {
	formatedUrl := strings.Replace(_mainConf.SearchUrl, "{KP_QUERY_PLACEHOLDER}", _query, 1)
	if pageNum != 0 {
		formatedUrl = strings.Replace(formatedUrl, "{KP_PAGE_NUMBER}", string(pageNum), 1)
	}
	return formatedUrl
}

func findMax(queryArray []int64) int64 {
	for j := 1; j < len(queryArray); j++ {
		if queryArray[0] < queryArray[j] {
			queryArray[0] = queryArray[j]
		}
	}
	return queryArray[0]
}

func createResultsFile() *os.File {
	os.Mkdir(_query, 0777)
	destination := _query + "/results.txt"
	file, err := os.Create(destination)
	if err != nil {
		panic(err)
	}
	return file
}

func fetchAds() {
	resp, err := http.Get(_requestUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	itemNames := _doc.Find(_mainConf.ItemNameNode)
	itemUrls := _doc.Find(_mainConf.ItemUrlNode)
	itemPrices := _doc.Find(_mainConf.ItemPriceNode).Contents()

	for _, nameNode := range itemNames.Contents().Nodes {
		_names = append(_names, strings.TrimSpace(nameNode.Data))
	}
	for _, urlNode := range itemUrls.Nodes {
		_urls = append(_urls, urlNode.Attr[0].Val)
	}
	for _, itemNode := range itemPrices.Nodes {
		_prices = append(_prices, strings.TrimSpace(itemNode.Data))
	}

	for index, name := range _names {
		ad := model.NewAdvertisement(
			name,
			_prices[index],
			_urls[index],
		)
		_ads = append(_ads, *ad)
	}
}

func fetchNumberOfPages() {
	resp, err := http.Get(_requestUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	numberOfPages := _doc.Find(_mainConf.NumberOfPagesNode).First().Children().Contents()
	pageNums := []int64{0, 0, 0, 0, 0, 0}
	for _, page := range numberOfPages.Nodes {

		if len(page.Attr) != 0 {
			//string is "Strana <num>" - cut "Strana " from string
			if s, err := strconv.ParseInt(page.Attr[1].Val[7:], 0, 0); err == nil {
				pageNums = append(pageNums, s)
			}
		}
	}
	_numOfPages = findMax(pageNums)
}

func buildRequestUrl(pageNum int) {
	var stringBuilder strings.Builder
	util.SanitizeQuery(_query)
	searchUrl := formatSearchUrl(pageNum)
	stringBuilder.WriteString(searchUrl)
	stringBuilder.WriteString(strconv.Itoa(pageNum))

	_requestUrl = stringBuilder.String()
}
