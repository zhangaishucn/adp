package interfaces

import (
	"context"
	"database/sql"
)

// JobAccess 定义 job 访问接口
//
//go:generate mockgen -source ../interfaces/job_access.go -destination ../interfaces/mock/mock_job_access.go
type JobAccess interface {
	CreateJob(ctx context.Context, tx *sql.Tx, jobInfo *JobInfo) error
	DeleteJobs(ctx context.Context, tx *sql.Tx, jobIDs []string) error
	DeleteTasks(ctx context.Context, tx *sql.Tx, jobIDs []string) error
	ListJobs(ctx context.Context, queryParams JobsQueryParams) ([]*JobInfo, error)
	GetJobsTotal(ctx context.Context, queryParams JobsQueryParams) (int64, error)
	ListTasks(ctx context.Context, queryParams TasksQueryParams) ([]*TaskInfo, error)
	GetTasksTotal(ctx context.Context, queryParams TasksQueryParams) (int64, error)

	CreateTasks(ctx context.Context, tx *sql.Tx, tasks map[string]*TaskInfo) error
	GetJob(ctx context.Context, jobID string) (*JobInfo, error)
	GetJobs(ctx context.Context, jobIDs []string) (map[string]*JobInfo, error)

	UpdateJobState(ctx context.Context, tx *sql.Tx, jobID string, stateInfo JobStateInfo) error
	UpdateTaskState(ctx context.Context, taskID string, stateInfo TaskStateInfo) error
}
