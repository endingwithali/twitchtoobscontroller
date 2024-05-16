package main

import (
	"fmt"
	"lightcontroler/redemptions"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
)

var MEME_FONT_LOCATION = "/assets/fonts/impact.ttf"
var IMG_HEIGHT = 800
var IMG_WIDTH = 600

/*
* Setting up listeners and clients for redemptions + OBS / Twitch Connections
 */
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// obs connection
	fmt.Println(os.Getenv("OBS_SERVER_PORT"))
	obsClient, err := goobs.New(fmt.Sprintf("%s:%s", os.Getenv("OBS_SERVER"), os.Getenv("OBS_SERVER_PORT")), goobs.WithPassword(os.Getenv("OBS_PASSWORD")))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer obsClient.Disconnect()

	ggContext := gg.NewContext(IMG_WIDTH, IMG_HEIGHT)
	// fontLocation, _ := filepath.Abs(MEME_FONT_LOCATION)
	pwd, _ := os.Getwd()
	fontLocation := pwd + MEME_FONT_LOCATION
	fmt.Println(fontLocation)
	if err := ggContext.LoadFontFace(fontLocation, 50); err != nil {
		log.Fatal("unable to load font")
		fmt.Print(err)
		panic(err)
	}

	// client handlers
	clients := redemptions.ClientHolder{
		OBSClient: obsClient,
		GGContext: ggContext,
	}
	defer obsClient.Disconnect()

	err = clients.TestFXN()
	if err != nil {
		panic(err)
	}
	// clients.Process()

	// // this needs to get deployed somewhere, and then have the callback url put below to handle receiving redirects
	// client, err := helix.NewClient(&helix.Options{
	// 	ClientID:     os.Getenv("TWITCH_CLIENT"),
	// 	ClientSecret: os.Getenv("TWITCH_SECRET"),
	// 	RedirectURI:  "https://localhost", //ngrok uri
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// tokenResponse, err := client.RequestAppAccessToken([]string{})
	// if err != nil {
	// 	panic(err)
	// }
	// client.SetAppAccessToken(tokenResponse.Data.AccessToken)

	// //	EventSubTypeChannelPointsCustomRewardRedemptionAdd    = "channel.channel_points_custom_reward_redemption.add"
	// _, err = client.CreateEventSubSubscription(&helix.EventSubSubscription{
	// 	Type:    helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
	// 	Version: "1",
	// 	Condition: helix.EventSubCondition{
	// 		BroadcasterUserID: os.Getenv("USER_ID"),
	// 	},
	// 	Transport: helix.EventSubTransport{ // using twitch event subscription, im telling it to use a url that i created to process the calls. when deployed change url here
	// 		Method:   "webhook",
	// 		Callback: "https://c682-2600-4041-54c0-9500-8c37-a651-7697-8db2.ngrok.io/process", //ngrok url
	// 		Secret:   os.Getenv("CALLBACK_SECRET"),
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// // now we're ready to make helix api calls
	// http.HandleFunc("/process", clients.Process)
	// http.ListenAndServe(":8090", nil)
}
