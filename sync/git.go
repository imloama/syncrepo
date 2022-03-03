package sync

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	cfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path"
)


var GitService = new(gitService)

type gitService struct{ }

type GitRepo struct {

}

func(svr *gitService)Clone(repoName string, override bool)error{
	repo := repos[repoName]
	if repo == nil{
		return errors.New(fmt.Sprintf("repo %s 不存在！", repoName))
	}
	folder := path.Join(config.Folder, repoName)
	_, err := os.Stat(folder)
	if err == nil { // 文件夹存在，
		// 判断是否为正常的git仓库
		_, err = git.PlainOpen(folder)
		if err == nil{
			if !override {
				fmt.Printf("文件夹【%s】已经存在，且设置不覆盖，直接返回！\n", folder)
				return nil
			}
		}
		// 删除后重新下载
		err = os.RemoveAll(folder)
		if err!=nil{
			return errors.New(fmt.Sprintf("删除文件夹[%s]失败！%v", folder, err))
		}
	}
	// 文件夹不存在，
	err = os.MkdirAll(folder, os.ModeDir)
	if err!=nil{
		return errors.New(fmt.Sprintf("创建文件夹[%s]失败！%v", folder, err))
	}
	gitrepo, err := git.PlainClone(folder, false, &git.CloneOptions{
		URL: repo.From,
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
		SingleBranch: true,
	})
	if err!=nil{
		return errors.New(fmt.Sprintf("%v", err))
	}
	_, err = gitrepo.CreateRemote(&cfg.RemoteConfig{
		Name: "target",
		URLs: []string{ repo.Target },
	})
	if err!=nil{
		return errors.New(fmt.Sprintf("创建待同步的地址失败！%v", err))
	}
	return nil
}

// 同步操作
func(svr *gitService)Sync(repoName string)error{
	folder := path.Join(config.Folder, repoName)
	err := svr.Clone(repoName, false)
	if err!=nil{
		return err
	}
	repo, err := git.PlainOpen(folder)
	if err != nil {
		fmt.Printf("文件夹【%s】打开错误！%v\n", folder, err)
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		fmt.Printf("文件夹【%s】打开worktree错误！%v\n", folder, err)
		return err
	}

	err = w.Pull(&git.PullOptions{ RemoteName: "origin" })
	if err!=nil{
		fmt.Printf("文件夹【%s】pull错误！%v\n", folder, err)
		return err
	}
	w.Pull(&git.PullOptions{ RemoteName: "target" })

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
	})
	err = repo.Push(&git.PushOptions{
		RemoteName: "target",
	})
	if err!=nil{
		fmt.Printf("文件夹【%s】push错误！%v\n", folder, err)
		return err
	}

	return nil
}