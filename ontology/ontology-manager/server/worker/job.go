package worker

import (
	"ontology-manager/interfaces"
)

type Job struct {
	mJobInfo     *interfaces.JobInfo
	mTasks       map[string]Task
	mFinishCount int
}

func (j *Job) GetJobInfo() *interfaces.JobInfo {
	return j.mJobInfo
}

type Task interface {
	GetTaskInfo() *interfaces.TaskInfo
}
