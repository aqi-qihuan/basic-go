package ioc

import (
	"basic-go/lmbook/internal/job"
	"basic-go/lmbook/internal/service"
	"basic-go/lmbook/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService, client *rlock.Client, l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, l, client, time.Second*30)
}

func InitJobs(l logger.LoggerV1, rjob *job.RankingJob) *cron.Cron {
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "lmbook",
		Name:      "cron_job",
		Help:      "定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1s", builder.Build(rjob))
	if err != nil {
		panic(err)
	}
	return expr
}
