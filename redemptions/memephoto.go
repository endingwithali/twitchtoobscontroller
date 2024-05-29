package redemptions

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/sources"
	"github.com/fogleman/gg"
	"github.com/nicklaw5/helix"
)

var MEME_FONT_LOCATION = "/assets/fonts/impact.ttf"
var IMG_HEIGHT = 600
var IMG_WIDTH = 800
var FRONT_CAM_SOURCE = "CAM_front"
var MEME_INPUT_NAME = "INPUT_memeimage"

/*
	need to load image from base
	then load font
	then draw on image
	pulled from file
	cant use context to do so
	need to experiment with that

*/

func (clients ClientHolder) memeGenerator_Redemption(rewardEvent helix.EventSubChannelPointsCustomRewardRedemptionEvent) error {
	// sceneClient := obsClient.SceneItems
	inputClient := clients.OBSClient.Inputs
	sourceClient := clients.OBSClient.Sources
	sceneClient := clients.OBSClient.Scenes
	sceneItemsClient := clients.OBSClient.SceneItems

	screenshotResponse, err := sourceClient.GetSourceScreenshot(&sources.GetSourceScreenshotParams{
		SourceName:              &FRONT_CAM_SOURCE,
		ImageCompressionQuality: &[]float64{-1}[0],
		ImageFormat:             &[]string{"png"}[0],
		// ImageHeight:             &[]float64{float64(IMG_HEIGHT)}[0],
		// ImageWidth:              &[]float64{float64(IMG_WIDTH)}[0],
	})
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}

	// remove the `data:image/png;base64,...` prefix
	data := screenshotResponse.ImageData[strings.IndexByte(screenshotResponse.ImageData, ',')+1:]

	// save image
	preFileName := fmt.Sprintf("pre%d.png", time.Now().Unix())
	pwd, _ := os.Getwd()
	preFileLocation := fmt.Sprintf("%s/generated_memes/%s", pwd, preFileName)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	image, _ := png.Decode(reader)
	f, _ := os.Create(preFileLocation)
	_ = png.Encode(f, image)
	//trying to put this into a relative file path - into generated_memes

	// STRING/STRING
	memeStringArr := strings.Split(rewardEvent.UserInput, "/")
	if len(memeStringArr) != 2 {
		log.Fatal("String too long")
		return errors.New("STRING MISFORMATTEED LONG")
	}

	postProcessFileName := clients.addText(preFileName, memeStringArr)
	fmt.Print(postProcessFileName)

	/*
		1) get screen shot of source
			4) execute method
			5) store variable and check for errors
		2) pull message from redemption info
		3) add text to image
			https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go


		4) display image on scene
		5) sleep(10 sec)
		6) hide image from scene
		7) save image to local file - done in 	ggClient.SavePNG(postFileName)
		8) return

	*/

	// currentSceneInfo, err := sceneClient.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})

	currentScene, err := sceneClient.GetCurrentProgramScene()
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}

	/*
		use set input settings to select file location

	*/
	// Compose the new input
	inputParams := inputs.
		NewCreateInputParams().
		WithSceneName(currentScene.CurrentProgramSceneName).
		WithInputKind("image_source").
		WithInputName(MEME_INPUT_NAME).
		WithInputSettings(map[string]interface{}{
			"file": postProcessFileName,
		})

	// Send the input over to OBS
	createdInput, err := inputClient.CreateInput(inputParams)
	if err != nil {
		panic(err)
	}

	valueTrue := true
	_, err = sceneItemsClient.CreateSceneItem(&sceneitems.CreateSceneItemParams{
		SceneItemEnabled: &valueTrue,
		SceneName:        &currentScene.SceneName,
		SourceName:       &MEME_INPUT_NAME,
		SourceUuid:       &createdInput.InputUuid,
	})
	time.Sleep(10)

	/*
		how to display in a scene
		- get current scene name
		- create scene item using souce
		- insert into scene
		- wait
		- delete scene item
	*/
	return nil
}

func (clients ClientHolder) addText(src string, memeString []string) string {
	ggClient := generateGGContextWithImage(src)
	// ggClient.LoadPNG()
	height := float64(ggClient.Height())
	width := float64(ggClient.Width())
	line1 := strings.Trim(memeString[0], " ")
	line2 := strings.Trim(memeString[1], " ")
	lineHeight := 1.5
	lineWidth := width - 50
	outlineSize := 6
	for dy := -outlineSize; dy <= outlineSize; dy++ {
		for dx := -outlineSize; dx <= outlineSize; dx++ {
			if dx*dx+dy*dy >= outlineSize*outlineSize {
				// give it rounded corners
				continue
			}
			x := width/2 + float64(dx)
			y1 := height*0.115 + float64(dy)
			y2 := height*0.85 + float64(dy)
			ggClient.DrawStringAnchored(line1, x, y1, 0.5, 0.5)
			ggClient.DrawStringAnchored(line2, x, y2, 0.5, 0.5)
		}
	}
	ggClient.SetRGB(1, 1, 1)
	ggClient.DrawStringWrapped(line1, width*0.5, height*0.115, 0.5, 0.5, lineWidth, lineHeight, gg.AlignCenter)
	ggClient.DrawStringWrapped(line2, width*0.5, height*0.85, 0.5, 0.5, lineWidth, lineHeight, gg.AlignCenter)

	postFileName := fmt.Sprintf("post%d.png", time.Now().Unix())
	pwd, _ := os.Getwd()
	postFileLocation := fmt.Sprintf("%s/generated_memes/%s", pwd, postFileName)
	ggClient.SavePNG(postFileLocation)
	return postFileLocation
}

func generateGGContextWithImage(preFileName string) gg.Context {

	pwd, _ := os.Getwd()
	fontLocation := pwd + MEME_FONT_LOCATION

	preFileLocation := fmt.Sprintf("%s/generated_memes/%s", pwd, preFileName)
	img, err := gg.LoadImage(preFileLocation)
	if err != nil {
		log.Fatal("unable to images")
		fmt.Println(err)
		panic(err)
	}
	ggContext := gg.NewContextForImage(img)

	if err := ggContext.LoadFontFace(fontLocation, 50); err != nil {
		log.Fatal("unable to load font")
		fmt.Println(err)
		panic(err)
	}

	return *ggContext
}
