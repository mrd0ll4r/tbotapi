// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

// Package tbotapi provides a Go wrapper for the Telegram Messenger Bot API.
//
// Note that, if the REST API returns an error, that error will be wrapped in a Go error.
//
// We currently only support long polling (i.e. no webhooks). Feature-wise, everything up to and including the January 20
// changes should be implemented.
//
// Examples are provided in the examples package, so check that out.
//
// The Bot API imposes certain limitations, these are especially interesting for inline query results and files. This
// library does not keep track of those limitations, so you'll have to perform checks yourself.
package tbotapi
