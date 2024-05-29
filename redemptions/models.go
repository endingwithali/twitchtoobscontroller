package redemptions

import (
	"encoding/json"

	"github.com/andreykaipov/goobs"
	"github.com/nicklaw5/helix"
)

type EventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

type ClientHolder struct {
	OBSClient *goobs.Client
}
