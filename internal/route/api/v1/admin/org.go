// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	api "github.com/CowellTech/go-gogs-client"

	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/route/api/v1/org"
	"github.com/CowellTech/gogs-0.12.3/internal/route/api/v1/user"
)

func CreateOrg(c *context.APIContext, form api.CreateOrgOption) {
	org.CreateOrgForUser(c, form, user.GetUserByParams(c))
}
