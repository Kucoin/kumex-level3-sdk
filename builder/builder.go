package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/JetBlink/orderbook/base"
	"github.com/JetBlink/orderbook/level3"
	"github.com/Kucoin/kumex-market/helper"
	"github.com/Kucoin/kumex-market/kumex"
	"github.com/Kucoin/kumex-market/level3stream"
	"github.com/Kucoin/kumex-market/utils/log"
	"github.com/Kucoin/kumex-market/utils/recovery"
	"github.com/shopspring/decimal"
)

type Builder struct {
	apiService *kumex.KuMEX
	symbol     string
	lock       *sync.RWMutex
	Messages   chan json.RawMessage

	fullOrderBook *level3.OrderBook
}

func NewOrderOfDepth(orderId string, side string, price string, size float64, time float64, info interface{}) (order *level3.Order, err error) {
	if err := base.CheckSide(side); err != nil {
		return nil, err
	}

	priceValue, err := decimal.NewFromString(price)
	if err != nil {
		return nil, fmt.Errorf("NewOrder failed, price: `%s`, error: %v", price, err)
	}

	sizeValue := decimal.NewFromFloat(size)

	order = &level3.Order{
		OrderId: orderId,
		Side:    side,
		Price:   priceValue,
		Size:    sizeValue,
		Time:    uint64(time),
		Info:    info,
	}
	return
}

func NewBuilder(apiService *kumex.KuMEX, symbol string) *Builder {
	return &Builder{
		apiService: apiService,
		symbol:     symbol,
		lock:       &sync.RWMutex{},
		Messages:   make(chan json.RawMessage, helper.MaxMsgChanLen*1024),
	}
}

func (b *Builder) resetOrderBook() {
	b.lock.Lock()
	b.fullOrderBook = level3.NewOrderBook()
	b.lock.Unlock()
}

func (b *Builder) ReloadOrderBook() {
	defer recovery.Recover(func(stack string) {
		log.Error("ReloadOrderBook panic : %v", stack)
		b.ReloadOrderBook()
	})()

	log.Warn("start running ReloadOrderBook, symbol: %s", b.symbol)
	b.resetOrderBook()

	b.playback()

	for msg := range b.Messages {
		l3Data, err := level3stream.NewStreamDataModel(msg)
		if err != nil {
			panic(err)
		}
		b.updateFromStream(l3Data)
	}
}

func (b *Builder) playback() {
	log.Warn("prepare playback...")
	const tempMsgChanMaxLen = 10240000
	tempMsgChan := make(chan *level3stream.StreamDataModel, tempMsgChanMaxLen)
	firstSequence := uint64(0)
	var fullOrderBook *DepthResponse

	for msg := range b.Messages {
		l3Data, err := level3stream.NewStreamDataModel(msg)
		if err != nil {
			panic(err)
		}

		tempMsgChan <- l3Data

		if firstSequence == 0 {
			firstSequence = l3Data.Sequence
			log.Error("firstSequence: %d", firstSequence)
		}

		if len(tempMsgChan) > 5 {
			if fullOrderBook == nil {
				log.Warn("start getting full level3 order book data, symbol: %s", b.symbol)
				fullOrderBook, err = b.GetAtomicFullOrderBook()
				if err != nil {
					panic(err)
					continue
				}
				log.Error("got full level3 order book data, Sequence: %d", fullOrderBook.Sequence)
			}

			if fullOrderBook != nil && fullOrderBook.Sequence < firstSequence {
				log.Error("获取 %d 全量数据太小", fullOrderBook.Sequence)
				fullOrderBook = nil
				continue
			}

			if fullOrderBook != nil && fullOrderBook.Sequence <= l3Data.Sequence { //string camp
				log.Warn("sequence match, start playback, tempMsgChan: %d", len(tempMsgChan))

				b.lock.Lock()
				b.AddDepthToOrderBook(fullOrderBook)
				b.lock.Unlock()

				n := len(tempMsgChan)
				for i := 0; i < n; i++ {
					b.updateFromStream(<-tempMsgChan)
				}

				log.Warn("finish playback.")
				break
			}

			if len(tempMsgChan) > tempMsgChanMaxLen-5 {
				panic("playback failed, tempMsgChan is too long, retry...")
			}
		}
	}
}

func (b *Builder) AddDepthToOrderBook(depth *DepthResponse) {
	b.fullOrderBook = b.FormatDepthToOrderBook(depth)
}

func (b *Builder) FormatDepthToOrderBook(depth *DepthResponse) *level3.OrderBook {
	fullOrderBook := level3.NewOrderBook()
	fullOrderBook.Sequence = depth.Sequence

	for _, elem := range depth.Asks {
		order, err := NewOrderOfDepth(elem[1].(string), base.AskSide, elem[2].(string), elem[3].(float64), elem[4].(float64), nil)
		if err != nil {
			panic(err)
		}

		if err := fullOrderBook.AddOrder(order); err != nil {
			panic(err)
		}
	}

	for _, elem := range depth.Bids {
		order, err := NewOrderOfDepth(elem[1].(string), base.BidSide, elem[2].(string), elem[3].(float64), elem[4].(float64), nil)
		if err != nil {
			panic(err)
		}

		if err := fullOrderBook.AddOrder(order); err != nil {
			panic(err)
		}
	}

	return fullOrderBook
}

