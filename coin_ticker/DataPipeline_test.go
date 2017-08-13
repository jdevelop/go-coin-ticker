package coin_ticker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEtherFetch(t *testing.T) {
	cm := MakeCoinMarket()
	ticker, err := cm.FetchCoins("ethereum")
	assert.NotNil(t, ticker)
	assert.Nil(t, err)
	assert.Condition(t, func() bool { return ticker.LastUpdated > 1000 })
	assert.Equal(t, ticker.Symbol, "ETH")
	assert.Equal(t, ticker.Id, "ethereum")
	assert.Equal(t, ticker.Name, "Ethereum")
}
