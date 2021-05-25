// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"encoding/json"
	"errors"
	"fmt"

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

	baseCommit, err := gitRepo.BranchCommit(base)
	if err != nil {
		c.Error(errors.New("ErrGetCommit"), "git cat-file failed")
		return
	}

	res := struct {
		Name   string      `json:"name"`
		Commit *git.Commit `json:"commit"`
	}{
		Name:   branchname,
		Commit: baseCommit,
	}

	c.JSONSuccess(&res)
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

func DiffBranch(c *context.Context) {
	userName := c.Repo.Owner.Name
	repoName := c.Repo.Repository.Name
	c.Repo.Owner.Name = userName
	c.Repo.Repository.Name = repoName
	branch1 := c.Params(":branch1")
	branch2 := c.Params(":branch2")
	if c.Repo.GitRepo == nil {
		// repoPath := models.RepoPath(c.Repo.Owner.Name, c.Repo.Repository.Name)
		var err error
		c.Repo.GitRepo, err = git.Open(c.Repo.Repository.RepoPath())
		if err != nil {
			res := git.DiffBranchInfo{
				Branch1: branch1,
				Branch2: branch2,
				Error: fmt.Sprintf("RepoRef Invalid repo ,请确认 项目:%s/%s 分支:%s和%s是否存在",
					userName, repoName, branch1, branch2),
			}
			c.JSON(200, res)
			return
		}
	}
	res, err := c.Repo.GitRepo.DiffBranch(branch1, branch2)
	if err != nil {
		res.Branch1 = branch1
		res.Branch2 = branch2
		res.Error = fmt.Sprintf("获取分支diff信息失败，请确项目:%s/%s 分支:%s和%s是否存在",
			userName, repoName, branch1, branch2)
	}
	res.Repo = repoName
	res.Owner = userName
	c.JSON(200, res)
}

func DiffBranchList(c *context.Context) {
	var (
		branchList []api.ProjectBranch
		diffList   []git.DiffBranchInfo
	)
	requestBody := c.Query("branchList")
	if len(requestBody) < 1 {
		requestBody, _ = c.Context.Req.Body().String()
	}

	err := json.Unmarshal([]byte(requestBody), &branchList)
	if err != nil {
		diffList = append(diffList, git.DiffBranchInfo{
			Error: "json解析失败，请效验参数",
		})
		c.JSON(200, diffList)
		return
	}

	for _, v := range branchList {
		//if c.Repo.GitRepo == nil { //todo 优化测试处理下
		// repoPath := models.RepoPath(v.Owner, v.Repo)
		var err error
		c.Repo.GitRepo, err = git.Open(c.Repo.Repository.RepoPath())
		fmt.Println("OpenRepository errr ", err)
		if err != nil {
			res := git.DiffBranchInfo{
				Branch1: v.Branch1,
				Branch2: v.Branch2,
				Error: fmt.Sprintf("RepoRef Invalid repo ,请确认 项目:%s/%s 分支:%s和%s是否存在",
					v.Owner, v.Repo, v.Branch1, v.Branch2),
			}
			diffList = append(diffList, res)
			continue
		}
		//}
		res, err := c.Repo.GitRepo.DiffBranch(v.Branch1, v.Branch2)
		fmt.Println("DiffBranch errr ", err)
		fmt.Println("DiffBranch errr ,res = ", res)
		if err != nil {
			res.Branch1 = v.Branch1
			res.Branch2 = v.Branch2
			res.Error = fmt.Sprintf("获取分支diff信息失败，请确项目:%s/%s 分支:%s和%s是否存在",
				v.Owner, v.Repo, v.Branch1, v.Branch2)
		}
		res.Owner = v.Owner
		res.Repo = v.Repo
		diffList = append(diffList, res)
	}

	c.JSON(200, diffList)
}