func (b *Builder) DepthResponse2FullOrderBook(atomicFullOrderBook *DepthResponse) (*FullOrderBook, error) {
	fullOrderBook := b.FormatDepthToOrderBook(atomicFullOrderBook)
	data, err := json.Marshal(fullOrderBook)
	if err != nil {
		return nil, err
	}

	ret := &FullOrderBook{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (b *Builder) updateFromStream(msg *level3stream.StreamDataModel) {
	//time.Now().UnixNano()
	//log.Info("msg: %s", string(msg.GetRawMessage()))

	b.lock.Lock()
	defer b.lock.Unlock()

	skip, err := b.updateSequence(msg)
	if err != nil {
		panic(err)
	}

	if !skip {
		b.updateOrderBook(msg)
	}
}

func (b *Builder) updateSequence(msg *level3stream.StreamDataModel) (bool, error) {
	fullOrderBookSequenceValue := b.fullOrderBook.Sequence

	if fullOrderBookSequenceValue+1 > msg.Sequence {
		return true, nil
	}

	if fullOrderBookSequenceValue+1 != msg.Sequence {
		return false, errors.New(fmt.Sprintf(
			"currentSequence: %d, msgSequence: %d, the sequence is not continuous, 当前chanLen: %d",
			b.fullOrderBook.Sequence,
			msg.Sequence,
			len(b.Messages),
		))
	}

	b.fullOrderBook.Sequence = msg.Sequence

	return false, nil
}

//todo 大单特别注意
func (b *Builder) updateOrderBook(msg *level3stream.StreamDataModel) {
	//[3]string{"orderId", "price", "size"}
	//var item = [3]string{msg.OrderId, msg.Price, msg.Size}

	switch msg.Type {
	case level3stream.MessageReceivedType:

	case level3stream.MessageOpenType:
		data := &level3stream.StreamDataOpenModel{}
		if err := json.Unmarshal(msg.GetRawMessage(), data); err != nil {
			panic(err)
		}

		if data.Price == "" || data.Size == 0 {
			return
		}

		side := ""
		switch data.Side {
		case level3stream.SellSide:
			side = base.AskSide
		case level3stream.BuySide:
			side = base.BidSide
		default:
			panic("error side: " + data.Side)
		}

		order, err := level3.NewOrder(data.OrderId, side, data.Price, strconv.FormatInt(data.Size, 10), data.Time, nil)
		if err != nil {
			panic(err)
		}
		if err := b.fullOrderBook.AddOrder(order); err != nil {
			panic(err)
		}

	case level3stream.MessageDoneType:
		data := &level3stream.StreamDataDoneModel{}
		if err := json.Unmarshal(msg.GetRawMessage(), data); err != nil {
			panic(err)
		}

		if err := b.fullOrderBook.RemoveByOrderId(data.OrderId); err != nil {
			panic(err)
		}

	case level3stream.MessageMatchType:
		data := &level3stream.StreamDataMatchModel{}
		if err := json.Unmarshal(msg.GetRawMessage(), data); err != nil {
			panic(err)
		}

		if err := b.fullOrderBook.ChangeOrder(data.MakerOrderId, decimal.New(data.Size, 0)); err != nil {
			panic(err)
		}

	case level3stream.MessageChangeType:
		data := &level3stream.StreamDataChangeModel{}
		if err := json.Unmarshal(msg.GetRawMessage(), data); err != nil {
			panic(err)
		}

		if err := b.fullOrderBook.ChangeOrder(data.OrderId, decimal.New(data.Size, 0)); err != nil {
			panic(err)
		}

	default:
		panic("error msg type: " + msg.Type)
	}
}

func (b *Builder) Snapshot() (*FullOrderBook, error) {
	data, err := b.SnapshotBytes()
	if err != nil {
		return nil, err
	}

	ret := &FullOrderBook{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (b *Builder) SnapshotBytes() ([]byte, error) {
	b.lock.RLock()
	data, err := json.Marshal(b.fullOrderBook)
	b.lock.RUnlock()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (b *Builder) GetPartOrderBook(number int) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("GetPartOrderBook panic : %v", r)
		}
	}()

	b.lock.RLock()
	defer b.lock.RUnlock()

	data, err := json.Marshal(map[string]interface{}{
		"sequence":   b.fullOrderBook.Sequence,
		base.AskSide: b.fullOrderBook.GetPartOrderBookBySide(base.AskSide, number),
		base.BidSide: b.fullOrderBook.GetPartOrderBookBySide(base.BidSide, number),
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}
