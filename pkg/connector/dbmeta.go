package connector

import (
	"sync/atomic"

	waTypes "go.mau.fi/whatsmeow/types"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-meta/messagix/cookies"
	"go.mau.fi/mautrix-meta/messagix/table"
	"go.mau.fi/mautrix-meta/messagix/types"
	"go.mau.fi/mautrix-meta/pkg/metaid"
)

func (m *MetaConnector) GetDBMetaTypes() database.MetaTypes {
	return database.MetaTypes{
		Portal: func() any {
			return &PortalMetadata{}
		},
		Ghost: func() any {
			return &GhostMetadata{}
		},
		Message: func() any {
			return &MessageMetadata{}
		},
		Reaction: nil,
		UserLogin: func() any {
			return &UserLoginMetadata{}
		},
	}
}

type PortalMetadata struct {
	ThreadType     table.ThreadType `json:"thread_type"`
	WhatsAppServer string           `json:"whatsapp_server,omitempty"`

	fetchAttempted atomic.Bool
}

func (meta *PortalMetadata) JID(id networkid.PortalID) waTypes.JID {
	jid := metaid.ParseWAPortalID(id, meta.WhatsAppServer)
	if jid.Server == "" {
		switch meta.ThreadType {
		case table.ENCRYPTED_OVER_WA_GROUP:
			jid.Server = waTypes.GroupServer
		//case table.ENCRYPTED_OVER_WA_ONE_TO_ONE:
		//	jid.Server = waTypes.DefaultUserServer
		default:
			jid.Server = waTypes.MessengerServer
		}
	}
	return jid
}

type GhostMetadata struct {
	Username string `json:"username"`
}

type MessageMetadata struct {
	EditTimestamp int64 `json:"edit_timestamp,omitempty"`
}

type UserLoginMetadata struct {
	Platform   types.Platform   `json:"platform"`
	Cookies    *cookies.Cookies `json:"cookies"`
	WADeviceID uint16           `json:"wa_device_id,omitempty"`
}
