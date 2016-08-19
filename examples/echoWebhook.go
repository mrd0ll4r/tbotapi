// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package main

import (
	"fmt"

	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := "123456789:Your_API_token_goes_here"
	botHost := "your.host.com"
	botPort := uint16(8443)
	privkey := "private.key"
	pubkey := "public.pem"

	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			if typ != tbotapi.TextMessage {
				// Ignore non-text messages for now.
				fmt.Println("Ignoring non-text message")
				return
			}
			// Note: Bots cannot receive from channels, at least no text messages. So we don't have to distinguish anything here.

			// Display the incoming message.
			// msg.Chat implements fmt.Stringer, so it'll display nicely.
			// We know it's a text message, so we can safely use the Message.Text pointer.
			fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, *msg.Text)

			// Now simply echo that back.
			outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), *msg.Text).Send()

			if err != nil {
				fmt.Printf("Error sending: %s\n", err)
				return
			}
			fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
		case tbotapi.InlineQueryUpdate:
			fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
		case tbotapi.ChosenInlineResultUpdate:
			fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
		default:
			fmt.Println("Ignoring unknown Update type.")
		}
	}

	// Run the bot, this will block.
	boilerplate.RunBotOnWebhook(apiToken, updateFunc, "Echo", "Echoes messages back", botHost, botPort, pubkey, privkey)
}
