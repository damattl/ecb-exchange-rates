package main

import "encoding/xml"

type ECBRatesEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Gesmes  string   `xml:"gesmes,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Subject string   `xml:"subject"`
	Sender  struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"Sender"`
	ExchangeRates struct {
		Text                 string `xml:",chardata"`
		ExchangeRatesForTime struct {
			Text             string `xml:",chardata"`
			Time             string `xml:"time,attr"`
			ExchangeRateInfo []struct {
				Text     string `xml:",chardata"`
				Currency string `xml:"currency,attr"`
				Rate     string `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}
