package config

const (
	LogDebugLevel = 0                // 日志等级 : 调试信息
	LogInfoLevel  = 1                // 日志等级 : 普通信息
	LogWarnLevel  = 2                // 日志等级 : 警告信息
	LogErrorLevel = 3                // 日志等级 : 错误信息
	LogFatalLevel = 4                // 日志等级 : 致命信息
	LogMaxLevel   = 5                // 日志最大等级
	LogLimitLevel = LogInfoLevel     // 显示这个等级之上的日志(控制台)
	LogBuffMax    = 20 * 1024 * 1024 // 日志缓冲
)

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
)

const (
	UpdateCurrTimeCount = 32 // 刷新时间戳变更上线
)

const (
	LogBuffSize = 10 * 1024 * 1024
	LogDir      = "../log"
	ProfFile    = "pangu_prof.log"
	LogFileName = LogDir + "/pangu.log"
)
