// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"github.com/CowellTech/gogs-0.12.3/internal/conf"
	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/db"
	"github.com/CowellTech/gogs-0.12.3/internal/route"
)

const (
	ORGS = "admin/org/list"
)

func Organizations(c *context.Context) {
	c.Data["Title"] = c.Tr("admin.organizations")
	c.Data["PageIsAdmin"] = true
	c.Data["PageIsAdminOrganizations"] = true

	route.RenderUserSearch(c, &route.UserSearchOptions{
		Type:     db.UserOrganization,
		Counter:  db.CountOrganizations,
		Ranger:   db.Organizations,
		PageSize: conf.UI.Admin.OrgPagingNum,
		OrderBy:  "id ASC",
		TplName:  ORGS,
	})
}
