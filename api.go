// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package tbotapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

// A TelegramBotAPI is an API Client for one Telegram bot.
// Create a new client by calling the New() function.
type TelegramBotAPI struct {
	ID       int            // The bots ID.
	Name     string         // The bots Name as seen by users.
	Username string         // The bots username.
	Updates  chan BotUpdate // A channel providing updates this bot receives.
	baseURIs map[method]string
	closed   chan struct{}
	c        *client
	wg       sync.WaitGroup
}

// BotUpdate represents an update the bot received.
// Always check if an error occurred before using the update.
type BotUpdate struct {
	update Update
	err    error
}

// Update returns the contained update.
func (u *BotUpdate) Update() Update {
	return u.update
}

// Error returns != nil, if an error occurred during retrieval of the
// update.
func (u *BotUpdate) Error() error {
	return u.err
}

const apiBaseURI string = "https://api.telegram.org/bot%s"

// New creates a new API Client for a Telegram bot using the apiKey
// provided.
// It will call the GetMe method to retrieve the bots id, name and
// username.
//
// This bot uses long polling to retrieve its updates. If a webhook was set
// for the given apiKey, this will remove it.
func New(apiKey string) (*TelegramBotAPI, error) {
	toReturn := TelegramBotAPI{
		Updates:  make(chan BotUpdate),
		baseURIs: createEndpoints(fmt.Sprintf(apiBaseURI, apiKey)),
		closed:   make(chan struct{}),
		c:        newClient(fmt.Sprintf(apiBaseURI, apiKey)),
	}
	user, err := toReturn.GetMe()
	if err != nil {
		return nil, err
	}
	toReturn.ID = user.User.ID
	toReturn.Name = user.User.FirstName
	toReturn.Username = *user.User.Username

	err = toReturn.removeWebhook()
	if err != nil {
		return nil, err
	}

	toReturn.wg.Add(1)
	go toReturn.updateLoop()

	return &toReturn, nil
}

// NewWithWebhook creates a new API client for a Telegram bot using the apiKey
// provided. It will call the GetMe method to retrieve the bots id, name and
// username.
// In addition to the API client, a http.HandlerFunc will be returned. This
// handler func reacts to webhook requests and will put updates into the
// Updates channel.
func NewWithWebhook(apiKey, webhookURL, certificate string) (*TelegramBotAPI, http.HandlerFunc, error) {
	toReturn := TelegramBotAPI{
		Updates:  make(chan BotUpdate),
		baseURIs: createEndpoints(fmt.Sprintf(apiBaseURI, apiKey)),
		closed:   make(chan struct{}),
		c:        newClient(fmt.Sprintf(apiBaseURI, apiKey)),
	}
	user, err := toReturn.GetMe()
	if err != nil {
		return nil, nil, err
	}
	toReturn.ID = user.User.ID
	toReturn.Name = user.User.FirstName
	toReturn.Username = *user.User.Username

	file, err := os.Open(certificate)
	if err != nil {
		return nil, nil, err
	}

	err = toReturn.setWebhook(webhookURL, certificate, file)
	if err != nil {
		return nil, nil, err
	}

	updateFunc := func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			toReturn.Updates <- BotUpdate{err: err}
			return
		}

		update := &Update{}
		err = json.Unmarshal(bytes, update)
		if err != nil {
			toReturn.Updates <- BotUpdate{err: err}
			return
		}

		toReturn.Updates <- BotUpdate{update: *update}
	}

	return &toReturn, updateFunc, nil
}

// Close shuts down this client.
// Until Close returns, new updates and errors will be put into the
// respective channels.
// Note that, if no updates are received, this function may block for up to
// one minute, which is the time interval
// for long polling.
func (api *TelegramBotAPI) Close() {
	select {
	case <-api.closed:
		return
	default:
	}
	close(api.closed)
	api.wg.Wait()
}

func (api *TelegramBotAPI) updateLoop() {
	defer api.wg.Done()
	updates, err := api.getUpdates()
	offset := -1

	for {
		select {
		case <-api.closed:
			return
		default:
		}

		if err != nil {
			api.Updates <- BotUpdate{err: err}
		} else {
			updates.sort()
			for _, update := range updates.Update {
				api.Updates <- BotUpdate{update: update}
				offset = update.ID
			}
		}

		if offset == -1 {
			updates, err = api.getUpdates()
		} else {
			updates, err = api.getUpdatesByOffset(offset + 1)
		}
	}
}

