package model

// -- 任务模型定义 --

// Task 任务模型
type Task struct {
	ID          string `json:"id"`          // 唯一标识
	Name        string `json:"name"`        // 任务名称
	Expr        string `json:"expr"`        // Cron表达式，*/5 * * * * *（支持秒级）
	Command     string `json:"command"`     // Shell命令
	Description string `json:"description"` // 任务描述
	Status      int    `json:"status"`      // 任务状态，1 - 运行中，0 - 停止
	Timeout     int    `json:"timeout"`     // 超时时间（秒）
}

// Log 运行日志模型
type Log struct {
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
	Output    string `json:"output"`     // 输出
	Success   bool   `json:"success"`    // 成功标识
}
