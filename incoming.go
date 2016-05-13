// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package tbotapi

import "fmt"
import "sort"

// BaseResponse contains the basic fields contained in every API response
type baseResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
}

// Audio represents an audio file to be treated as music
type Audio struct {
	FileBase
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
}

// Chat contains information about the chat a message originated from
type Chat struct {
	ID        int     `json:"id"`         // Unique identifier for this chat
	Type      string  `json:"type"`       // Type of chat, can be either "private", "group" or "channel". Check Is(PrivateChat|GroupChat|Channel)() methods
	Title     *string `json:"title"`      // Title for channels and group chats
	Username  *string `json:"username"`   // Username for private chats and channels if available
	FirstName *string `json:"first_name"` // First name of the other party in a private chat
	LastName  *string `json:"last_name"`  // Last name of the other party in a private chat
}

// IsPrivateChat checks if the chat is a private chat
func (c Chat) IsPrivateChat() bool {
	return c.Type == "private"
}

// IsGroupChat checks if the chat is a group chat
func (c Chat) IsGroupChat() bool {
	return c.Type == "group"
}

// IsSupergroup checks if the chat is a supergroup chat
func (c Chat) IsSupergroup() bool {
	return c.Type == "supergroup"
}

// IsChannel checks if the chat is a channel
func (c Chat) IsChannel() bool {
	return c.Type == "channel"
}

func (c Chat) String() string {
	toReturn := fmt.Sprint(c.ID)
	if c.IsPrivateChat() {
		toReturn += " (P) "
	} else if c.IsGroupChat() {
		toReturn += " (G) "
	} else {
		toReturn += " (C) "
	}

	if c.Title != nil {
		toReturn += "\"" + *c.Title + "\" "
	}

	if c.FirstName != nil {
		toReturn += *c.FirstName + " "
	}

	if c.LastName != nil {
		toReturn += *c.LastName + " "
	}

	if c.Username != nil {
		toReturn += "(@" + *c.Username + ")"
	}

	return toReturn
}

// Contact represents a phone contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	ID          int    `json:"user_id"`
}

// Document represents a general file
type Document struct {
	FileBase
	Thumbnail PhotoSize `json:"thumb"`
	Name      string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
}

// FileBase contains all the fields present in every file-like API response
type FileBase struct {
	ID   string `json:"file_id"`
	Size int    `json:"file_size"`
}

// File represents a file ready to be downloaded
type File struct {
	FileBase
	Path string `json:"file_path"`
}

// FileResponse represents the response sent by the API when requesting a file for download
type FileResponse struct {
	baseResponse
	File File `json:"result"`
}

// Location represents a point on the map
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// MessageResponse represents the response sent by the API on successful messages sent
type MessageResponse struct {
	baseResponse
	Message Message `json:"result"`
}

// Message represents a message
type Message struct {
	noReplyMessage
	ReplyToMessage *noReplyMessage `json:"reply_to_message"`
}

// IsForwarded checks if the message was forwarded
func (m *Message) IsForwarded() bool {
	return m.ForwardFrom != nil
}

// IsReply checks if the message is a reply
func (m *Message) IsReply() bool {
	return m.ReplyToMessage != nil
}

// Type determines the type of the message.
// Note that, for all these types, messages can still be replies or forwarded.
func (m *Message) Type() MessageType {
	if m.Text != nil {
		return TextMessage
	} else if m.Audio != nil {
		return AudioMessage
	} else if m.Document != nil {
		return DocumentMessage
	} else if m.Photo != nil {
		return PhotoMessage
	} else if m.Sticker != nil {
		return StickerMessage
	} else if m.Video != nil {
		return VideoMessage
	} else if m.Voice != nil {
		return VoiceMessage
	} else if m.Contact != nil {
		return ContactMessage
	} else if m.Location != nil {
		return LocationMessage
	} else if m.NewChatMember != nil {
		return NewChatMember
	} else if m.LeftChatMember != nil {
		return LeftChatMember
	} else if m.NewChatTitle != nil {
		return NewChatTitle
	} else if m.NewChatPhoto != nil {
		return NewChatPhoto
	} else if m.DeleteChatPhoto {
		return DeletedChatPhoto
	} else if m.GroupChatCreated {
		return GroupChatCreated
	} else if m.SupergroupChatCreated {
		return SupergroupChatCreated
	} else if m.ChannelChatCreated {
		return ChannelChatCreated
	} else if m.MigrateToChatID != nil {
		return MigrationToSupergroup
	} else if m.MigrateFromChatID != nil {
		return MigrationFromGroup
	} else if m.Venue != nil {
		return VenueMessage
	} else if m.PinnedMessage != nil {
		return PinnedMessage
	}

	return UnknownMessage
}