func (api *TelegramBotAPI) getUpdates() (*updateResponse, error) {
	resp := &updateResponse{}
	response, err := api.c.getQuerystring(getUpdates, resp, map[string]string{"timeout": fmt.Sprint(60)})

	if err != nil {
		if response != nil {
			if response.StatusCode() < 500 {
				return nil, err
			}
			//Telegram server problems, retry later...
			time.Sleep(time.Duration(5) * time.Second)
			return api.getUpdates()
		}
		return nil, err
	}
	err = check(&resp.baseResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *TelegramBotAPI) getUpdatesByOffset(offset int) (*updateResponse, error) {
	resp := &updateResponse{}
	response, err := api.c.getQuerystring(getUpdates, resp, map[string]string{
		"timeout": fmt.Sprint(60),
		"offset":  fmt.Sprint(offset),
	})

	if err != nil {
		if response != nil {
			if response.StatusCode() < 500 {
				return nil, err
			}
			//Telegram server problems, retry later...
			time.Sleep(time.Duration(5) * time.Second)
			return api.getUpdatesByOffset(offset)
		}
		return nil, err
	}
	err = check(&resp.baseResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *TelegramBotAPI) setWebhook(url, fileName string, r io.Reader) error {
	req := outgoingSetWebhook{
		URL: url,
		outgoingFileBase: outgoingFileBase{
			fileName: fileName,
			r:        r,
		},
	}
	resp := &baseResponse{}

	_, err := api.c.uploadFile(setWebhook, resp, file{fieldName: "certificate", fileName: req.fileName, r: req.r}, &req)
	if err != nil {
		return err
	}

	return check(resp)
}

func (api *TelegramBotAPI) removeWebhook() error {
	req := outgoingSetWebhook{
		URL: "",
	}
	resp := &baseResponse{}

	_, err := api.c.postJSON(setWebhook, resp, req)
	if err != nil {
		return err
	}

	return check(resp)
}

// GetMe returns basic information about the bot in form of a UserResponse.
func (api *TelegramBotAPI) GetMe() (*UserResponse, error) {
	resp := &UserResponse{}
	_, err := api.c.get(getMe, resp)

	if err != nil {
		return nil, err
	}
	err = check(&resp.baseResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetFile returns a FileResponse containing a Path string needed to
// download a file.
// You will have to construct the download link manually like
// https://api.telegram.org/file/bot<token>/<file_path>, where <file_path>
// is taken from the response.
func (api *TelegramBotAPI) GetFile(fileID string) (*FileResponse, error) {
	resp := &FileResponse{}
	_, err := api.c.getQuerystring(getFile, resp, map[string]string{"file_id": fileID})

	if err != nil {
		return nil, err
	}
	err = check(&resp.baseResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func check(br *baseResponse) error {
	if br.Ok {
		return nil
	}

	return fmt.Errorf("tbotapi: API error: %d - %s", br.ErrorCode, br.Description)
}

// ErrNoFileSpecified is returned in case neither a file name + io.Reader
// nor a fileID were specified.
var ErrNoFileSpecified = errors.New("tbotapi: Neither a fileID nor a fileName/reader were specified")

func (api *TelegramBotAPI) send(s sendable) (resp *MessageResponse, err error) {
	resp = &MessageResponse{}

	switch s := s.(type) {
	case *OutgoingMessage:
		_, err = api.c.postJSON(sendMessage, resp, s)
	case *OutgoingLocation:
		_, err = api.c.postJSON(sendLocation, resp, s)
	case *OutgoingVenue:
		_, err = api.c.postJSON(sendVenue, resp, s)
	case *OutgoingForward:
		_, err = api.c.postJSON(forwardMessage, resp, s)
	case *OutgoingVideo:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendVideo, resp, file{fieldName: "video", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingVideo
				Video string `json:"video"`
			}{
				OutgoingVideo: *s,
				Video:         s.fileID,
			}
			_, err = api.c.postJSON(sendVideo, resp, toSend)
		}
	case *OutgoingPhoto:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendPhoto, resp, file{fieldName: "photo", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingPhoto
				Photo string `json:"photo"`
			}{
				OutgoingPhoto: *s,
				Photo:         s.fileID,
			}
			_, err = api.c.postJSON(sendPhoto, resp, toSend)
		}
	case *OutgoingVoice:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendVoice, resp, file{fieldName: "audio", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingVoice
				Audio string `json:"audio"`
			}{
				OutgoingVoice: *s,
				Audio:         s.fileID,
			}
			_, err = api.c.postJSON(sendVoice, resp, toSend)
		}
	case *OutgoingAudio:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendAudio, resp, file{fieldName: "audio", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingAudio
				Audio string `json:"audio"`
			}{
				OutgoingAudio: *s,
				Audio:         s.fileID,
			}
			_, err = api.c.postJSON(sendAudio, resp, toSend)
		}
	case *OutgoingDocument:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendDocument, resp, file{fieldName: "document", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingDocument
				Document string `json:"document"`
			}{
				OutgoingDocument: *s,
				Document:         s.fileID,
			}
			_, err = api.c.postJSON(sendDocument, resp, toSend)
		}
	case *OutgoingSticker:
		if !s.valid() {
			return nil, ErrNoFileSpecified
		}
		if s.isUpload() {
			_, err = api.c.uploadFile(sendSticker, resp, file{fieldName: "sticker", fileName: s.fileName, r: s.r}, s)
		} else {
			toSend := struct {
				OutgoingSticker
				Sticker string `json:"sticker"`
			}{
				OutgoingSticker: *s,
				Sticker:         s.fileID,
			}
			_, err = api.c.postJSON(sendSticker, resp, toSend)
		}
	default:
		panic(fmt.Sprintf("tbotapi: internal: unexpected type for send(): %T", s))
	}

	if err != nil {
		return nil, err
	}
	err = check(&resp.baseResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
