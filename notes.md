# Custom Plug Ins for twitch.tv/endingwithali 


type of event
EventSubTypeChannelPointsCustomRewardRedemptionAdd
info on that redemption:
https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_custom_reward_redemptionadd 


reauth
https://id.twitch.tv/oauth2/authorize
    ?response_type=token
    &client_id=0m0zg67tw1i3mmavvc7sdjqgk5ynh1
    &redirect_uri=https://localhost
    &scope=channel%3Amanage%3Aredemptions



# GOOBS EXPLORATION

HOW TO MAP EVENTS TO REQUESTS
https://github.com/andreykaipov/goobs/blob/main/docs/request-mapping.md


Events we are about
got event: *events.InputCreated
got event: *events.SceneItemSelected
got event: *events.SceneItemCreated
got event: *events.SceneItemRemoved
got event: *events.InputRemoved



input created: &{map[linear_alpha:false unload:false] image_source Image map[]  image_source}
SceneItemsSelected: &{6 Scene 2 }
SceneItemsCreated: &{6 0 Scene 2  Image }
SceneItemsRemoved: &{6 Scene 2  Image }
InputRemoved: &{Image }




EVENTS TO GOOBS MAPPING
*events.InputCreated - client.Inputs.CreateInput(...)
*events.SceneItemSelected - client.SceneItems.CreateSceneItem(...)
*events.SceneItemCreated
*events.SceneItemRemoved
*events.InputRemoved

https://stackoverflow.com/questions/78003747/obs-websocket-change-source-text-and-image


Scene Item Selected?
A scene item has been selected in the Ui.

Complexity Rating: 2/5
Latest Supported RPC Version: 1
Added in v5.0.0
Data Fields:

Name	Type	Description
sceneName	String	Name of the scene the item is in
sceneUuid	String	UUID of the scene the item is in
sceneItemId	Number	Numeric ID of the scene item

Scene Item Created
A scene item has been created.

Complexity Rating: 3/5
Latest Supported RPC Version: 1
Added in v5.0.0
Data Fields:

Name	Type	Description
sceneName	String	Name of the scene the item was added to
sceneUuid	String	UUID of the scene the item was added to
sourceName	String	Name of the underlying source (input/scene)
sourceUuid	String	UUID of the underlying source (input/scene)
sceneItemId	Number	Numeric ID of the scene item
sceneItemIndex	Number	Index position of the item





## Transforming scene items 
type SetSceneItemTransformParams  - https://pkg.go.dev/github.com/andreykaipov/goobs@v1.3.0/api/requests/sceneitems#SetSceneItemTransformParams
prob want to use the fucntion SetSceneItemTransform???  https://pkg.go.dev/github.com/andreykaipov/goobs@v1.3.0/api/requests/sceneitems#Client.SetSceneItemTransform
\

get scene item
generate transform param
setsceneitemtranform 
then show scene item 


insert scene item using createsceneitem  


// Send the input over to OBS
	createdInput, err := inputClient.CreateInput(inputParams)
	if err != nil {
		fmt.Println("CreatePanic")
		log.Fatal(err)
	}

returns 
type CreateInputResponse struct {

	// UUID of the newly created input
	InputUuid string `json:"inputUuid,omitempty"`

	// ID of the newly created scene item
	SceneItemId int `json:"sceneItemId,omitempty"`
	// contains filtered or unexported fields
}



```
	valueTrue := false
	_, err = sceneItemsClient.CreateSceneItem(&sceneitems.CreateSceneItemParams{
		SceneItemEnabled: &valueTrue,
		SceneName:        &currentScene.SceneName,
		SourceName:       &MEME_INPUT_NAME,
		SourceUuid:       &createdInput.InputUuid,
	})
	time.Sleep(10)



type SceneItemTransform ¶
type SceneItemTransform struct {
	Alignment       float64 `json:"alignment"`
	BoundsAlignment float64 `json:"boundsAlignment"`
	BoundsHeight    float64 `json:"boundsHeight"`
	BoundsType      string  `json:"boundsType"`
	BoundsWidth     float64 `json:"boundsWidth"`
	CropBottom      float64 `json:"cropBottom"`
	CropLeft        float64 `json:"cropLeft"`
	CropRight       float64 `json:"cropRight"`
	CropTop         float64 `json:"cropTop"`
	Height          float64 `json:"height"`
	PositionX       float64 `json:"positionX"`
	PositionY       float64 `json:"positionY"`
	Rotation        float64 `json:"rotation"`
	ScaleX          float64 `json:"scaleX"`
	ScaleY          float64 `json:"scaleY"`
	SourceHeight    float64 `json:"sourceHeight"`
	SourceWidth     float64 `json:"sourceWidth"`
	Width           float64 `json:"width"`
}
```


```
func (*Client) SetSceneItemTransform ¶
func (c *Client) SetSceneItemTransform(params *SetSceneItemTransformParams) (*SetSceneItemTransformResponse, error)
```