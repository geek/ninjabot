package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/rodrigo-brito/ninjabot/pkg/model"
)

type Kucoin struct {
	srv *kucoin.ApiService
}

func NewKucoin() *Kucoin {
	// The following env variables are supported
	// export API_KEY=key
	// export API_SECRET=secret
	// export API_PASSPHRASE=passphrase
	srv := kucoin.NewApiServiceFromEnv()

	return &Kucoin{
		srv: srv,
	}
}

func periodToKucoin(period string) string {
	if strings.HasSuffix(period, "w") {
		return fmt.Sprintf("%seek", period)
	}

	if strings.HasSuffix(period, "d") {
		return fmt.Sprintf("%say", period)
	}

	if strings.HasSuffix(period, "h") {
		return fmt.Sprintf("%sour", period)
	}

	if strings.HasSuffix(period, "m") {
		return fmt.Sprintf("%sin", period)
	}

	return period
}

func parseKlines(res *kucoin.ApiResponse, symbol string) ([]model.Candle, error) {
	klines := kucoin.KLinesModel{}
	if err := res.ReadData(&klines); err != nil {
		return nil, err
	}

	candles := make([]model.Candle, 0)
	for _, kline := range klines {
		// [0] time	Start time of the candle cycle
		// [1] open	Opening price
		// [2] close	Closing price
		// [3] high	Highest price
		// [4] low	Lowest price
		// [5] volume	Transaction volume
		// [6] turnover	Transaction amount
		vals := *kline
		t, err := strconv.ParseInt(vals[0], 10, 0)
		// shouldn't happen, but better to have correct times
		if err != nil {
			return nil, err
		}

		open, err := strconv.ParseFloat(vals[1], 64)
		if err != nil {
			return nil, err
		}

		close, err := strconv.ParseFloat(vals[2], 64)
		if err != nil {
			return nil, err
		}

		high, err := strconv.ParseFloat(vals[3], 64)
		if err != nil {
			return nil, err
		}

		low, err := strconv.ParseFloat(vals[4], 64)
		if err != nil {
			return nil, err
		}

		volume, err := strconv.ParseFloat(vals[5], 64)
		if err != nil {
			return nil, err
		}

		candles = append(candles, model.Candle{
			Symbol: symbol,
			Time:   time.Unix(t, 0),
			Open:   open,
			Close:  close,
			Low:    low,
			High:   high,
			Volume: volume,
		})
	}

	return candles, nil
}

func (k *Kucoin) CandlesByPeriod(ctx context.Context, pair, period string, start, end time.Time) ([]model.Candle, error) {
	res, err := k.srv.KLines(pair, periodToKucoin(period), start.Unix(), end.Unix())
	if err != nil {
		return nil, err
	}

	return parseKlines(res, pair)
}

func (k *Kucoin) CandlesByLimit(ctx context.Context, pair, period string, limit int) ([]model.Candle, error) {
	res, err := k.srv.KLines(pair, periodToKucoin(period), 0, 0)
	if err != nil {
		return nil, err
	}

	return parseKlines(res, pair)
}

func (k *Kucoin) CandlesSubscription(pair, timeframe string) (chan model.Candle, chan error) {
	return nil, nil
}

func (k *Kucoin) Account() (model.Account, error) {
	return model.Account{}, nil
}

func (k *Kucoin) Position(symbol string) (asset, quote float64, err error) {
	return 0, 0, nil
}

func (k *Kucoin) Order(symbol string, id int64) (model.Order, error) {
	return model.Order{}, nil
}

func (k *Kucoin) OrderOCO(side model.SideType, symbol string, size, price, stop, stopLimit float64) ([]model.Order, error) {
	return nil, nil
}

func (k *Kucoin) OrderLimit(side model.SideType, symbol string, size float64, limit float64) (model.Order, error) {
	return model.Order{}, nil
}

func (k *Kucoin) OrderMarket(side model.SideType, symbol string, size float64) (model.Order, error) {
	return model.Order{}, nil
}

func (k *Kucoin) Cancel(model.Order) error {
	return nil
}
