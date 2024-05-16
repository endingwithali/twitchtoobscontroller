package redemptions

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
)

/*
* Process input from REST request received from Twitch
 */
func (clients ClientHolder) Process(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer req.Body.Close()

	// part of the twitch handshake
	if !helix.VerifyEventSubNotification(os.Getenv("CALLBACK_SECRET"), req.Header, string(body)) {
		log.Println("no valid signature on subscription")
		return
	} else {
		log.Println("verified signature for subscription")
	}

	// unmarshall event sub value received from body
	var eventSubValues EventSubNotification
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&eventSubValues)
	if err != nil {
		log.Println(err)
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if eventSubValues.Challenge != "" {
		w.Write([]byte(eventSubValues.Challenge))
		return
	}

	// Pull Values from Channel Point Redemption event
	var rewardEvent helix.EventSubChannelPointsCustomRewardRedemptionEvent
	err = json.NewDecoder(bytes.NewReader(eventSubValues.Event)).Decode(&rewardEvent)
	if err != nil {
		panic(err)

	}

	switch rewardEvent.Reward.Title {
	case "HueBulbParty":
		clients.lightControl_Redemption(rewardEvent)
	case "MemeMaker":
		clients.memeGenerator_Redemption(rewardEvent)
	}
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
