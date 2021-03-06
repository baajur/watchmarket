package watchmarket

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"math"
	"strconv"
	"strings"
	"time"
)

type (
	Rate struct {
		Currency         string  `json:"currency"`
		PercentChange24h float64 `json:"percent_change_24h,omitempty"`
		Provider         string  `json:"provider,omitempty"`
		Rate             float64 `json:"rate"`
		Timestamp        int64   `json:"timestamp"`
	}

	Rates []Rate

	CoinType string

	Response struct {
		Currency string  `json:"currency"`
		Docs     Tickers `json:"docs"`
	}

	Ticker struct {
		Coin       uint      `json:"coin"`
		CoinName   string    `json:"-"`
		TokenId    string    `json:"token_id,omitempty"`
		CoinType   CoinType  `json:"type,omitempty"`
		Price      Price     `json:"price,omitempty"`
		LastUpdate time.Time `json:"-"`
		Error      string    `json:"error,omitempty"`
		Volume     float64   `json:"-"`
		MarketCap  float64   `json:"-"`
	}

	Price struct {
		Change24h float64 `json:"change_24h"`
		Currency  string  `json:"-"`
		Provider  string  `json:"provider,omitempty"`
		Value     float64 `json:"value"`
	}

	Tickers []Ticker

	Chart struct {
		Provider string       `json:"provider,omitempty"`
		Prices   []ChartPrice `json:"prices,omitempty"`
		Error    string       `json:"error,omitempty"`
	}

	ChartPrice struct {
		Price float64 `json:"price"`
		Date  int64   `json:"date"`
	}

	CoinDetails struct {
		Provider          string  `json:"provider,omitempty"`
		ProviderURL       string  `json:"provider_url,omitempty"`
		Vol24             float64 `json:"volume_24"`
		MarketCap         float64 `json:"market_cap"`
		CirculatingSupply float64 `json:"circulating_supply"`
		TotalSupply       float64 `json:"total_supply"`
		Info              *Info   `json:"info,omitempty"`
	}

	Info struct {
		Name             string       `json:"name,omitempty"`
		Website          string       `json:"website,omitempty"`
		SourceCode       string       `json:"source_code,omitempty"`
		WhitePaper       string       `json:"white_paper,omitempty"`
		Description      string       `json:"description,omitempty"`
		ShortDescription string       `json:"short_description,omitempty"`
		Research         string       `json:"research,omitempty"`
		Explorer         string       `json:"explorer,omitempty"`
		Socials          []SocialLink `json:"socials,omitempty"`
	}

	SocialLink struct {
		Name   string `json:"name"`
		Url    string `json:"url"`
		Handle string `json:"handle"`
	}
)

const (
	UnknownCoinID                 = 111111
	DefaultPrecision              = 10
	Coin                 CoinType = "coin"
	Token                CoinType = "token"
	DefaultCurrency               = "USD"
	DefaultMaxChartItems          = 64

	ErrNotFound   = "not found"
	ErrBadRequest = "bad request"
	ErrInternal   = "internal"

	coinPrefix  = 'c'
	tokenPrefix = 't'
)

func (d Chart) IsEmpty() bool {
	return len(d.Prices) == 0
}

func (i CoinDetails) IsEmpty() bool {
	return i.Info.Name == ""
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func TruncateWithPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UnixToDuration(unixTime uint) time.Duration {
	return time.Duration(unixTime * 1000000000)
}

func DurationToUnix(duration time.Duration) uint {
	return uint(duration.Seconds())
}

func ParseID(id string) (uint, string, error) {
	rawResult := strings.Split(id, "_")
	resLen := len(rawResult)
	if resLen < 1 {
		return 0, "", errors.E("bad ID")
	}

	coin, err := findCoinID(rawResult)
	if err != nil {
		return 0, "", errors.E("bad ID")
	}

	token := findTokenID(rawResult)

	if token != "" {
		return coin, token, nil
	}

	return coin, "", nil
}

func findCoinID(words []string) (uint, error) {
	for _, w := range words {
		if w[0] == coinPrefix {
			rawCoin := removeFirstChar(w)
			coin, err := strconv.Atoi(rawCoin)
			if err != nil {
				return 0, errors.E("bad coin")
			}
			return uint(coin), nil
		}
	}
	return 0, errors.E("no coin")
}

func findTokenID(words []string) string {
	for _, w := range words {
		if w[0] == tokenPrefix {
			token := removeFirstChar(w)
			if len(token) > 0 {
				return token
			}
			return ""
		}
	}
	return ""
}

func removeFirstChar(input string) string {
	if len(input) <= 1 {
		return ""
	}
	return string([]rune(input)[1:])
}

func BuildID(coin uint, token string) string {
	c := strconv.Itoa(int(coin))
	if token != "" {
		return string(coinPrefix) + c + "_" + string(tokenPrefix) + token
	}
	return string(coinPrefix) + c
}
