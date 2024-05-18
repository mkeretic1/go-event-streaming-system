package eventstream

import "encoding/json"

// NewPublishAction creates a new Publish action for sending messages to a Stream.
func NewPublishAction(stream string, message string) (*Action, error) {
	data, err := json.Marshal(&Publish{
		Stream:  stream,
		Message: message,
	})
	if err != nil {
		return nil, err
	}
	return &Action{
		ActionType: ActionTypePublish,
		Data:       data,
	}, nil
}

// NewSubscribeAction creates a new Subscribe action for subscribing to a Stream.
func NewSubscribeAction(stream string, mode SubscriptionMode, afterId ...int) (*Action, error) {
	subscribe := &Subscribe{
		Stream:           stream,
		SubscriptionMode: mode,
	}
	if len(afterId) > 0 {
		subscribe.AfterId = afterId[0]
	}

	data, err := json.Marshal(subscribe)
	if err != nil {
		return nil, err
	}
	return &Action{
		ActionType: ActionTypeSubscribe,
		Data:       data,
	}, nil
}

// removeValueFromSlice removes all instances of a specific string from a slice.
func removeValueFromSlice(s []string, value string) []string {
	result := make([]string, 0, len(s))
	for _, item := range s {
		if item != value {
			result = append(result, item)
		}
	}
	return result
}
