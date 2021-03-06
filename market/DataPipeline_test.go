package market

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEtherFetch(t *testing.T) {
	cm := MakeCoinMarket()
	tickerMap, err := cm.FetchCoins("eth")
	fmt.Println(tickerMap)
	ticker := tickerMap["eth"]
	assert.NotNil(t, ticker)
	assert.Nil(t, err)
	assert.Condition(t, func() bool { return ticker.LastUpdated > 1000 })
	assert.Equal(t, ticker.Symbol, "ETH")
	assert.Equal(t, ticker.ID, "ethereum")
	assert.Equal(t, ticker.Name, "Ethereum")
}
