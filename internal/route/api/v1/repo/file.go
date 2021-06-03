// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"fmt"

	"github.com/CowellTech/git-module-1.1.2"
	api "github.com/CowellTech/go-gogs-client"

	"github.com/CowellTech/gogs-0.12.3/internal/context"
	"github.com/CowellTech/gogs-0.12.3/internal/db"
	"github.com/CowellTech/gogs-0.12.3/internal/gitutil"
	"github.com/CowellTech/gogs-0.12.3/internal/route/repo"
)

func GetRawFile(c *context.APIContext) {
	if !c.Repo.HasAccess() {
		c.NotFound()
		return
	}

	if c.Repo.Repository.IsBare {
		c.NotFound()
		return
	}

	blob, err := c.Repo.Commit.Blob(c.Repo.TreePath)
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get blob")
		return
	}
	if err = repo.ServeBlob(c.Context, blob); err != nil {
		c.Error(err, "serve blob")
	}
}

func GetArchive(c *context.APIContext) {
	repoPath := db.RepoPath(c.Params(":username"), c.Params(":reponame"))
	gitRepo, err := git.Open(repoPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}
	c.Repo.GitRepo = gitRepo

	repo.Download(c.Context)
}

func GetEditorconfig(c *context.APIContext) {
	ec, err := c.Repo.Editorconfig()
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get .editorconfig")
		return
	}

	fileName := c.Params("filename")
	def, err := ec.GetDefinitionForFilename(fileName)
	if err != nil {
		c.Error(err, "get definition for filename")
		return
	}
	if def == nil {
		c.NotFound()
		return
	}
	c.JSONSuccess(def)
}

func GetRawFiles(c *context.APIContext, fileList []api.DiffFileList) {
	var (
		// fileList       []api.DiffFileList
		retrunFileList []api.ReturnDiffFile
	)
	// requestBody := c.Query("fileList")
	// if len(requestBody) < 1 {
	// 	requestBody, _ = c.Context.Req.Body().String()
	// }

	// err := json.Unmarshal([]byte(requestBody), &fileList)
	// if err != nil {
	// 	fmt.Println(err)
	// 	//todo
	// }

	for _, v := range fileList {
		diffFileInfo := api.ReturnDiffFile{
			BaseInfo: v,
		}
		if v.IsBinary == false {
			repoPath := db.RepoPath(v.ProjectOwner, v.Project)
			// repoPath := models.RepoPath(v.ProjectOwner, v.Project)
			//todo 容错判断
			gitrepo, err := git.Open(repoPath)
			if err != nil {
				fmt.Println("GetRawFiles OpenRepository err ", err)
			}
			c.Repo.GitRepo = gitrepo
			c.Repo.TreePath = v.File

			c.Repo.CommitID = v.BaseDiffBranchCommitID
			c.Repo.Commit, err = c.Repo.GitRepo.CommitByRevision(v.BaseDiffBranchCommitID)
			// blob, err := c.Repo.Commit.GetBlobByPath(c.Repo.TreePath)
			blob, err := c.Repo.Commit.Blob(c.Repo.TreePath)
			if err != nil {
				fmt.Println("GetBlobByPath err ", err)
				diffFileInfo.BaseDiffFile = ""
			} else {
				r, err := blob.Blob().Bytes()
				if err != nil {
					fmt.Println("r, err := blob.Blob().Data() ", err)
				}
				diffFileInfo.BaseDiffFile = fmt.Sprintf("%s", r)
			}

			c.Repo.CommitID = v.DeployBranchCommitID
			// c.Repo.Commit, err = c.Repo.GitRepo.GetCommit(v.DeployBranchCommitID)
			c.Repo.Commit, err = c.Repo.GitRepo.CommitByRevision(v.DeployBranchCommitID)
			blob2, _ := c.Repo.Commit.Blob(c.Repo.TreePath)
			// r2, _ := blob2.Blob().Data()
			r2, _ := blob2.Blob().Bytes()

			diffFileInfo.BranchDiffFile = fmt.Sprintf("%s", r2)
		} else {
			diffFileInfo.BaseDiffFile = "[gogs提示]二进制文件暂不提供显示。"
			diffFileInfo.BranchDiffFile = ""
		}
		retrunFileList = append(retrunFileList, diffFileInfo)
	}
	c.JSON(200, retrunFileList)
	return
}
