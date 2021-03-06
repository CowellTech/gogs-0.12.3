// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/CowellTech/git-module-1.1.2"
	api "github.com/CowellTech/go-gogs-client"

	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/db"
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
	repoPath := c.Repo.Repository.RepoPath()

	if git.RepoHasBranch(repoPath, branchname) {
		c.Error(errors.New("ErrBranchExisted"), "Branch is existed")
		return
	}

	gitRepo, err := git.Open(repoPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	base := form.Base

	if !git.RepoHasBranch(repoPath, base) {
		c.Error(errors.New("ErrBranchNotFound"), "Base is not existed")
		return
	}

	err = gitRepo.CreateBranch(branchname, base)
	if err != nil {
		c.Error(err, "CreatBranch failed")
		return
	}

	branch := &db.Branch{
		Name:     branchname,
		RepoPath: repoPath,
	}
	// baseCommitID, err := gitRepo.BranchCommitID(base)
	// if err != nil {
	// 	c.Error(errors.New("ErrGitShowRef"), "git show-ref failed")
	// 	return
	// }

	// baseCommit, err := gitRepo.CommitByRevision(baseCommitID)
	// if err != nil {
	// 	c.Error(errors.New("ErrRevisionNotExist"), "bad revision")
	// 	return
	// }
	commit, err := branch.GetCommit()
	if err != nil {
		c.Error(err, "get commit")
		return
	}

	// res := struct {
	// 	Name   string `json:"name"`
	// 	Commit string `json:"commit"`
	// 	Msg    string `json:"message"`
	// }{
	// 	Name:   branchname,
	// 	Commit: baseCommitID,
	// 	Msg:    baseCommit.Message,
	// }

	// c.JSONSuccess(&res)
	c.JSONSuccess(convert.ToBranch(branch, commit))
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
				Error: fmt.Sprintf("RepoRef Invalid repo ,????????? ??????:%s/%s ??????:%s???%s????????????",
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
		res.Error = fmt.Sprintf("????????????diff???????????????????????????:%s/%s ??????:%s???%s????????????",
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
			Error: "json??????????????????????????????",
		})
		c.JSON(200, diffList)
		return
	}

	for _, v := range branchList {
		//if c.Repo.GitRepo == nil { //todo ?????????????????????
		// repoPath := models.RepoPath(v.Owner, v.Repo)
		var err error
		c.Repo.GitRepo, err = git.Open(c.Repo.Repository.RepoPath())
		fmt.Println("OpenRepository errr ", err)
		if err != nil {
			res := git.DiffBranchInfo{
				Branch1: v.Branch1,
				Branch2: v.Branch2,
				Error: fmt.Sprintf("RepoRef Invalid repo ,????????? ??????:%s/%s ??????:%s???%s????????????",
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
			res.Error = fmt.Sprintf("????????????diff???????????????????????????:%s/%s ??????:%s???%s????????????",
				v.Owner, v.Repo, v.Branch1, v.Branch2)
		}
		res.Owner = v.Owner
		res.Repo = v.Repo
		diffList = append(diffList, res)
	}

	c.JSON(200, diffList)
}

func GetCommitsOfBranch(c *context.Context) {
	branch := c.Params(":branch")
	ps := c.Params(":pagesize")
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		c.Error(err, "??????????????????")
		return
	}

	repoPath := c.Repo.Repository.RepoPath()

	type CommitResp struct {
		ID     string         `json:"id"`
		Author *git.Signature `json:"author"`
		// The committer of the commit.
		Committer *git.Signature `json:"committer"`
		// The full commit message.
		Message string `json:"message"`
	}
	gitRepo, err := git.Open(repoPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	if !git.RepoHasBranch(repoPath, branch) {
		c.Error(errors.New("ErrBranchNotFound"), "branch is not existed")
		return
	}

	commits, err := gitRepo.CommitsByPage(branch, 1, pagesize)
	var cr []CommitResp
	for _, v := range commits {
		var c CommitResp
		c.Author = v.Author
		c.Committer = v.Committer
		c.Message = v.Message
		c.ID = v.ID.String()
		cr = append(cr, c)
	}

	if err != nil {
		c.Error(err, "Get Commits failed")
		return
	}

	c.JSON(200, cr)
}
