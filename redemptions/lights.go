package redemptions

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/crazy3lf/colorconv"
	"github.com/nicklaw5/helix"
)

func (clients ClientHolder) lightControl_Redemption(rewardEvent helix.EventSubChannelPointsCustomRewardRedemptionEvent) {
	color := rewardEvent.UserInput
	strippedColor := strings.Replace(color, "#", "", -1)

	conv, err := colorconv.HexToColor(strippedColor)
	for err != nil {
		conv, err = colorconv.HexToColor(randomcolor.GetRandomColorInHex())
	}
	h1, s1, v1 := colorconv.ColorToHSV(conv)

	hue := int(h1 * (65535 / 360))
	sat := int(s1 * 255)
	val := int(v1 * 255)
	callURL := os.Getenv("HUE_URL") + `api/` + os.Getenv("HUE_USER") + `/groups/0/action`
	dataString := fmt.Sprintf(`{"sat":%d, "bri":%d,"hue":%d, "alert": "select"}`, sat, val, hue)
	data := []byte(dataString)

	req, err := http.NewRequest(http.MethodPut, callURL, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
