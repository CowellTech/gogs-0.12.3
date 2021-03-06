// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"github.com/CowellTech/gogs-0.12.3/internal/db"
)

type APIOrganization struct {
	Organization *db.User
	Team         *db.Team
}
