package sync

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func init()  {
	initConfig()
	initSync()
}

var repoSyncIds = make(map[string]cron.EntryID)



func initSync(){
	if len(config.Repos) == 0 {
		return
	}
	crdata := cron.New(cron.WithSeconds())
	for _, repo := range config.Repos {
		id,err := crdata.AddJob(repo.Cron, NewSyncJob(repo))
		repoSyncIds[repo.Name] = id
		if err !=nil{
			panic(err)
			return
		}
	}
	crdata.Start()
}


type SyncJob struct{
	Repo Repo
}

func NewSyncJob(repo Repo)SyncJob{
	return SyncJob{
		Repo: repo,
	}
}

func (job SyncJob)Run(){
	fmt.Printf("进行git仓库[%s]的同步\n", job.Repo.Name)
	err := GitService.Sync(job.Repo.Name)
	fmt.Printf("进行git仓库[%s]的同步结束\n", job.Repo.Name)
	if err!=nil{
		fmt.Printf("进行git仓库[%s]的同步发生错误！%v\n", job.Repo.Name, err)
	}
}
