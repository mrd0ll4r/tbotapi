// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package tbotapi

import (
	"encoding/json"
	"fmt"
	"io"
)

type outgoingSetWebhook struct {
	URL string `json:"url"`
	outgoingFileBase
}

// querystring implements querystringer to represent the outgoing certificate
// file.
func (ow *outgoingSetWebhook) querystring() querystring {
	toReturn := make(map[string]string)

	if ow.URL != "" {
		toReturn["url"] = ow.URL
	}

	return querystring(toReturn)
}

// outgoingBase contains fields shared by most of the outgoing requests.
type outgoingBase struct {
	api       *TelegramBotAPI
	Recipient Recipient `json:"chat_id"`
}

// outgoingMessageBase contains fields shared by most of the outgoing
// messages.
type outgoingMessageBase struct {
	outgoingBase
	ReplyToMessageID    int         `json:"reply_to_message_id,omitempty"`
	ReplyMarkup         ReplyMarkup `json:"reply_markup,omitempty"`
	DisableNotification bool        `json:"disable_notification,omitempty"`
	replyToMessageIDSet bool
	replyMarkupSet      bool
}

// SetDisableNotification sets whether notifications should be disabled for
// this message (optional).
func (op *outgoingMessageBase) SetDisableNotification(to bool) {
	op.DisableNotification = to
}

// SetReplyToMessageID sets the ID for the message to reply to (optional).
func (op *outgoingMessageBase) SetReplyToMessageID(to int) {
	op.ReplyToMessageID = to
	op.replyToMessageIDSet = true
}

// SetReplyKeyboardMarkup sets the ReplyKeyboardMarkup (optional)
// Note that only one of ReplyKeyboardMarkup, ReplyKeyboardHide or
// ForceReply can be set.
// Attempting to set any of the other two or re-setting this will cause a
// panic.
func (op *outgoingMessageBase) SetReplyKeyboardMarkup(to ReplyKeyboardMarkup) {
	if op.replyMarkupSet {
		panic("Outgoing: Only one of ReplyKeyboardMarkup, ReplyKeyboardHide or ForceReply can be set")
	}

	op.ReplyMarkup = ReplyMarkup(to)
}

// SetReplyKeyboardHide sets the ReplyKeyboardHide (optional)
// Note that only one of ReplyKeyboardMarkup, ReplyKeyboardHide or
// ForceReply can be set.
// Attempting to set any of the other two or re-setting this will cause a
// panic.
func (op *outgoingMessageBase) SetReplyKeyboardHide(to ReplyKeyboardHide) {
	if !to.HideKeyboard {
		return
	}

	if op.replyMarkupSet {
		panic("Outgoing: Only one of ReplyKeyboardMarkup, ReplyKeyboardHide or ForceReply can be set")
	}

	op.ReplyMarkup = ReplyMarkup(to)
}

// SetForceReply sets ForceReply for this message (optional)
// Note that only one of ReplyKeyboardMarkup, ReplyKeyboardHide or
// ForceReply can be set.
// Attempting to set any of the other two or re-setting this will cause a
// panic.
func (op *outgoingMessageBase) SetForceReply(to ForceReply) {
	if !to.ForceReply {
		return
	}

	if op.replyMarkupSet {
		panic("Outgoing: Only one of ReplyKeyboardMarkup, ReplyKeyboardHide or ForceReply can be set")
	}

	op.ReplyMarkup = ReplyMarkup(to)
}

// getBaseQueryString gets a Querystring representing this message.
func (op *outgoingBase) getBaseQueryString() querystring {
	toReturn := map[string]string{}
	if op.Recipient.isChannel() {
		//Channel
		toReturn["chat_id"] = fmt.Sprint(*op.Recipient.ChannelID)
	} else {
		toReturn["chat_id"] = fmt.Sprint(*op.Recipient.ChatID)
	}

	return querystring(toReturn)
}

