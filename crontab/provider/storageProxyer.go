package provider

import (
	"gogs.xlh/tools/util/crontab/job"
	"time"
)

type StorageProxyer interface {
	ReadJobs() (jbs []job.BasicJob, err error)
	WriteJobs(jbs []job.BasicJob) error
	WriteAJob(j job.BasicJob) error
	UpdateAJob(j job.BasicJob) error
	DeleteAJob(j job.BasicJob) error
	DeleteAJobById(jId string, t time.Time, repeat bool) error
}
