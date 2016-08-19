// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := "123456789:Your_API_token_goes_here"

	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			if typ.IsChatAction() {
				// Ignore chat actions, like group creates or joins.
				fmt.Println("Ignoring chat action")
				return
			}
			if msg.Chat.IsChannel() {
				// Ignore messages _about_ channels (bots cannot receive from channels, but will probably get events like channel creations).
				fmt.Println("Ignoring channel message")
				return
			}

			// Display the incoming message.
			// msg.Chat implements fmt.Stringer, so it'll display nicely.
			// MessageType implements fmt.Stringer, so it'll display nicely.
			fmt.Printf("<-%d, From:\t%s, Type: %s \n", msg.ID, msg.Chat, typ)

			// Create a message with some keyboard markup.
			toSend := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), "What time is it where I am?")
			toSend.SetReplyKeyboardMarkup(tbotapi.ReplyKeyboardMarkup{
				Keyboard:        [][]tbotapi.KeyboardButton{[]tbotapi.KeyboardButton{{Text: time.Now().Format(time.RFC1123Z)}}},
				OneTimeKeyboard: true,
			})

			// Send it.
			outMsg, err := toSend.Send()

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
	boilerplate.RunBot(apiToken, updateFunc, "KeyboardMarkup", "Demonstrates keyboard markup")
}