// getMessageBaseQueryString gets a Querystring representing this message.
func (op *outgoingMessageBase) getBaseQueryString() querystring {
	toReturn := map[string]string{}
	if op.Recipient.isChannel() {
		//Channel.
		toReturn["chat_id"] = fmt.Sprint(*op.Recipient.ChannelID)
	} else {
		toReturn["chat_id"] = fmt.Sprint(*op.Recipient.ChatID)
	}

	if op.replyToMessageIDSet {
		toReturn["reply_to_message_id"] = fmt.Sprint(op.ReplyToMessageID)
	}

	if op.replyMarkupSet {
		b, err := json.Marshal(op.ReplyMarkup)
		if err != nil {
			panic(err)
		}
		toReturn["reply_markup"] = string(b)
	}

	if op.DisableNotification {
		toReturn["disable_notification"] = fmt.Sprint(op.DisableNotification)
	}

	return querystring(toReturn)
}

type outgoingFileBase struct {
	fileName string
	r        io.Reader
	fileID   string
}

func (b outgoingFileBase) valid() bool {
	return b.isUpload() || b.isResend()
}

func (b outgoingFileBase) isUpload() bool {
	return b.fileName != "" && b.r != nil && b.fileID == ""
}

func (b outgoingFileBase) isResend() bool {
	return b.fileName == "" && b.r == nil && b.fileID != ""
}

// OutgoingAudio represents an outgoing audio file.
type OutgoingAudio struct {
	outgoingMessageBase
	outgoingFileBase
	Duration  int    `json:"duration,omitempty"`
	Title     string `json:"title,omitempty"`
	Performer string `json:"performer,omitempty"`
}

// SetDuration sets a duration for the audio file (optional).
func (oa *OutgoingAudio) SetDuration(to int) *OutgoingAudio {
	oa.Duration = to
	return oa
}

// SetPerformer sets a performer for the audio file (optional).
func (oa *OutgoingAudio) SetPerformer(to string) *OutgoingAudio {
	oa.Performer = to
	return oa
}

// SetTitle sets a title for the audio file (optional).
func (oa *OutgoingAudio) SetTitle(to string) *OutgoingAudio {
	oa.Title = to
	return oa
}

// querystring implements querystringer to represent the audio file.
func (oa *OutgoingAudio) querystring() querystring {
	toReturn := map[string]string(oa.getBaseQueryString())

	if oa.Duration != 0 {
		toReturn["duration"] = fmt.Sprint(oa.Duration)
	}

	if oa.Performer != "" {
		toReturn["performer"] = oa.Performer
	}

	if oa.Title != "" {
		toReturn["title"] = oa.Title
	}

	return querystring(toReturn)
}

// OutgoingDocument represents an outgoing file.
type OutgoingDocument struct {
	outgoingMessageBase
	outgoingFileBase
}

// querystring implements querystringer to represent the outgoing file.
func (od *OutgoingDocument) querystring() querystring {
	return od.getBaseQueryString()
}

// OutgoingForward represents an outgoing, forwarded message.
type OutgoingForward struct {
	outgoingMessageBase
	FromChatID Recipient `json:"from_chat_id"`
	MessageID  int       `json:"message_id"`
}

