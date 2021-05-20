// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"errors"

	"github.com/CowellTech/git-module-1.1.2"
	api "github.com/CowellTech/go-gogs-client"

	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/route/api/v1/convert"
)

// https://github.com/CowellTech/go-gogs-client/wiki/Repositories#get-branch
func GetBranch(c *context.APIContext) {
	branch, err := c.Repo.Repository.GetBranch(c.Params("*"))
	if err != nil {
		c.NotFoundOrError(err, "get branch")
		return
	}

	commit, err := branch.GetCommit()
	if err != nil {
		c.Error(err, "get commit")
		return
	}

	c.JSONSuccess(convert.ToBranch(branch, commit))
}

// https://github.com/CowellTech/go-gogs-client/wiki/Repositories#list-branches
func ListBranches(c *context.APIContext) {
	branches, err := c.Repo.Repository.GetBranches()
	if err != nil {
		c.Error(err, "get branches")
		return
	}

	apiBranches := make([]*api.Branch, len(branches))
	for i := range branches {
		commit, err := branches[i].GetCommit()
		if err != nil {
			c.Error(err, "get commit")
			return
		}
		apiBranches[i] = convert.ToBranch(branches[i], commit)
	}

	c.JSONSuccess(&apiBranches)
}

func CreateBranch(c *context.APIContext, form api.CreateBranchOption) {
	branchname := form.BranchName

	if git.RepoHasBranch(c.Repo.Repository.RepoPath(), branchname) {
		c.Error(errors.New("ErrBranchExisted"), "Branch is existed")
		return
	}

	gitRepo, err := git.Open(c.Repo.Repository.RepoPath())
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	base := form.Base

	if !git.RepoHasBranch(c.Repo.Repository.RepoPath(), base) {
		c.Error(errors.New("ErrBranchNotFound"), "Base is not existed")
		return
	}

	err = gitRepo.CreateBranch(branchname, base)
	if err != nil {
		c.Error(err, "CreatBranch failed")
		return
	}

	c.JSONSuccess(&form)
}

func DeleteBranch(c *context.APIContext) {
	branchname := c.Params(":name")

	if !git.RepoHasBranch(c.Repo.Repository.RepoPath(), branchname) {
		c.Error(errors.New("ErrBranchNotFound"), "branch is not existed")
		return
	}

	gitRepo, err := git.Open(c.Repo.Repository.RepoPath())

	if err != nil {
		c.Error(err, "open repository")
		return
	}

	err = gitRepo.DeleteBranch(branchname)
	if err != nil {
		c.Error(err, "branch deleted failed")
		return
	}

	c.NoContent()
}
