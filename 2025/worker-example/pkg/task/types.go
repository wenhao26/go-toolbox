package task

// Task 任务结构体
type Task struct {
	ID      string `json:"id"`      // 任务ID
	Payload string `json:"payload"` // 实际业务载荷
}
