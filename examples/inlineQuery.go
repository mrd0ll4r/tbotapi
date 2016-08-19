// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"unicode"

	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := "123456789:Your_API_token_goes_here"

	// Note: For this example to work, you'll have to enable inline queries for your bot (chat with @BotFather).

	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		switch update.Type() {
		case tbotapi.MessageUpdate:
			fmt.Println("Ignoring received message of type:", update.Message.Type().String())
		case tbotapi.InlineQueryUpdate:
			query := update.InlineQuery
			fmt.Printf("<-%s (query), From:\t%s, Query: %s \n", query.ID, query.From, query.Query)
			var results []tbotapi.InlineQueryResult

			for i, s := range query.Query {
				if len(results) >= 50 {
					// The API accepts up to 50 results.
					break
				}
				if !unicode.IsSpace(s) {
					// Don't set mandatory fields to whitespace.
					results = append(results, tbotapi.NewInlineQueryResultArticle(fmt.Sprint(i), string(s), string(s)))
				}
			}

			err := api.NewInlineQueryAnswer(query.ID, results).Send()
			if err != nil {
				fmt.Printf("Err: %s\n", err)
			}
		case tbotapi.ChosenInlineResultUpdate:
			// id, not value.
			fmt.Println("Chosen inline query result (ID):", update.ChosenInlineResult.ID)
		default:
			fmt.Println("Ignoring unknown Update type.")
		}
	}

	// Run the bot, this will block.
	boilerplate.RunBot(apiToken, updateFunc, "InlineQuery", "Demonstrates inline queries by splitting words")
}
