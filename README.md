# TBotAPI - Telegram Bot-API for Go #

This is a simple wrapper for the Telegram Bot-API for Go.
It provides a high-level API in Go that works with the Telegram REST API.

[![Build Status](https://travis-ci.org/mrd0ll4r/tbotapi.svg?branch=master)](https://travis-ci.org/mrd0ll4r/tbotapi)
[![GoDoc](https://godoc.org/github.com/mrd0ll4r/tbotapi?status.svg)](https://godoc.org/github.com/mrd0ll4r/tbotapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrd0ll4r/tbotapi)](https://goreportcard.com/report/github.com/mrd0ll4r/tbotapi)

The implementation is pretty raw, i.e. you will just send and receive messages - you have to handle any command parsing or stuff yourself.

### How do I get set up? ###

A simple

    go get -u github.com/mrd0ll4r/tbotapi

should do it.

### Example ###

See `examples/` for some simple bots.
Start with the echo bot.

### API-stableness ###

Is the API stable? **No**

Why is the API not stable?
Because Telegram changes its bot API frequently.
We try to only add things, not remove them.
So everything you did with this library before should also work after API changes.

**Exception**: I did one major overhaul so far, where I removed the model package and unified most sending methods.
I'm pretty happy with the architecture so far, so I hope I won't have to do that again.

### What do we use? ###

We use

* [resty] for REST calls


[resty]: https://github.com/go-resty/resty

### Contribution guidelines ###

If you want to contribute code, make sure to read the [contributing guidelines].


[contributing guidelines]: https://github.com/mrd0ll4r/tbotapi/blob/master/CONTRIBUTING.md

### Things that need to be done ###

* Implement a test
* Implement the API 2.0 changes
* Fix godoc line length
* Move these TODOs to GitHub issues

### License
This work is licensed under the MIT License. A copy of the MIT License can be found in the `LICENSE` file.

Feel free to use this library for any bot whatsoever.
If you find any bugs, have any ideas about improvements or just want to show me what you've done with this, please contact me through [this bot].


[this bot]: (https://telegram.me/tbotapibot).
