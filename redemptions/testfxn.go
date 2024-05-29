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
	"github.com/andreykaipov/goobs/api/requests/sources"
)

func (clients ClientHolder) TestFXN() error {
	// sceneClient := obsClient.SceneItems
	inputClient := clients.OBSClient.Inputs
	sourceClient := clients.OBSClient.Sources
	sceneClient := clients.OBSClient.Scenes
	// sceneItemsClient := clients.OBSClient.SceneItems

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
	memeStringArr := strings.Split("BEANS A / BEANS B", "/")
	if len(memeStringArr) != 2 {
		log.Fatal("String too long")
		return errors.New("STRING MISFORMATTEED LONG")
	}

	postProcessFileName := clients.addText(preFileName, memeStringArr)
	fmt.Println(postProcessFileName)

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
	fmt.Printf("Current Scene name: %s", currentScene.CurrentProgramSceneName)

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
		fmt.Println("CreatePanic")
		log.Fatal(err)
	}

	fmt.Println("Waiting...")
	time.Sleep(10 * time.Second)
	fmt.Println("Deleting Image")

	removeParams := inputs.NewRemoveInputParams().WithInputName(inputName).WithInputUuid(createdInput.InputUuid)
	_, err = inputClient.RemoveInput(removeParams)
	if err != nil {
		fmt.Println("CreatePanic")
		log.Fatal(err)
	}
	return nil
}
