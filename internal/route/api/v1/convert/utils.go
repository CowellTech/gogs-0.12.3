// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package convert

import (
	"github.com/CowellTech/gogs-0.12.3/internal/conf"
)

// ToCorrectPageSize makes sure page size is in allowed range.
func ToCorrectPageSize(size int) int {
	if size <= 0 {
		size = 10
	} else if size > conf.API.MaxResponseItems {
		size = conf.API.MaxResponseItems
	}
	return size
}
