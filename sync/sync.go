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
		id,err := crdata.AddFunc(fmt.Sprintf("0/%d * * * * *", repo.Interval), func() {
			fmt.Printf("进行git仓库[%s]的同步\n", repo.Name)
			err := GitService.Sync(repo.Name)
			fmt.Printf("进行git仓库[%s]的同步结束\n", repo.Name)
			if err!=nil{
				fmt.Printf("进行git仓库[%s]的同步发生错误！%v\n", repo.Name, err)
			}
		})
		repoSyncIds[repo.Name] = id
		if err !=nil{
			panic(err)
			return
		}
	}
	crdata.Start()
}