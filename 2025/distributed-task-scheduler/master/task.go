package master

// - 定义任务结构

// Task 任务信息
type Task struct {
	ID      string // 任务的唯一标识符
	Content string // 任务的内容(表达式或文本处理指令)
	Result  string // 任务的执行结果
	Status  string // 任务的状态，有效值：待定=pending | 处理中=processing | 完成=completed | 失败=failed
}
