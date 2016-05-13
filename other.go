// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package tbotapi

import "io"

// Querystring is a type to represent querystring-applicable data
type querystring map[string]string

type querystringer interface {
	querystring() querystring
}

type file struct {
	fieldName string
	fileName  string
	r         io.Reader
}
