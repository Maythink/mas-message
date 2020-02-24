package mas

const (
	ProductionHost = "http://112.35.1.155:1992"
)

const (
	NorURL = "/sms/norsubmit" // 一对一/一对多普通短信(一种短信内容)
	TmpURL = "/sms/tmpsubmit" // 模板短信
)

var (
	PostRetryTimes = 3 //重试次数
)
