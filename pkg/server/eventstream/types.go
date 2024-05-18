package eventstream

import (
	"encoding/json"
	"net"
)

// SubscriptionMode defines the modes of subscription for messages within the event stream.
type SubscriptionMode string

const (
	SubscriptionModeAll   SubscriptionMode = "SUBSCRIPTION_MODE_ALL"
	SubscriptionModeNew   SubscriptionMode = "SUBSCRIPTION_MODE_NEW"
	SubscriptionModeAfter SubscriptionMode = "SUBSCRIPTION_MODE_AFTER"
)

// ActionType specifies the type of action being requested.
type ActionType string

const (
	ActionTypePublish   ActionType = "PUBLISH"
	ActionTypeSubscribe ActionType = "SUBSCRIBE"
)

// Action represents a generic action within the event stream, used for handling polymorphic message types.
type Action struct {
	ActionType ActionType      `json:"actionType"`
	Data       json.RawMessage `json:"data"`
}

// Publish is used to send messages to a specified stream.
type Publish struct {
	Stream  string `json:"stream"`
	Message string `json:"message"`
}

// Subscribe is used to request subscriptions to a specified stream, supporting different subscription modes.
type Subscribe struct {
	Stream           string           `json:"stream"`
	SubscriptionMode SubscriptionMode `json:"subscriptionMode"`
	AfterId          int              `json:"afterId"` // Specifies the Message ID after which messages should be received.
}

// EsStream represents a stream within the event system, containing events and subscribers.
type EsStream struct {
	Events      []EsEvent `json:"events"`
	Subscribers []string
}

// EsEvent defines an individual event within a stream.
type EsEvent struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

// EsClient represents a client connected to the event system.
type EsClient struct {
	uid  string
	conn net.Conn
}
