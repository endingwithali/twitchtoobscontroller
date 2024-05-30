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
)

func (clients ClientHolder) TestFXN() error {
	// sceneClient := obsClient.SceneItems
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

	currentScene, err := sceneClient.GetCurrentProgramScene()
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}
	// fmt.Printf("Current Scene name: %s", currentScene.CurrentProgramSceneName)

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

	getSceneItemsSourceResponse, err := sceneItemsClient.GetSceneItemTransform(&sceneitems.GetSceneItemTransformParams{
		SceneItemId: &createdInput.SceneItemId,
		SceneName:   &currentScene.CurrentProgramSceneName,
		SceneUuid:   &currentScene.CurrentProgramSceneUuid,
	})
	if err != nil {
		log.Fatal(err)
		// return errors.New(err.Error())
	}

	transformParams := getSceneItemsSourceResponse.SceneItemTransform
	log.Println(transformParams)
	transformParams.Alignment = 0
	transformParams.ScaleX = 0.35
	transformParams.ScaleY = 0.35
	transformParams.BoundsWidth = 100
	transformParams.BoundsHeight = 100
	transformParams.BoundsType = "OBS_BOUNDS_NONE"
	transformParams.BoundsAlignment = 0
	// assuming a default resolution 1920x1080
	//the right way to do this is to programmatically get the height and width of the current scene and do the divison to get the accurate x and y location for any obs scene. however for this i am hard coding it in
	transformParams.PositionX = 1920.0 / 2
	transformParams.PositionY = 1080.0 / 2

	log.Println(transformParams)

	_, err = sceneItemsClient.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneItemId:        &createdInput.SceneItemId,
		SceneItemTransform: transformParams,
		SceneName:          &currentScene.CurrentProgramSceneName,
		SceneUuid:          &currentScene.CurrentProgramSceneUuid,
	})
	if err != nil {
		log.Fatal(err)
		// return errors.New(err.Error())
	}

	log.Println("Waiting...")
	time.Sleep(10 * time.Second)
	log.Println("Deleting Image")

	removeParams := inputs.NewRemoveInputParams().WithInputName(inputName).WithInputUuid(createdInput.InputUuid)
	_, err = inputClient.RemoveInput(removeParams)
	if err != nil {
		fmt.Println("CreatePanic")
		log.Fatal(err)
	}
	return nil
}
