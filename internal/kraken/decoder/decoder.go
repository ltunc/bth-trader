package decoder

import (
	"bth-trader/internal/entities"
	"bytes"
	"encoding/json"
	"log"
	"strconv"
)

type msgType string

const (
	msgHeartbeat         msgType = "heartbeat"
	msgSubStatus         msgType = "subStatus"
	msgSysStatus         msgType = "sysStatus"
	msgAddOrderStatus    msgType = "addOrderStatus"
	msgCancelOrderStatus msgType = "cancelOrderStatus"
	msgOrder             msgType = "order"
	msgTrade             msgType = "trade"
	msgUnknown           msgType = "unknown"
)

// Outputs is a container for holding output channels for different messages from decoded stream
type Outputs struct {
	//HeartBeats chan []byte
	Orders chan *entities.Order
	Trades chan *entities.Trade
}

// DecodeStream decodes messages from channel,
// splits messages by their type and send them to appropriate output channel
func DecodeStream(in <-chan json.RawMessage, out *Outputs) {
	for m := range in {
		d := json.NewDecoder(bytes.NewReader(m))
		d.UseNumber()
		var rawData any
		err := d.Decode(&rawData)
		if err != nil {
			log.Printf("cannot decode message, error: %v", err)
			continue
		}
		switch detectType(rawData) {
		case msgHeartbeat:
			//log.Printf("heartbeat")
		case msgSubStatus:
			//log.Printf("subscribtion status")
		case msgSysStatus:
			log.Printf("system status")
		case msgAddOrderStatus:
			addOrder := parseAddOrderStatus(rawData)
			order := &entities.Order{
				RefId: addOrder.ReqId,
			}
			if addOrder.Status != "ok" {
				order.Status = "error"
				order.Error = addOrder.Error
			} else {
				order.Status = "open"
				order.OrderId = addOrder.TxId
			}
			out.Orders <- order
		case msgCancelOrderStatus:
			// processing of canceling order status is not a priority
			//log.Printf("cancel order status: %v", rawData)
		case msgOrder:
			for _, order := range parseOrders(rawData) {
				select {
				case out.Orders <- order:
					// the order is sent to output for further processing
				default:
					log.Printf("cannot send order to output channel, output is full.")
				}
			}
		case msgTrade:
			log.Printf("trade")
		case msgUnknown:
			log.Printf("unknown on unsupported message: %s", m)
		}
	}
}

func detectType(rawData any) msgType {
	if lstData, ok := rawData.([]any); ok {
		if str, ok := lstData[len(lstData)-2].(string); ok {
			switch str {
			case "ownTrades":
				return msgTrade
			case "openOrders":
				return msgOrder
			}
		}
	}
	if mapData, ok := rawData.(map[string]any); ok {
		if event, ok := mapData["event"]; ok {
			switch event {
			case "heartbeat":
				return msgHeartbeat
			case "subscriptionStatus":
				return msgSubStatus
			case "systemStatus":
				return msgSysStatus
			case "addOrderStatus", "addOrderStatusStatus":
				return msgAddOrderStatus
			case "cancelOrderStatus":
				return msgCancelOrderStatus
			}
		}
	}
	return msgUnknown
}

func parseOrders(rawData any) []*entities.Order {
	lstData, ok := rawData.([]any)
	if !ok {
		log.Printf("unexpected format of orders message, expected list: %v", rawData)
		return nil
	}
	rawOrders, ok := lstData[0].([]any)
	if !ok {
		log.Printf("wrong format of message, no list of orders: %T", lstData[0])
		return nil
	}
	var listOrders []*entities.Order
	for _, r := range rawOrders {
		orderMap, ok := r.(map[string]any)
		if !ok {
			log.Printf("unexpected format of orders map: %v", r)
			continue
		}
		for orderId, r := range orderMap {
			info, ok := r.(map[string]any)
			if !ok {
				log.Printf("unexpected format of order info: %#v, (from %v)", r, orderMap)
				continue
			}
			// We care only about the ID and the status of an order, and refId if present
			order := &entities.Order{
				OrderId: orderId,
			}
			if s, ok := info["status"]; ok {
				order.Status = s.(string)
			} else {
				// kraken sends messages to openOrders after each trade on the order,
				// but without status since only traded volume changing
				// we shall consider this as "open" status
				order.Status = "open"
			}
			if rawRef, ok := info["userref"]; ok {
				if refId, err := strconv.Atoi(rawRef.(json.Number).String()); err != nil {
					log.Printf("cannot parse userref from the order: %v", err)
				} else {
					order.RefId = refId
				}
			}
			listOrders = append(listOrders, order)
		}
	}
	return listOrders
}

type addOrderStatus struct {
	ReqId  int
	TxId   string
	Status string
	Error  string
}

func parseAddOrderStatus(rawData any) addOrderStatus {
	rawMap := rawData.(map[string]any)
	result := addOrderStatus{}
	result.Status = rawMap["status"].(string)
	reqId := rawMap["reqid"].(json.Number)
	result.ReqId, _ = strconv.Atoi(reqId.String())
	if txId, ok := rawMap["txid"]; ok {
		result.TxId = txId.(string)
	}
	if errMsg, ok := rawMap["errorMessage"]; ok {
		result.Error = errMsg.(string)
	}
	return result
}
