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
