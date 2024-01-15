package matrix

import (
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

func IsValidRoomId(id string) bool {
	return strings.HasPrefix(id, "!") || strings.HasPrefix(id, "#")
}

func LoginToHomeServer(homeserver, username, password string) (*mautrix.Client, error) {

	client, err := mautrix.NewClient(homeserver, "", "")
	if err != nil {
		return nil, err
	}
	_, err = client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: username},
		Password:         password,
		StoreCredentials: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetValidAndJoinableRooms(client *mautrix.Client, roomAliases []string) []id.RoomID {

	rooms := make([]id.RoomID, 0)
	for _, rid := range roomAliases {
		var roomID id.RoomID
		if strings.HasPrefix(rid, "#") {
			resp, err := client.ResolveAlias(id.RoomAlias(rid))
			if err != nil {
				continue
			}
			roomID = resp.RoomID
		} else {
			roomID = id.RoomID(rid)
		}
		if _, err := client.JoinRoomByID(roomID); err != nil {
			continue
		} else {
			rooms = append(rooms, roomID)
		}
	}
	return rooms
}
