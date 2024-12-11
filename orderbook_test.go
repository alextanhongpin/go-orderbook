package orderbook_test

import (
	"testing"

	"github.com/alextanhongpin/go-orderbook"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	ob := orderbook.New()

	is := assert.New(t)
	is.Nil(ob.Add(orderbook.Order{
		ID:     "1",
		Type:   orderbook.Buy,
		Price:  100,
		Amount: 1,
	}))
	is.Nil(ob.Add(orderbook.Order{
		ID:     "2",
		Type:   orderbook.Buy,
		Price:  200,
		Amount: 1,
	}))
	is.Nil(ob.Add(orderbook.Order{
		ID:     "3",
		Type:   orderbook.Buy,
		Price:  150,
		Amount: 1,
	}))

	buys, sells := ob.Get()
	is.Len(buys, 3, "added 3 buy orders")
	is.Len(sells, 0, "no sell orders")

	is.Equal("2", buys[0].ID, "highest price first")
	is.Equal("3", buys[1].ID)
	is.Equal("1", buys[2].ID, "lowest price last")
}

func TestCancel(t *testing.T) {
	ob := orderbook.New()

	is := assert.New(t)
	var orderErr *orderbook.OrderError
	is.ErrorAs(ob.Cancel("1"), &orderErr, "order not found")
	is.Equal("1", orderErr.ID)
	is.Equal("not found", orderErr.Reason)

	is.Nil(ob.Add(orderbook.Order{
		ID:     "1",
		Type:   orderbook.Buy,
		Price:  100,
		Amount: 1,
	}))
	is.Nil(ob.Cancel("1"))
}

func TestMatch(t *testing.T) {
	type args struct {
		buyAmount   float64
		buyPrice    float64
		sellAmount  float64
		sellPrice   float64
		tradeAmount float64
		tradePrice  float64
	}
	test := func(name string, a args) {
		t.Helper() // Required for nested tests.
		t.Run(name, func(t *testing.T) {
			t.Helper() // Also required.

			ob := orderbook.New()
			is := assert.New(t)
			is.Nil(ob.Add(orderbook.Order{
				ID:     "1",
				Type:   orderbook.Buy,
				Price:  a.buyPrice,
				Amount: a.buyAmount,
			}))
			is.Nil(ob.Add(orderbook.Order{
				ID:     "2",
				Type:   orderbook.Sell,
				Price:  a.sellPrice,
				Amount: a.sellAmount,
			}))

			trades, err := ob.Match()
			is.Len(trades, 1)
			is.Nil(err)

			trade := trades[0]
			is.Equal("1", trade.BuyOrderID)
			is.Equal("2", trade.SellOrderID)
			is.Equal(a.tradePrice, trade.Price)
			is.Equal(a.tradeAmount, trade.Amount)
		})
	}

	test("equal", args{
		buyAmount:   100,
		buyPrice:    100,
		sellAmount:  100,
		sellPrice:   100,
		tradeAmount: 100,
		tradePrice:  100,
	})
	test("more buy amount", args{
		buyAmount:   100,
		buyPrice:    100,
		sellAmount:  1,
		sellPrice:   100,
		tradeAmount: 1,
		tradePrice:  100,
	})
	test("more sell amount", args{
		buyAmount:   1,
		buyPrice:    100,
		sellAmount:  100,
		sellPrice:   100,
		tradeAmount: 1,
		tradePrice:  100,
	})
}
