package events

import (
	"encoding/json"
	"sync"

	"github.com/Kucoin/kumex-level3-sdk/helper"
	"github.com/Kucoin/kumex-level3-sdk/level3stream"
	"github.com/Kucoin/kumex-level3-sdk/service"
	"github.com/Kucoin/kumex-level3-sdk/utils/log"
)

type Watcher struct {
	Messages  chan json.RawMessage
	redisPool *service.Redis
	lock      *sync.RWMutex

	orderIds   map[string]map[string]bool
	clientOids map[string]map[string]bool
}

func NewWatcher(redisPool *service.Redis) *Watcher {
	return &Watcher{
		Messages:  make(chan json.RawMessage, helper.MaxMsgChanLen),
		redisPool: redisPool,
		lock:      &sync.RWMutex{},

		orderIds:   make(map[string]map[string]bool),
		clientOids: make(map[string]map[string]bool),
	}
}

func (w *Watcher) Run() {
	log.Warn("start running Watcher")

	for msg := range w.Messages {
		if !w.existEventOrderIds() {
			continue
		}

		l3Data, err := level3stream.NewStreamDataModel(msg)
		if err != nil {
			panic(err)
		}

		switch l3Data.Type {
		case level3stream.MessageReceivedType:
			data := &level3stream.StreamDataReceivedModel{}
			if err := json.Unmarshal(l3Data.GetRawMessage(), data); err != nil {
				panic(err)
			}

			w.migrationClientOidToOrderIds(data.ClientOid, data.OrderId)

			w.publish(data.OrderId, string(l3Data.GetRawMessage()))

		case level3stream.MessageOpenType:
			data := &level3stream.StreamDataOpenModel{}
			if err := json.Unmarshal(l3Data.GetRawMessage(), data); err != nil {
				panic(err)
			}

			w.publish(data.OrderId, string(l3Data.GetRawMessage()))

		case level3stream.MessageMatchType:
			data := &level3stream.StreamDataMatchModel{}
			if err := json.Unmarshal(l3Data.GetRawMessage(), data); err != nil {
				panic(err)
			}

			w.publish(data.MakerOrderId, string(l3Data.GetRawMessage()))
			w.publish(data.TakerOrderId, string(l3Data.GetRawMessage()))

		case level3stream.MessageDoneType:
			data := &level3stream.StreamDataDoneModel{}
			if err := json.Unmarshal(l3Data.GetRawMessage(), data); err != nil {
				panic(err)
			}

			w.publish(data.OrderId, string(l3Data.GetRawMessage()))
			w.removeEventOrderId(data.OrderId)

		case level3stream.MessageChangeType:
			data := &level3stream.StreamDataChangeModel{}
			if err := json.Unmarshal(l3Data.GetRawMessage(), data); err != nil {
				panic(err)
			}

			w.publish(data.OrderId, string(l3Data.GetRawMessage()))

		default:
			panic("error msg type: " + l3Data.Type)
		}
	}
}

func (w *Watcher) migrationClientOidToOrderIds(clientOid, orderId string) {
	w.lock.RLock()
	channelsMap, ok := w.clientOids[clientOid]
	var channels []string
	if ok {
		channels = getMapKeys(channelsMap)
	}
	w.lock.RUnlock()
	if ok {
		w.removeEventClientOid(clientOid)
		w.AddEventOrderIdsToChannels(map[string][]string{
			orderId: channels,
		})
	}
}

func getMapKeys(data map[string]bool) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	return keys
}

func (w *Watcher) publish(orderId string, message string) {
	w.lock.RLock()
	channelsMap, ok := w.orderIds[orderId]
	var channels []string
	if ok {
		channels = getMapKeys(channelsMap)
	}
	w.lock.RUnlock()

	if ok {
		for _, channel := range channels {
			if err := w.redisPool.Publish(channel, message); err != nil {
				log.Error("redis publish to %s, msg: %s, error: %s", channel, message, err.Error())
				return
			}
		}
	}
}

func (w *Watcher) existEventOrderIds() bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if len(w.orderIds) == 0 && len(w.clientOids) == 0 {
		return false
	}

	return true
}

func (w *Watcher) AddEventOrderIdsToChannels(data map[string][]string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	for orderId, channels := range data {
		for _, channel := range channels {
			if w.orderIds[orderId] == nil {
				w.orderIds[orderId] = make(map[string]bool)
			}
			w.orderIds[orderId][channel] = true
		}
	}
}

func (w *Watcher) AddEventClientOidsToChannels(data map[string][]string) {
	w.lock.Lock()
	for clientOid, channels := range data {
		for _, channel := range channels {
			if w.clientOids[clientOid] == nil {
				w.clientOids[clientOid] = make(map[string]bool)
			}
			w.clientOids[clientOid][channel] = true
		}
	}
	w.lock.Unlock()
}

func (w *Watcher) removeEventOrderId(orderId string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	delete(w.orderIds, orderId)
}

func (w *Watcher) removeEventClientOid(clientOid string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	delete(w.clientOids, clientOid)
}
