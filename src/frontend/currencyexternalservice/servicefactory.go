package currencyexternalservice

import (
	"github.com/kurtosis-tech/online-boutique-demo/frontend/currencyexternalservice/config/freecurrency"
	"github.com/kurtosis-tech/online-boutique-demo/frontend/currencyexternalservice/config/ghgist"
)

func CreateService() *Service {
	primaryApi := NewCurrencyAPI(freecurrency.FreeCurrencyAPIConfig)
	secondaryApi := NewCurrencyAPI(ghgist.GHGistCurrencyAPIConfig)

	service := NewService(primaryApi, secondaryApi)

	return service
}
