package currencyexternalservice

import (
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config/freecurrency"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config/ghgist"
)

func CreateService() *Service {
	primaryApi := currencyexternalapi.NewCurrencyAPI(freecurrency.FreeCurrencyAPIConfig)
	secondaryApi := currencyexternalapi.NewCurrencyAPI(ghgist.GHGistCurrencyAPIConfig)

	service := NewService(primaryApi, secondaryApi)

	return service
}
