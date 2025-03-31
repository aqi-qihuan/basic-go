package domain

type AsyncSMS struct {
	Id      int64
	TplId   string
	Args    []string
	Numbers []string
	// 重试的配置
	RetryMax int
}