type noReplyMessage struct {
	Chat                  Chat             `json:"chat"`                    // information about the chat
	ID                    int              `json:"message_id"`              // message id
	From                  User             `json:"from"`                    // sender
	Date                  int              `json:"date"`                    // timestamp
	ForwardFrom           *User            `json:"forward_from"`            // forwarded from who
	ForwardDate           *int             `json:"forward_date"`            // forwarded from when
	Text                  *string          `json:"text"`                    // the actual text content
	Entities              *[]MessageEntity `json:"entities"`                // For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text (optional)
	Caption               *string          `json:"caption"`                 // caption for photo or video messages
	Audio                 *Audio           `json:"audio"`                   // information about audio contents
	Document              *Document        `json:"document"`                // information about file contents
	Photo                 *[]PhotoSize     `json:"photo"`                   // information about photo contents
	Sticker               *Sticker         `json:"sticker"`                 // information about sticker contents
	Video                 *Video           `json:"video"`                   // information about video contents
	Voice                 *Voice           `json:"voice"`                   // information about voice message contents
	Contact               *Contact         `json:"contact"`                 // information about contact contents
	Location              *Location        `json:"location"`                // information about location contents
	Venue                 *Venue           `json:"venue"`                   // information about venue contents
	NewChatMember         *User            `json:"new_chat_member"`         // information about a new chat participant
	LeftChatMember        *User            `json:"left_chat_member"`        // information about a chat participant who left
	NewChatTitle          *string          `json:"new_chat_title"`          // information about changes in the group name
	NewChatPhoto          *[]PhotoSize     `json:"new_chat_photo"`          // information about a new chat photo
	DeleteChatPhoto       bool             `json:"delete_chat_photo"`       // information about a deleted chat photo
	GroupChatCreated      bool             `json:"group_chat_created"`      // information about a created group chat
	SupergroupChatCreated bool             `json:"supergroup_chat_created"` // information about a created supergroup chat
	ChannelChatCreated    bool             `json:"channel_chat_created"`    // information about a created channel
	MigrateToChatID       *int             `json:"migrate_to_chat_id"`      // indicates the chat ID the group chat was migrated to (is now a supergroup)
	MigrateFromChatID     *int             `json:"migrate_from_chat_id"`    // indicates the chat ID the now supergroup chat was migrated from
	PinnedMessage         *noReplyMessage  `json:"pinned_message"`
}

// MessageEntityType is the type of an entity contained in a message.
type MessageEntityType string

// Entity types.
const (
	EntityTypeMention    = MessageEntityType("mention")
	EntityTypeHashtag    = MessageEntityType("hashtag")
	EntityTypeBotCommand = MessageEntityType("bot_command")
	EntityTypeURL        = MessageEntityType("url")
	EntityTypeEmail      = MessageEntityType("email")
	EntityTypeBold       = MessageEntityType("bold")
	EntityTypeItalic     = MessageEntityType("italic")
	EntityTypeCode       = MessageEntityType("code")
	EntityTypePre        = MessageEntityType("pre")
	EntityTypeTextLink   = MessageEntityType("text_link")
)

// MessageEntity represents an entity contained in a text message.
type MessageEntity struct {
	Type   MessageEntityType `json:"type"`
	Offset int               `json:"offset"`
	Length int               `json:"length"`
	URL    *string           `json:"url"`
}

// Venue represents a venue contained in a message.
type Venue struct {
	Location     Location `json:"location"`     // venue location
	Title        string   `json:"title"`        // Name of the venue
	Address      string   `json:"address"`      // Address of the venue
	FoursquareID *string  `json:"foursqare_id"` // Foursqare ID of the venue (optional)
}

// MessageType is the type of a message
type MessageType int

// Message types
const (
	TextMessage     MessageType = iota // text messages
	PinnedMessage                      // pinned messages
	AudioMessage                       // audio messages
	DocumentMessage                    // files
	PhotoMessage                       // photos
	StickerMessage                     // stickers
	VideoMessage                       // videos
	VoiceMessage                       // voice messages
	ContactMessage                     // contact information
	LocationMessage                    // locations
	VenueMessage                       // venues

	chatActionsBegin
	NewChatMember         // joined chat participants
	LeftChatMember        // left chat participants
	NewChatTitle          // chat title changes
	NewChatPhoto          // new chat photos
	DeletedChatPhoto      // deleted chat photos
	GroupChatCreated      // creation of a group chat
	SupergroupChatCreated // creation of a supergroup chat
	ChannelChatCreated    // createion of a channel
	MigrationToSupergroup // migration to supergroup
	MigrationFromGroup    // migration from group (to supergroup)
	chatActionsEnd

	UnknownMessage // unknown (probably new due to API changes)
)

var messageTypes = map[MessageType]string{
	TextMessage:     "Text",
	PinnedMessage:   "Pinned",
	AudioMessage:    "Audio",
	DocumentMessage: "Document",
	PhotoMessage:    "Photo",
	StickerMessage:  "Sticker",
	VideoMessage:    "Video",
	VoiceMessage:    "Voice",
	ContactMessage:  "Contact",
	LocationMessage: "Location",
	VenueMessage:    "Venue",

	NewChatMember:    "NewChatMember",
	LeftChatMember:   "LeftChatMember",
	NewChatTitle:     "NewChatTitle",
	NewChatPhoto:     "NewChatPhoto",
	DeletedChatPhoto: "DeletedChatPhoto",
	GroupChatCreated: "GroupChatCreated",

	UnknownMessage: "Unknown",
}

