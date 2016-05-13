// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package tbotapi

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v0"
	"net/http"
)

type method string

const (
	getMe                = method("GetMe")
	sendMessage          = method("SendMessage")
	forwardMessage       = method("ForwardMessage")
	sendPhoto            = method("SendPhoto")
	sendAudio            = method("SendAudio")
	sendDocument         = method("SendDocument")
	sendSticker          = method("SendSticker")
	sendVideo            = method("SendVideo")
	sendVoice            = method("SendVoice")
	sendLocation         = method("SendLocation")
	sendVenue            = method("SendVenue")
	sendChatAction       = method("SendChatAction")
	getUserProfilePhotos = method("GetUserProfilePhotos")
	getUpdates           = method("GetUpdates")
	setWebhook           = method("SetWebhook")
	getFile              = method("GetFile")
	answerInlineQuery    = method("AnswerInlineQuery")
	kickChatMember       = method("KickChatMember")
	unbanChatMember      = method("UnbanChatMember")
	answerCallbackQuery  = method("AnswerCallbackQuery")
)

type client struct {
	c         *resty.Client
	endpoints map[method]string
}

func newClient(baseURI string) *client {
	toReturn := &client{
		c:         resty.New().SetHTTPMode().OnAfterResponse(parseResponseBody).OnAfterResponse(checkHTTPStatus),
		endpoints: createEndpoints(baseURI),
	}

	return toReturn
}

func (c *client) get(m method, result interface{}) (*resty.Response, error) {
	return c.c.R().SetResult(result).Get(c.getEndpoint(m))
}

func (c *client) getQuerystring(m method, result interface{}, querystring map[string]string) (*resty.Response, error) {
	return c.c.R().SetQueryParams(querystring).SetResult(result).Get(c.getEndpoint(m))
}

func (c *client) postJSON(m method, result interface{}, data interface{}) (*resty.Response, error) {
	return c.c.R().SetBody(data).SetResult(result).Post(c.getEndpoint(m))
}

func (c *client) uploadFile(m method, result interface{}, data file, fields querystringer) (*resty.Response, error) {
	return c.c.R().SetFileReader(data.fieldName, data.fileName, data.r).SetResult(result).SetFormData(map[string]string(fields.querystring())).Post(c.getEndpoint(m))
}

func parseResponseBody(c *resty.Client, res *resty.Response) (err error) {
	// Handles only JSON
	ct := res.Header().Get(http.CanonicalHeaderKey("Content-Type"))
	if resty.IsJSONType(ct) {
		// Considered as Result
		if res.StatusCode() > 199 && res.StatusCode() < 500 {
			if res.Request.Result != nil {
				err = json.Unmarshal(res.Body(), res.Request.Result)
			}
		}
	}

	return
}

func checkHTTPStatus(_ *resty.Client, res *resty.Response) error {
	if res.StatusCode() >= 500 {
		return fmt.Errorf("API: Server error: returned %s when requesting %s", res.Status(), res.Request.URL)
	}
	return nil
}

func (c *client) getEndpoint(method method) string {
	endpoint, ok := c.endpoints[method]
	if !ok {
		panic(fmt.Errorf("tbotapi: internal: Endpoint for method %s not found", string(method)))
	}
	return endpoint
}

func createEndpoints(baseURI string) map[method]string {
	toReturn := map[method]string{}

	toReturn[getMe] = fmt.Sprint(baseURI, "/", string(getMe))
	toReturn[sendMessage] = fmt.Sprint(baseURI, "/", string(sendMessage))
	toReturn[forwardMessage] = fmt.Sprint(baseURI, "/", string(forwardMessage))
	toReturn[sendPhoto] = fmt.Sprint(baseURI, "/", string(sendPhoto))
	toReturn[sendAudio] = fmt.Sprint(baseURI, "/", string(sendAudio))
	toReturn[sendDocument] = fmt.Sprint(baseURI, "/", string(sendDocument))
	toReturn[sendSticker] = fmt.Sprint(baseURI, "/", string(sendSticker))
	toReturn[sendVideo] = fmt.Sprint(baseURI, "/", string(sendVideo))
	toReturn[sendVoice] = fmt.Sprint(baseURI, "/", string(sendVoice))
	toReturn[sendLocation] = fmt.Sprint(baseURI, "/", string(sendLocation))
	toReturn[sendChatAction] = fmt.Sprint(baseURI, "/", string(sendChatAction))
	toReturn[getUserProfilePhotos] = fmt.Sprint(baseURI, "/", string(getUserProfilePhotos))
	toReturn[getUpdates] = fmt.Sprint(baseURI, "/", string(getUpdates))
	toReturn[setWebhook] = fmt.Sprint(baseURI, "/", string(setWebhook))
	toReturn[getFile] = fmt.Sprint(baseURI, "/", string(getFile))
	toReturn[answerInlineQuery] = fmt.Sprint(baseURI, "/", string(answerInlineQuery))

	return toReturn
}
