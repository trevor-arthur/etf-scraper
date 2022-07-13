package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type EtfInfo struct {
	Title              string
	Replication        string
	Earnings           string
	TotalExpenseRatio  string
	TrackingDifference string
	FundSize           string
}

func cleanDesc(s string) string {
	// A function that cleans the descriptions of any whitespaces
	return strings.TrimSpace(s)
}

func scrapeUrl(isin string) string {
	// A function which returns the base website path with each given ISIN
	return "https://www.trackingdifferences.com/ETF/ISIN/" + isin
}

func main() {

	// Different ETF ISIN's for the URL
	isins := []string{
		"IE00B1XNHC34",
		"IE00B4L5Y983",
		"LU1838002480",
	}

	// Create empty EtfInfo instance
	etfInfo := EtfInfo{}

	// EtfInfo Slice will hold all the different EtfInfos from the scraping
	etfInfos := make([]EtfInfo, 0, 1)

	// Created a New Collector and allows only domains within scope
	c := colly.NewCollector(colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com"))

	c.OnRequest(func(r *colly.Request) {
		// Sets header to English and prints the URL that is getting scraped
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("[*] Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		// Error handling
		fmt.Printf("[!] Error: %s\n", e.Error())
	})

	c.OnHTML("h1.page-title", func(h *colly.HTMLElement) {
		// Search within the <h1 class="page-title"> for etfInfo's Title text
		etfInfo.Title = h.Text
	})

	c.OnHTML("div.descfloat p.desc", func(h *colly.HTMLElement) {
		// Search within the <div class="descfloat"> and then within <p class="desc">

		selection := h.DOM

		childNodes := selection.Children().Nodes
		if len(childNodes) == 3 {
			// Checks to make sure there are only 3 HTML lines (nodes)
			//	Node 1:	<span class="desctitle desctitle-high">Replication</span>
			//	Node 2:	<br>
			//	Node 3:	<span class="catvalue-default phys">Physical</span>

			// Search within the <span class="desctitle">
			description := cleanDesc(selection.Find("span.desctitle").Text())

			// Sets the value as the 3rd Node (the 3rd HTML line)
			value := selection.FindNodes(childNodes[2]).Text()

			switch description {
			// A switch for each of the description fields
			case "Replication":
				etfInfo.Replication = value
				break
			case "TER":
				etfInfo.TotalExpenseRatio = value
				break
			case "TD":
				etfInfo.TrackingDifference = value
				break
			case "Earnings":
				etfInfo.Earnings = value
				break
			case "Fund size":
				etfInfo.FundSize = value
				break
			}

		}

	})

	c.OnScraped(func(r *colly.Response) {
		etfInfos = append(etfInfos, etfInfo)
		etfInfo = EtfInfo{}
	})

	for _, isin := range isins {
		// Visits each ISIN URL in the ISINs slice
		// https://www.trackingdifferences.com/ETF/ISIN/ + isins[0]
		// https://www.trackingdifferences.com/ETF/ISIN/ + isins[1]
		// https://www.trackingdifferences.com/ETF/ISIN/ + isins[2]
		c.Visit(scrapeUrl(isin))
	}

	// Create json out of response
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(etfInfos)

}
