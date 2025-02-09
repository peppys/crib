package redfin

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
}

type SearchPropertiesResponse struct {
	Version int
	Payload struct {
		Sections []struct {
			Rows []struct {
				ID      string
				Name    string
				SubName string
				URL     string
				Active  bool
			}
		}
		ExactMatch struct {
			ID      string
			Name    string
			SubName string
			URL     string
			Active  bool
		}
	}
}

type AVMResponse struct {
	Version int
	Payload struct {
		Root struct {
			AVMInfo struct {
				PredictedValue float64
			}
		} `json:"__root"`
	}
}

func NewClient(httpClient *http.Client) *Client {
	mustConfigureAuthCookies(httpClient)

	return &Client{httpClient}
}

func (c *Client) SearchProperties(address string) (*SearchPropertiesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.redfin.com/stingray/do/location-autocomplete?location=%s&v=2", url.QueryEscape(address)), nil)
	if err != nil {
		return nil, fmt.Errorf("error building redfin api request: %w", err)
	}
	req.Header.Set("User-Agent", "Redfin/556.0.0.23067 CFNetwork/1240.0.4 Darwin/20.5.0")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying redfin api: %w", err)
	}
	defer resp.Body.Close()
	var result *SearchPropertiesResponse

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading redfin api response: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error querying redfin api: %s", string(b))
	}
	err = json.Unmarshal([]byte(strings.TrimPrefix(string(b), "{}&&")), &result)
	if err != nil {
		return nil, fmt.Errorf("error json decoding redfin api response: %w", err)
	}

	return result, nil
}

func (c *Client) GetAutomatedValuationModel(propertyId string) (*AVMResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.redfin.com/stingray/mobile/api/home/details/avm?propertyId=%s&accessLevel=2", url.QueryEscape(propertyId)), nil)
	if err != nil {
		return nil, fmt.Errorf("error building redfin api request: %w", err)
	}
	req.Header.Set("User-Agent", "Redfin/402.1.0.6002 CFNetwork/1240.0.4 Darwin/20.5.0")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying redfin api: %w", err)
	}
	defer resp.Body.Close()
	var result *AVMResponse

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading redfin api response: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error querying redfin api: %s", string(b))
	}
	err = json.Unmarshal([]byte(strings.TrimPrefix(string(b), "{}&&")), &result)
	if err != nil {
		return nil, fmt.Errorf("error json decoding redfin api response: %w", err)
	}

	return result, nil
}

func mustConfigureAuthCookies(httpClient *http.Client) {
	if httpClient.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("failed to initialize redfin client cookie jar: %v", err)
		}
		httpClient.Jar = jar
	}

	httpClient.Jar.SetCookies(
		&url.URL{Scheme: "http", Host: "www.redfin.com"},
		[]*http.Cookie{{
			Name:  "aws-waf-token",
			Value: uuid.New().String(),
		}},
	)
}
