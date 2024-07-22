package currencyexternalservice

import (
	"context"
	"github.com/kurtosis-tech/online-boutique-demo/frontend/currencyexternalservice/config/ghgist"
	demo "github.com/kurtosis-tech/online-boutique-demo/frontend/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test(t *testing.T) {
	currencyAPI := NewCurrencyAPI(ghgist.GHGistCurrencyAPIConfig)

	supported, err := currencyAPI.GetSupportedCurrencies(context.Background())
	require.NoError(t, err)
	require.NotNil(t, supported)

	in := &demo.CurrencyConversionRequest{
		From: &demo.Money{
			CurrencyCode: "USD",
		},
		ToCode: "BRL",
	}

	code, units, nanos, err := currencyAPI.Convert(context.Background(), in.From.CurrencyCode, in.From.Units, in.From.Nanos, in.ToCode)
	require.NoError(t, err)
	require.NotNil(t, code)
	require.NotNil(t, units)
	require.NotNil(t, nanos)
}