// IsChatAction checks if the MessageType is about changes in group chats
func (mt MessageType) IsChatAction() bool {
	return mt > chatActionsBegin && mt < chatActionsEnd
}

func (mt MessageType) String() string {
	val, ok := messageTypes[mt]
	if !ok {
		return messageTypes[UnknownMessage]
	}
	return val
}

// PhotoSize represents one size of a photo or a thumbnail
type PhotoSize struct {
	FileBase
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Sticker represents a sticker
type Sticker struct {
	FileBase
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Thumbnail PhotoSize `json:"thumb"`
}

// UpdateResponse represents the response sent by the API for a GetUpdates request
type updateResponse struct {
	baseResponse
	Update []Update `json:"result"`
}

// ByID is a wrapper to sort an []Update by ID
type byID []Update

func (a byID) Len() int           { return len(a) }
func (a byID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Sort sorts all the updates contained in an UpdateResponse by their ID
func (resp *updateResponse) sort() {
	sort.Sort(byID(resp.Update))
}

// Update represents an incoming update
type Update struct {
	ID                 int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
}

// Type returns the type of the update
func (u *Update) Type() UpdateType {
	if u.Message != nil {
		return MessageUpdate
	} else if u.InlineQuery != nil {
		return InlineQueryUpdate
	} else if u.ChosenInlineResult != nil {
		return ChosenInlineResultUpdate
	}
	return UnknownUpdate
}

// UpdateType represents an update type
type UpdateType int

// Update types
const (
	MessageUpdate            UpdateType = iota // message update
	InlineQueryUpdate                          // inline query
	ChosenInlineResultUpdate                   // chosen inline result

	UnknownUpdate // unkown, probably due to API changes
)

var updateTypes = map[UpdateType]string{
	MessageUpdate:            "Message",
	InlineQueryUpdate:        "InlineQuery",
	ChosenInlineResultUpdate: "ChosenInlineResult",

	UnknownUpdate: "Unknown",
}

func (t UpdateType) String() string {
	val, ok := updateTypes[t]
	if !ok {
		return updateTypes[UnknownUpdate]
	}
	return val
}

// UserResponse represents the response sent by the API on a GetMe request
type UserResponse struct {
	baseResponse
	User User `json:"result"`
}

// User represents a Telegram user or bot
type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"username"`
}

func (u User) String() string {
	if u.LastName != nil && u.Username != nil {
		return fmt.Sprintf("%d/%s %s (@%s)", u.ID, u.FirstName, *u.LastName, *u.Username)
	} else if u.LastName != nil {
		return fmt.Sprintf("%d/%s %s", u.ID, u.FirstName, *u.LastName)
	} else if u.Username != nil {
		return fmt.Sprintf("%d/%s (@%s)", u.ID, u.FirstName, *u.Username)
	}
	return fmt.Sprintf("%d/%s", u.ID, u.FirstName)
}

// UserProfilePhotosResponse represents the response sent by the API on a GetUserProfilePhotos request
type UserProfilePhotosResponse struct {
	baseResponse
	UserProfilePhotos UserProfilePhotos `json:"result"`
}

// UserProfilePhotos represents a users profile pictures
type UserProfilePhotos struct {
	TotalCount int         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

// Video represents a video file
type Video struct {
	FileBase
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	Thumbnail PhotoSize `json:"thumb"`
	MimeType  string    `json:"mime_type"`
	Caption   string    `json:"caption"`
}

// Voice represents a voice note
type Voice struct {
	FileBase
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
}

// InlineQuery represents an incoming inline query
type InlineQuery struct {
	ID     string `json:"id"`     // unique identifier for this query
	From   User   `json:"from"`   // sender
	Query  string `json:"query"`  // text of the query
	Offset string `json:"offset"` // offset of the results to be returned, can be controlled by the bot
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner
type ChosenInlineResult struct {
	ID    string `json:"result_id"` // unique identifier for the result that was chosen
	From  User   `json:"from"`      // user that chose the result
	Query string `json:"query"`     // query that was used to obtain the result
}

// CallbackQuery represents an incoming callback query from a button of an
// inline keyboard.
type CallbackQuery struct {
	ID              string   `json:"id"`                // Unique identifier for this query
	From            User     `json:"from"`              // Sender
	Message         *Message `json:"message"`           // Message with the callback button that originated the query. (optional)
	InlineMessageID *string  `json:"inline_message_id"` // Identifier of the message sent via the bot in inline mode, that originated the query
	Data            string   `json:"data"`              // Data associated with the callback button. Be aware that a bad client can send arbitrary data in this field.
}
