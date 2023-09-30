package zillow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
}

type SearchPropertiesResponse struct {
	Results []struct {
		Display    string
		ResultType string
		MetaData   struct {
			AddressType  string
			StreetNumber string
			StreetName   string
			City         string
			State        string
			Country      string
			ZipCode      string
			Zpid         int
			Lat          float64
			Lng          float64
			MaloneId     int
		}
	}
}

type PropertyLookupResponse struct {
	LookupResults []struct {
		Zpid      int
		Estimates struct {
			Zestimate     int64
			RentZestimate int64
		}
	}
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient}
}

func (c *Client) SearchProperties(address string) (*SearchPropertiesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.zillowstatic.com/autocomplete/v3/suggestions?q=%s", url.QueryEscape(address)), nil)
	if err != nil {
		return nil, fmt.Errorf("error building zillow autocomplete request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying zillow autocomplete: %w", err)
	}
	defer resp.Body.Close()
	var result *SearchPropertiesResponse

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading zillow autocomplete response: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error querying zillow autocomplete: %s", string(b))
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("error json decoding zillow autocomplete response: %w", err)
	}

	return result, nil
}

func (c *Client) LookupProperty(propertyId int) (*PropertyLookupResponse, error) {
	payload := map[string]interface{}{
		"homeDetailsUriParameters": map[string]interface{}{
			"googleMaps":           false,
			"platform":             "iphone",
			"showFactsAndFeatures": true,
			"streetView":           false,
		},
		"listingCategoryFilter": "all",
		"propertyIds":           []int{propertyId},
		"sortAscending":         true,
		"sortOrder":             "recentlyChanged",
	}
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("error encoding json for zillow api request: %w", err)
	}
	req, err := http.NewRequest("POST", "https://zm.zillow.com/api/public/v2/mobile-search/homes/lookup", payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("error building zillow api request: %w", err)
	}
	req.Header.Set("User-Agent", "ZillowMap/15.12.0.0 CFNetwork/1240.0.4 Darwin/20.5.0")
	req.Header.Set("X-Client", "com.zillow.ZillowMap")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying zillow api: %w", err)
	}
	defer resp.Body.Close()
	var result *PropertyLookupResponse

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading zillow api response: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error querying zillow api: %s", string(b))
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("error json decoding zillow api response: %w", err)
	}

	return result, nil
}
