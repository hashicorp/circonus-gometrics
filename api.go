package circonusgometrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Call Circonus API
func (m *CirconusMetrics) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)

	// default to SSL
	proto := "https://"
	// allow override with explict "http://" in ApiHost
	if m.ApiHost[0:5] == "http:" {
		proto = ""
	}

	url := fmt.Sprintf("%s%s%s", proto, m.ApiHost, reqPath)

	req, err := retryablehttp.NewRequest(reqMethod, url, dataReader)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] creating API request: %s %+v", url, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
	req.Header.Add("X-Circonus-App-Name", m.ApiApp)

	client := retryablehttp.NewClient()
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 50 * time.Millisecond
	client.RetryMax = 3
	client.Logger = m.Log

	resp, err := client.Do(req)
	if err != nil {
		standard_client := &http.Client{}
		dataReader.Seek(0, 0)
		standard_req, _ := http.NewRequest(reqMethod, url, dataReader)
		standard_req.Header.Add("Accept", "application/json")
		standard_req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
		standard_req.Header.Add("X-Circonus-App-Name", m.ApiApp)
		resp, err := standard_client.Do(standard_req)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if m.Debug {
				m.Log.Printf("[DEBUG] %v\n", string(body))
			}
			return nil, fmt.Errorf("[ERROR] %s", string(body))
		}
		return nil, fmt.Errorf("[ERROR] fetching %s: %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] reading body %+v", err)
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("API response code %d: %s", resp.StatusCode, string(body))
		if m.Debug {
			m.Log.Printf("[DEBUG] %s\n", msg)
		}

		return nil, fmt.Errorf("[ERROR] %s", msg)
	}

	return body, nil
}