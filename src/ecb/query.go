package ecb

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v\n", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v\n", err)
	}

	return data, nil
}

func GetTodaysRates() (*ECBRatesEnvelope, error) {
	if xmlBytes, err := getXML(ecbURL); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get XML: %v\n", err))
	} else {
		var result ECBRatesEnvelope
		if err = xml.Unmarshal(xmlBytes, &result); err != nil {
			return nil, errors.New(fmt.Sprintf("Could not parse data: %v\n", err))
		}
		return &result, nil
	}
}