// OutgoingLocation represents an outgoing location on a map.
type OutgoingLocation struct {
	outgoingMessageBase
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// OutgoingVenue represents an outgoing venue message.
type OutgoingVenue struct {
	outgoingMessageBase
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Title        string  `json:"title"`
	Address      string  `json:"address"`
	FoursquareID string  `json:"foursquare_id,omitempty"`
}

// SetFoursquareID sets the foursquare ID for the venue (optional).
func (ov *OutgoingVenue) SetFoursquareID(to string) *OutgoingVenue {
	ov.FoursquareID = to
	return ov
}

// OutgoingMessage represents an outgoing message.
type OutgoingMessage struct {
	outgoingMessageBase
	Text                  string    `json:"text"`
	DisableWebPagePreview bool      `json:"disable_web_page_preview,omitempty"`
	ParseMode             ParseMode `json:"parse_mode,omitempty"`
}

// SetMarkdown sets or resets whether the message should be parsed as
// markdown or plain text (optional).
func (om *OutgoingMessage) SetMarkdown(to bool) *OutgoingMessage {
	if to {
		om.ParseMode = ModeMarkdown
	} else {
		om.ParseMode = ModeDefault
	}
	return om
}

// SetHTML sets or resets whether the message should be parsed as HTML or
// plain text (optional).
func (om *OutgoingMessage) SetHTML(to bool) *OutgoingMessage {
	if to {
		om.ParseMode = ModeHTML
	} else {
		om.ParseMode = ModeDefault
	}
	return om
}

// SetDisableWebPagePreview disables web page previews for the message
// (optional).
func (om *OutgoingMessage) SetDisableWebPagePreview(to bool) *OutgoingMessage {
	om.DisableWebPagePreview = to
	return om
}

// OutgoingPhoto represents an outgoing photo.
type OutgoingPhoto struct {
	outgoingMessageBase
	outgoingFileBase
	Caption string `json:"caption,omitempty"`
}

// SetCaption sets a caption for the photo (optional).
func (op *OutgoingPhoto) SetCaption(to string) *OutgoingPhoto {
	op.Caption = to
	return op
}

// querystring implements querystringer to represent the photo.
func (op *OutgoingPhoto) querystring() querystring {
	toReturn := map[string]string(op.getBaseQueryString())

	if op.Caption != "" {
		toReturn["caption"] = op.Caption
	}

	return querystring(toReturn)
}

// OutgoingSticker represents an outgoing sticker message.
type OutgoingSticker struct {
	outgoingMessageBase
	outgoingFileBase
}

// querystring implements querystringer to represent the sticker message.
func (os *OutgoingSticker) querystring() querystring {
	return os.getBaseQueryString()
}

// OutgoingKickChatMember represents a request to kick a chat member.
type OutgoingKickChatMember struct {
	api       *TelegramBotAPI
	Recipient Recipient `json:"chat_id"`
	UserID    int       `json:"user_id"`
}

// OutgoingUnbanChatMember represents a request to unban a chat member.
type OutgoingUnbanChatMember struct {
	api       *TelegramBotAPI
	Recipient Recipient `json:"chat_id"`
	UserID    int       `json:"user_id"`
}

// OutgoingRestrictChatMember represents a request to restrict a chat member.
type OutgoingRestrictChatMember struct {
	api                   *TelegramBotAPI
	Recipient             Recipient `json:"chat_id"`
	UserID                int       `json:"user_id"`
	UntilDate             int       `json:"until_date,omitempty"`
	CanSendMessages       bool      `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool      `json:"can_send_media_messages,omitempty"`
	CanSendOtherMessages  bool      `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool      `json:"can_add_web_page_previews,omitempty"`
}

// OutgoingCallbackQueryResponse represents a response to a callback query.
type OutgoingCallbackQueryResponse struct {
	api             *TelegramBotAPI
	CallbackQueryID string `json:"callback_query_id"` // ID of the callback query.
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
}

// OutgoingUserProfilePhotosRequest represents a request for a users
// profile photos.
type OutgoingUserProfilePhotosRequest struct {
	api    *TelegramBotAPI
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// SetOffset sets an offset for the request (optional).
func (op *OutgoingUserProfilePhotosRequest) SetOffset(to int) *OutgoingUserProfilePhotosRequest {
	op.Offset = to
	return op
}

// SetLimit sets a limit for the request (optional).
func (op *OutgoingUserProfilePhotosRequest) SetLimit(to int) *OutgoingUserProfilePhotosRequest {
	op.Limit = to
	return op
}

// querystring implements querystringer to represent the request.
func (op *OutgoingUserProfilePhotosRequest) querystring() querystring {
	toReturn := map[string]string{}
	toReturn["user_id"] = fmt.Sprint(op.UserID)

	if op.Offset != 0 {
		toReturn["offset"] = fmt.Sprint(op.Offset)
	}

	if op.Limit != 0 {
		toReturn["limit"] = fmt.Sprint(op.Limit)
	}

	return querystring(toReturn)
}

// OutgoingVideo represents an outgoing video file.
type OutgoingVideo struct {
	outgoingMessageBase
	outgoingFileBase
	Duration int    `json:"duration,omitempty"`
	Caption  string `json:"caption,omitempty"`
}

// SetCaption sets a caption for the video file (optional).
func (ov *OutgoingVideo) SetCaption(to string) *OutgoingVideo {
	ov.Caption = to
	return ov
}

// SetDuration sets a duration for the video file (optional).
func (ov *OutgoingVideo) SetDuration(to int) *OutgoingVideo {
	ov.Duration = to
	return ov
}

// querystring implements querystringer to represent the outgoing video
// file.
func (ov *OutgoingVideo) querystring() querystring {
	toReturn := map[string]string(ov.getBaseQueryString())

	if ov.Caption != "" {
		toReturn["caption"] = ov.Caption
	}

	if ov.Duration != 0 {
		toReturn["duration"] = fmt.Sprint(ov.Duration)
	}

	return querystring(toReturn)
}

// OutgoingVoice represents an outgoing voice note.
type OutgoingVoice struct {
	outgoingMessageBase
	outgoingFileBase
	Duration int `json:"duration,omitempty"`
}

// SetDuration sets a duration of the voice note (optional).
func (ov *OutgoingVoice) SetDuration(to int) *OutgoingVoice {
	ov.Duration = to
	return ov
}

// querystring implements querystringer to represent the outgoing voice
// note.
func (ov *OutgoingVoice) querystring() querystring {
	toReturn := map[string]string(ov.getBaseQueryString())

	if ov.Duration != 0 {
		toReturn["duration"] = fmt.Sprint(ov.Duration)
	}

	return querystring(toReturn)
}

// ReplyMarkup is s marker interface for ReplyMarkups.
// It is implemented by ReplyKeyboard(Hide|Markup) and ForceReply.
type ReplyMarkup interface {
	replyMarkup()
}

func (ReplyKeyboardHide) replyMarkup()    {}
func (ReplyKeyboardMarkup) replyMarkup()  {}
func (ForceReply) replyMarkup()           {}
func (InlineKeyboardMarkup) replyMarkup() {}

// InlineKeyboardMarkup represents markup for an inline query response.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents a button for an inline query keyboard.
type InlineKeyboardButton struct {
	Text              string `json:"text"`
	URL               string `json:"url,omitempty"`
	CallbackData      string `json:"callback_data,omitempty"`
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
}

// ForceReply represents the values sent by a bot so that clients will be
// presented with a forced reply,
// see https://core.telegram.org/bots/api#forcereply
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options to
// be presented to clients.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"` // Slice of keyboard lines.
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

// KeyboardButton represents a button on a reply keyboard.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// ReplyKeyboardHide contains the fields necessary to hide a custom
// keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

// ParseMode describes how a message should be parsed client-side.
type ParseMode string

//ParseModes.
const (
	ModeMarkdown = ParseMode("Markdown") // Parse as Markdown.
	ModeHTML     = ParseMode("HTML")     // Parse as HTML.
	ModeDefault  = ParseMode("")         // Parse as text.
)

// OutgoingChatAction represents an outgoing chat action.
type OutgoingChatAction struct {
	outgoingBase
	Action ChatAction `json:"action"`
}

// ChatAction represents an action to be shown to clients, indicating
// activity of the bot.
type ChatAction string

// Represents all the possible ChatActions to be sent,
// see https://core.telegram.org/bots/api#sendchataction
const (
	ChatActionTyping         ChatAction = "typing"
	ChatActionUploadPhoto               = "upload_photo"
	ChatActionRecordVideo               = "record_video"
	ChatActionUploadVideo               = "upload_video"
	ChatActionRecordAudio               = "record_audio"
	ChatActionUploadAudio               = "upload_audio"
	ChatActionUploadDocument            = "upload_document"
	ChatActionFindLocation              = "find_location"
)

// InlineQueryAnswer represents a response to an inline query.
// For limitations, check the API documentation.
type InlineQueryAnswer struct {
	api        *TelegramBotAPI
	QueryID    string              `json:"inline_query_id"`       // Unique identifier for the answered query.
	Results    []InlineQueryResult `json:"results"`               // Results for the query.
	CacheTime  int                 `json:"cache_time,omitempty"`  // The maximum amount of time in seconds that the result of the query may be cached.
	Personal   bool                `json:"is_personal,omitempty"` // If set to true, results will be cached for that user onl.
	NextOffset string              `json:"next_offset,omitempty"` // The offset that a client should send in the next query with the same text to receive more results.
}

// InlineQueryResultType represents a type of an inline query result.
type InlineQueryResultType string

// Inline query result type constants.
const (
	ArticleResult  = InlineQueryResultType("article")
	PhotoResult    = InlineQueryResultType("photo")
	GifResult      = InlineQueryResultType("gif")
	Mpeg4GifResult = InlineQueryResultType("mpeg4_gif")
	VideoResult    = InlineQueryResultType("video")
)

// InlineQueryResultBase is the base for all InlineQueryResults.
type InlineQueryResultBase struct {
	Type                  InlineQueryResultType `json:"type"`                               // Type of the result.
	ID                    string                `json:"id"`                                 // Unique identifier for this result, 1-64 bytes.
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`               // Indicates how to parse client-side (optional).
	DisableWebPagePreview bool                  `json:"disable_web_page_preview,omitempty"` // Disables link previews (optional).
}

// InlineQueryResult is a marker interface for query results.
// It is implemented by pointers to
// InlineQueryResult(Article|Photo|Gif|Mpeg4Gif|Video).
type InlineQueryResult interface {
	result()
}

func (*InlineQueryResultArticle) result()  {}
func (*InlineQueryResultPhoto) result()    {}
func (*InlineQueryResultGif) result()      {}
func (*InlineQueryResultMpeg4Gif) result() {}
func (*InlineQueryResultVideo) result()    {}

// InlineQueryResultArticle represents a link to an article or web page.
type InlineQueryResultArticle struct {
	InlineQueryResultBase
	Title       string `json:"title"`                  // Title of the result.
	Text        string `json:"message_text"`           // Text of the message to be sent.
	URL         string `json:"url,omitempty"`          // URL of the result (optional).
	HideURL     bool   `json:"hide_url,omitempty"`     // Whether to hide the URL in the message (optional).
	Description string `json:"description,omitempty"`  // Short description of the result (optional).
	ThumbURL    string `json:"thumb_url,omitempty"`    // URL of the thumbnail for the result (optional).
	ThumbWidth  int    `json:"thumb_width,omitempty"`  // Thumbnail width (optional).
	ThumbHeight int    `json:"thumb_height,omitempty"` // Thumbnail height (optional).
}

// NewInlineQueryResultArticle returns a new InlineQueryResultArticle with
// all mandatory fields set.
func NewInlineQueryResultArticle(id, title, text string) *InlineQueryResultArticle {
	return &InlineQueryResultArticle{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: ArticleResult,
			ID:   id,
		},
		Title: title,
		Text:  text,
	}
}

// InlineQueryResultFileOptionals contains optional fields that all inline
// query file-like results support.
type InlineQueryResultFileOptionals struct {
	Caption string `json:"caption,omitempty"`      // Caption of the file to be sent, for limitations check the API documentation (optional).
	Text    string `json:"message_text,omitempty"` // Text of a message to be sent instead of the file, for limitations check the API documentation (optional).
}

// InlineQueryResultPhoto represents a link to a photo.
type InlineQueryResultPhoto struct {
	InlineQueryResultBase
	PhotoURL    string `json:"photo_url"`              // Valid URL of the photo.
	ThumbURL    string `json:"thumb_url"`              // URL of the thumbnail for the photo.
	PhotoWidth  int    `json:"photo_width,omitempty"`  // Width of the photo (optional).
	PhotoHeight int    `json:"photo_height,omitempty"` // Height of the photo (optional).
	Title       string `json:"title,omitempty"`        // Title for the result (optional).
	Description string `json:"description,omitempty"`  // Description of the result (optional).
	InlineQueryResultFileOptionals
}

// NewInlineQueryResultPhoto returns a new InlineQueryResultPhoto with all
// mandatory fields set.
func NewInlineQueryResultPhoto(id, photoURL, thumbURL string) *InlineQueryResultPhoto {
	return &InlineQueryResultPhoto{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: PhotoResult,
			ID:   id,
		},
		PhotoURL: photoURL,
		ThumbURL: thumbURL,
	}
}

// InlineQueryResultGif represents a link to an animated GIF file.
type InlineQueryResultGif struct {
	InlineQueryResultBase
	GifURL    string `json:"gif_url"`              // Valid URL for the GIF file.
	ThumbURL  string `json:"thumb_url"`            // URL of the static thumbnail for the result.
	GifWidth  int    `json:"gif_width,omitempty"`  // Width of the GIF (optional).
	GifHeight int    `json:"gif_height,omitempty"` // Height of the GIF (optional).
	Title     string `json:"title,omitempty"`      // Title for the result (optional).
	InlineQueryResultFileOptionals
}

// NewInlineQueryResultGif returns a new InlineQueryResultGif with all
// mandatory fields set.
func NewInlineQueryResultGif(id, gifURL, thumbURL string) *InlineQueryResultGif {
	return &InlineQueryResultGif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: GifResult,
			ID:   id,
		},
		GifURL:   gifURL,
		ThumbURL: thumbURL,
	}
}

// InlineQueryResultMpeg4Gif represents a link to a video animation
// (without sound).
type InlineQueryResultMpeg4Gif struct {
	InlineQueryResultBase
	Mpeg4URL    string `json:"mpeg4_url"`              // Valid URL for the MP4 file.
	ThumbURL    string `json:"thumb_url"`              // URL of the static thumbnail for the result.
	Mpeg4Width  int    `json:"mpeg4_width,omitempty"`  // Video width (optional).
	Mpeg4Height int    `json:"mpeg4_height,omitempty"` // Video height (optional).
	Title       string `json:"title,omitempty"`        // Title for the result (optional).
	InlineQueryResultFileOptionals
}

// NewInlineQueryResultMpeg4Gif returns a new InlineQueryResultMpeg4Gif
// with all mandatory fields set.
func NewInlineQueryResultMpeg4Gif(id, mpeg4URL, thumbURL string) *InlineQueryResultMpeg4Gif {
	return &InlineQueryResultMpeg4Gif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: PhotoResult,
			ID:   id,
		},
		Mpeg4URL: mpeg4URL,
		ThumbURL: thumbURL,
	}
}

// MIMEType represents a MIME type.
type MIMEType string

// MIME type constants for an InlineQueryResultVideo.
const (
	MIMETextHTML = MIMEType("text/html")
	MIMEVideoMP4 = MIMEType("video/mp4")
)

// InlineQueryResultVideo represents a link to a video player/file.
type InlineQueryResultVideo struct {
	InlineQueryResultBase
	VideoURL      string   `json:"video_url"`                // Valid URL for the video player/file.
	MIMEType      MIMEType `json:"mime_type"`                // MIME type of the content of the URL.
	ThumbURL      string   `json:"thumb_url"`                // URL of the static thumbnail.
	Title         string   `json:"title"`                    // Title for the result.
	Text          string   `json:"message_text"`             // Text of the message to be sent with the video, for limitations check the API documentation.
	VideoWidth    int      `json:"video_width,omitempty"`    // Video width (optional).
	VideoHeight   int      `json:"video_height,omitempty"`   // Video height (optional).
	VideoDuration int      `json:"video_duration,omitempty"` // Video duration in seconds (optional).
	Description   string   `json:"description"`              // Short description of the result (optional).
	InlineQueryResultFileOptionals
}

// NewInlineQueryResultVideo returns a new InlineQueryResultVideo with all
// mandatory fields set.
func NewInlineQueryResultVideo(id, videoURL, thumbURL, title, text string, mimeType MIMEType) *InlineQueryResultVideo {
	return &InlineQueryResultVideo{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: PhotoResult,
			ID:   id,
		},
		VideoURL: videoURL,
		MIMEType: mimeType,
		ThumbURL: thumbURL,
		Title:    title,
		Text:     text,
	}
}
