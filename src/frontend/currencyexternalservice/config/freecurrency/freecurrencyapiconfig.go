package freecurrency

import (
	"fmt"
	"github.com/kurtosis-tech/online-boutique-demo/frontend/currencyexternalservice/config"
	"net/url"
	"strings"
	"time"
)

const (
	apiBaseURL          = "https://api.freecurrencyapi.com/v1/"
	apiKeyQueryParamKey = "apikey"
	//TODO make it dynamic config
	apiKeyQueryParamValue   = "fca_live_VKZlykCWEiFcpBHnw74pzd4vLi04q1h9JySbVHDF"
	currenciesQueryParamKey = "currencies"
	currenciesEndpointPath  = "currencies"
	latestRatesEndpointPath = "latest"
)

var FreeCurrencyAPIConfig = config.NewCurrencyAPIConfig(
	// saving the response for a week because app.freecurrencyapi.com has a low limit
	// and this is a demo project, it's not important to have the latest data
	168*time.Hour,
	getCurrenciesURL,
	getLatestRatesURL,
)

func getCurrenciesURL() (*url.URL, error) {
	currenciesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, currenciesEndpointPath)

	currenciesEndpointUrl, err := url.Parse(currenciesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	currenciesEndpointQuery := currenciesEndpointUrl.Query()

	currenciesEndpointQuery.Set(apiKeyQueryParamKey, apiKeyQueryParamValue)

	currenciesEndpointUrl.RawQuery = currenciesEndpointQuery.Encode()

	return currenciesEndpointUrl, nil
}

func getLatestRatesURL(from string, to string) (*url.URL, error) {
	latestRatesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, latestRatesEndpointPath)

	latestRatesEndpointUrl, err := url.Parse(latestRatesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	latestRatesEndpointQuery := latestRatesEndpointUrl.Query()

	currenciesQueryParamValue := strings.Join([]string{strings.ToUpper(from), strings.ToUpper(to)}, ",")

	latestRatesEndpointQuery.Set(apiKeyQueryParamKey, apiKeyQueryParamValue)
	latestRatesEndpointQuery.Set(currenciesQueryParamKey, currenciesQueryParamValue)

	latestRatesEndpointUrl.RawQuery = latestRatesEndpointQuery.Encode()

	return latestRatesEndpointUrl, nil
}
