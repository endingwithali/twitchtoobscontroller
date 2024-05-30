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

func (clients ClientHolder) memeGenerator_Redemption(rewardEvent helix.EventSubChannelPointsCustomRewardRedemptionEvent) error {
	inputClient := clients.OBSClient.Inputs
	sourceClient := clients.OBSClient.Sources
	sceneClient := clients.OBSClient.Scenes
	sceneItemsClient := clients.OBSClient.SceneItems

	screenshotResponse, err := sourceClient.GetSourceScreenshot(&sources.GetSourceScreenshotParams{
		SourceName:              &FRONT_CAM_SOURCE,
		ImageCompressionQuality: &[]float64{-1}[0],
		ImageFormat:             &[]string{"png"}[0],
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
	//puts this into a relative file path - into generated_memes
	f, _ := os.Create(preFileLocation)
	_ = png.Encode(f, image)

	// Pulling STRING/STRING from twitch reqard
	memeStringArr := strings.Split(rewardEvent.UserInput, "/")
	if len(memeStringArr) != 2 {
		log.Fatal("String too long")
		return errors.New("STRING MISFORMATTEED LONG")
	}

	//For testing without having helix defined
	// // STRING/STRING
	// memeStringArr := strings.Split("BEANS A / BEANS B", "/")
	// if len(memeStringArr) != 2 {
	// 	log.Fatal("String too long")
	// 	return errors.New("STRING MISFORMATTEED LONG")
	// }

	postProcessFileName := clients.addText(preFileName, memeStringArr)
	log.Print(postProcessFileName)

	currentScene, err := sceneClient.GetCurrentProgramScene()
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}

	// log.Printf("Current Scene name: %s", currentScene.CurrentProgramSceneName)

	// Compose the new input
	inputName := fmt.Sprintf("%s_%d", MEME_INPUT_NAME, time.Now().Unix())
	inputParams := inputs.
		NewCreateInputParams().
		WithSceneName(currentScene.CurrentProgramSceneName).
		WithInputKind("image_source").
		WithInputName(inputName).
		WithInputSettings(map[string]interface{}{
			"file": postProcessFileName,
		})

	// Send the input over to OBS
	createdInput, err := inputClient.CreateInput(inputParams)
	if err != nil {
		panic(err)
	}

	getSceneItemsSourceResponse, err := sceneItemsClient.GetSceneItemTransform(&sceneitems.GetSceneItemTransformParams{
		SceneItemId: &createdInput.SceneItemId,
		SceneName:   &currentScene.CurrentProgramSceneName,
		SceneUuid:   &currentScene.CurrentProgramSceneUuid,
	})
	if err != nil {
		log.Fatal(err)
	}

	transformParams := getSceneItemsSourceResponse.SceneItemTransform
	log.Println(transformParams)

	// we can set the X and Y position to center the image to be the height and width because the image is a screenshot that is the proper height and width of the scene.
	transformParams.PositionX = transformParams.Width / 2
	transformParams.PositionY = transformParams.Height / 2
	transformParams.Alignment = 0
	transformParams.ScaleX = 0.35
	transformParams.ScaleY = 0.35
	transformParams.BoundsWidth = 100
	transformParams.BoundsHeight = 100

	log.Println(transformParams)

	_, err = sceneItemsClient.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneItemId:        &createdInput.SceneItemId,
		SceneItemTransform: transformParams,
		SceneName:          &currentScene.CurrentProgramSceneName,
		SceneUuid:          &currentScene.CurrentProgramSceneUuid,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting...")
	time.Sleep(10 * time.Second)
	log.Println("Deleting Image")

	removeParams := inputs.NewRemoveInputParams().WithInputName(inputName).WithInputUuid(createdInput.InputUuid)
	_, err = inputClient.RemoveInput(removeParams)
	if err != nil {
		log.Fatal(err)
	}
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

	if err := ggContext.LoadFontFace(fontLocation, 100); err != nil {
		log.Fatal("unable to load font")
		fmt.Println(err)
		panic(err)
	}

	return *ggContext
}
