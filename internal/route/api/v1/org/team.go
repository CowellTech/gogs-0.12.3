// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package org

import (
	api "github.com/CowellTech/go-gogs-client"

	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/route/api/v1/convert"
)

func ListTeams(c *context.APIContext) {
	org := c.Org.Organization
	if err := org.GetTeams(); err != nil {
		c.Error(err, "get teams")
		return
	}

	apiTeams := make([]*api.Team, len(org.Teams))
	for i := range org.Teams {
		apiTeams[i] = convert.ToTeam(org.Teams[i])
	}
	c.JSONSuccess(apiTeams)
}
