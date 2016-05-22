// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package boilerplate

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mrd0ll4r/tbotapi"
)

// BotFunc describes how the bot handles an update.
type BotFunc func(update tbotapi.Update, api *tbotapi.TelegramBotAPI)

// RunBot runs a bot.
// THIS IS JUST FOR DEMONSTRATION! NOT TO BE USED IN PRODUCTION!
// It will block until either something very bad happens or closing is
// closed.
func RunBot(apiKey string, bot BotFunc, name, description string) {
	fmt.Printf("%s: %s\n", name, description)
	fmt.Println("Starting...")

	api, err := tbotapi.New(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Just to show its working.
	fmt.Printf("User ID: %d\n", api.ID)
	fmt.Printf("Bot Name: %s\n", api.Name)
	fmt.Printf("Bot Username: %s\n", api.Username)

	closed := make(chan struct{})
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closed:
				return
			case update := <-api.Updates:
				if update.Error() != nil {
					fmt.Printf("Update error: %s\n", update.Error())
					continue
				}

				bot(update.Update(), api)
			}
		}
	}()

	// Ensure a clean shutdown.
	closing := make(chan struct{})
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		signal.Stop(shutdown)
		close(shutdown)
		close(closing)
	}()

	fmt.Println("Bot started. Press CTRL-C to close...")

	// Wait for the signal.
	<-closing
	fmt.Println("Closing...")

	// Always close the API first, let it clean up the update loop.
	api.Close() // This might take a while.
	close(closed)
	wg.Wait()
}

// RunBotOnWebhook runs the given BotFunc with a webhook.
func RunBotOnWebhook(apiKey string, bot BotFunc, name, description, webhookHost string, webhookPort uint16, pubkey, privkey string) {
	fmt.Printf("%s: %s\n", name, description)
	fmt.Println("Starting...")
	u := url.URL{
		Host:   webhookHost + ":" + fmt.Sprint(webhookPort),
		Scheme: "https",
		Path:   apiKey,
	}

	api, handler, err := tbotapi.NewWithWebhook(apiKey, u.String(), pubkey)
	if err != nil {
		log.Fatal(err)
	}

	// Just to show its working.
	fmt.Printf("User ID: %d\n", api.ID)
	fmt.Printf("Bot Name: %s\n", api.Name)
	fmt.Printf("Bot Username: %s\n", api.Username)

	closed := make(chan struct{})
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closed:
				return
			case update := <-api.Updates:
				if update.Error() != nil {
					fmt.Printf("Update error: %s\n", update.Error())
					continue
				}

				bot(update.Update(), api)
			}
		}
	}()

	http.HandleFunc("/"+apiKey, handler)

	fmt.Println("Starting webhook...")
	go func() {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+fmt.Sprint(webhookPort), pubkey, privkey, nil))
	}()

	// Ensure a clean shutdown.
	closing := make(chan struct{})
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		signal.Stop(shutdown)
		close(shutdown)
		close(closing)
	}()

	fmt.Println("Bot started. Press CTRL-C to close...")

	// Wait for the signal.
	<-closing
	fmt.Println("Closing...")

	// Always close the API first.
	api.Close() // This is instant for webhook based bots.
	close(closed)
	wg.Wait()
}
