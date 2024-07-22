package externalcurrencyservice

import (
	"context"
	"encoding/json"
	"fmt"
	fep "github.com/kurtosis-tech/online-boutique-demo/frontend/proto"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"math"
	"net/http"
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

type CurrenciesResponse struct {
	Data map[string]Currency `json:"data"`
}

type Currency struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	SymbolNative  string `json:"symbol_native"`
	DecimalDigits int    `json:"decimal_digits"`
	Rounding      int    `json:"rounding"`
	Code          string `json:"code"`
	NamePlural    string `json:"name_plural"`
	Type          string `json:"type"`
}

type LatestRatesResponse struct {
	Data LatestRates `json:"data"`
}

type LatestRates map[string]float64

type FreeCurrencyAPI struct {
	httpClient *http.Client
	cache      *Cache
}

func NewFreeCurrencyAPI() *FreeCurrencyAPI {
	return &FreeCurrencyAPI{httpClient: http.DefaultClient, cache: NewCache()}
}

func (c *FreeCurrencyAPI) GetSupportedCurrencies(ctx context.Context, _ *fep.Empty, _ ...client.CallOption) (*fep.GetSupportedCurrenciesResponse, error) {

	currenciesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, currenciesEndpointPath)

	currenciesEndpointUrl, err := url.Parse(currenciesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	currenciesEndpointQuery := currenciesEndpointUrl.Query()

	currenciesEndpointQuery.Set(apiKeyQueryParamKey, apiKeyQueryParamValue)

	currenciesEndpointUrl.RawQuery = currenciesEndpointQuery.Encode()

	httpRequest := &http.Request{
		Method: http.MethodGet,
		URL:    currenciesEndpointUrl,
	}
	httpRequestWithContext := httpRequest.WithContext(ctx)

	httpResponseBodyBytes, err := c.doHttpRequest(httpRequestWithContext)
	if err != nil {
		return nil, err
	}

	currenciesResp := &CurrenciesResponse{}
	if err = json.Unmarshal(httpResponseBodyBytes, currenciesResp); err != nil {
		return nil, err
	}

	currencyCodes := []string{}

	for code := range currenciesResp.Data {
		currencyCodes = append(currencyCodes, code)
	}

	response := &fep.GetSupportedCurrenciesResponse{
		CurrencyCodes: currencyCodes,
	}

	return response, nil
}

func (c *FreeCurrencyAPI) Convert(ctx context.Context, in *fep.CurrencyConversionRequest, _ ...client.CallOption) (*fep.Money, error) {

	fromCode := strings.ToUpper(in.From.CurrencyCode)
	toCode := strings.ToUpper(in.ToCode)

	currencies, err := c.getLatestRatesFromAPI(ctx, fromCode, toCode)
	if err != nil {
		return nil, err
	}
	fromCurrency, found := currencies[fromCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.From.CurrencyCode)
	}
	toCurrency, found := currencies[toCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.ToCode)
	}

	money := &fep.Money{}

	money.CurrencyCode = in.ToCode
	total := int64(math.Floor(float64(in.From.Units*10^9+int64(in.From.Nanos)) / fromCurrency * toCurrency))
	money.Units = total / 1e9
	money.Nanos = int32(total % 1e9)

	return money, nil
}

func (c *FreeCurrencyAPI) getLatestRatesFromAPI(ctx context.Context, from string, to string) (map[string]float64, error) {
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

	httpRequest := &http.Request{
		Method: http.MethodGet,
		URL:    latestRatesEndpointUrl,
	}
	httpRequestWithContext := httpRequest.WithContext(ctx)

	httpResponseBodyBytes, err := c.doHttpRequest(httpRequestWithContext)
	if err != nil {
		return nil, err
	}

	latestRatesResp := &LatestRatesResponse{}
	if err = json.Unmarshal(httpResponseBodyBytes, latestRatesResp); err != nil {
		return nil, err
	}

	return latestRatesResp.Data, nil
}

func (c *FreeCurrencyAPI) doHttpRequest(
	request *http.Request,
) (
	resultResponseBodyBytes []byte,
	resultErr error,
) {

	var (
		httpResponseBodyBytes []byte
		err                   error
		ok                    bool
		urlStr                = request.URL.String()
	)

	if httpResponseBodyBytes, ok = c.cache.Get(urlStr); ok {
		fmt.Println("Cache hit for", urlStr)
		return httpResponseBodyBytes, nil
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusOK {
		httpResponseBodyBytes, err = io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	// saving the response for a week because app.freecurrencyapi.com has a low limit
	// and this is a demo project, it's not important to have the latest data
	c.cache.Set(urlStr, httpResponseBodyBytes, 168*time.Hour)

	return httpResponseBodyBytes, nil
}
