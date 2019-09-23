package level3stream

import (
	"encoding/json"
)

type StreamDataModel struct {
	Sequence   uint64 `json:"sequence"`
	Symbol     string `json:"symbol"`
	Type       string `json:"type"`
	rawMessage json.RawMessage
}

func NewStreamDataModel(msgData json.RawMessage) (*StreamDataModel, error) {
	l3Data := &StreamDataModel{}

	if err := json.Unmarshal(msgData, l3Data); err != nil {
		return nil, err
	}
	l3Data.rawMessage = msgData

	return l3Data, nil
}

func (l3Data *StreamDataModel) GetRawMessage() json.RawMessage {
	return l3Data.rawMessage
}

const (
	BuySide  = "buy"
	SellSide = "sell"

	MessageReceivedType = "received"
	MessageOpenType     = "open"
	MessageDoneType     = "done"
	MessageMatchType    = "match"
	MessageChangeType   = "update"
)

type StreamDataReceivedModel struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

type StreamDataOpenModel struct {
	Side    string `json:"side"`
	Size    int64  `json:"size"`
	OrderId string `json:"orderId"`
	Price   string `json:"price"`
	Time    uint64 `json:"ts"`
}

type StreamDataDoneModel struct {
	OrderId string `json:"orderId"`
	Reason  string `json:"reason"`
	Time    uint64 `json:"ts"`
}

type StreamDataMatchModel struct {
	TradeId      string `json:"tradeId"`
	Side         string `json:"side"`
	TakerOrderId string `json:"takerOrderId"`
	MakerOrderId string `json:"makerOrderId"`
	Price        string `json:"price"`
	MatchSize    uint64 `json:"matchSize"`
	Size         int64  `json:"size"`
	Time         uint64 `json:"ts"`

	//OrderId   string `json:"orderId"`
}

type StreamDataChangeModel struct {
	OrderId string `json:"orderId"`
	Price   string `json:"price"`
	Size    int64  `json:"size"`
	OldSize int64  `json:"oldSize"`
	Time    uint64 `json:"ts"`
}
