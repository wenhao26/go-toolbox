package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 常量与业务配置定义
const (
	RegisterDayThreshold = 5
	ReadingBookThreshold = 3
	Day5InSeconds        = 5 * 24 * 3600
	FetchLimit           = 3000

	// Dsn MySQL 数据源配置信息
	Dsn = "root:123456@tcp(192.168.33.10:3306)/nbnovel_test?charset=utf8mb4&parseTime=True&loc=Local"

	// WorkerPoolSize 控制并发协程数量，严禁无限并发
	// 针对数据库迁移这类 I/O 密集兼带行锁的操作，建议设置 3~8
	WorkerPoolSize = 5
)

// UserModel 用户模型
type UserModel struct {
	ID         int64
	UUID       string
	CreateTime int64
}

// Migrator 迁移处理器
type Migrator struct {
	db *sql.DB
}

// Run 处理业务
func (m *Migrator) Run(ctx context.Context) error {
	// 预载拉取满足条件的用户
	users, err := m.fetchTargetUsers(ctx)
	if err != nil {
		return fmt.Errorf("拉取用户失败: %w", err)
	}

	totalItems := len(users)
	if totalItems == 0 {
		fmt.Println("暂无更多待处理的数据~")
		return nil
	}

	fmt.Printf("👉 需要检查用户数量为[%d]，数据已准备，即将执行\n", totalItems)
	fmt.Printf("  - 注册时间距离当前时间 ≤%d 天\n", RegisterDayThreshold)
	fmt.Printf("  - 用户阅读书籍数量 ≤%d 本且最近一次阅读章节时间 ≤%d 秒\n", ReadingBookThreshold, Day5InSeconds)
	fmt.Printf("  - 用户是否存在订单历史\n")
	time.Sleep(2 * time.Second)

	// 管道 与 WorkerPool 协程控制并发初始化
	userChan := make(chan UserModel, totalItems)
	var wg sync.WaitGroup

	// 将用户投递到并发通道中
	for _, user := range users {
		userChan <- user
	}
	close(userChan)

	// 启动受限的协程池处理具体业务
	for i := 0; i < WorkerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			m.workerProcessor(ctx, workerID, userChan)
		}(i)
	}

	wg.Wait()
	return nil
}

// fetchTargetUsers 获取目标用户
func (m *Migrator) fetchTargetUsers(ctx context.Context) ([]UserModel, error) {
	query := "SELECT id,uuid,create_time FROM nb_user WHERE user_segment = 3 ORDER BY id ASC LIMIT ?"
	rows, err := m.db.QueryContext(ctx, query, FetchLimit)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	var users []UserModel

	for rows.Next() {
		var user UserModel

		if err := rows.Scan(&user.ID, &user.UUID, &user.CreateTime); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// workerProcessor Worker 处理器，独立的协程流
func (m *Migrator) workerProcessor(ctx context.Context, workerID int, userChan chan UserModel) {
	now := time.Now().Unix()

	for user := range userChan {
		// 时刻拦截上下文状态，一旦收到优雅关闭信号，则立刻中断并不再领取新任务
		select {
		case <-ctx.Done():
			log.Printf("[info] Worker [%d] 收到终止信号，安全退出", workerID)
			return
		default:
		}

		userID := user.ID
		fmt.Printf("[Worker-%d] 📑 检查用户[%d]是否满足筛选条件...\n", workerID, userID)

		// 规则一：排除内部广告审核特定前缀用户
		if strings.Contains(user.UUID, "fb-ads-review-a000") {
			fmt.Printf("[Worker-%d] 👤 用户UUID[%s]属于内部FB审核用户池成员，跳过！\n\n", workerID, user.UUID)
			continue
		}

		// 规则二：计算注册周期阻断条件
		registrationIntervalDays := float64((now - user.CreateTime) / 86400)
		registrationIntervalDays = math.Floor(registrationIntervalDays)
		isRegistrationPeriodMet := registrationIntervalDays <= float64(RegisterDayThreshold)

		// 规则三: 下单记录条件拦截
		hasOrderRecord, err := m.checkHasOrderRecord(ctx, userID)
		if err != nil {
			log.Printf("[error] Worker-%d 用户[%d]订单数据校验异常: %v", workerID, userID, err)
			continue
		}

		// 规则四: 阅读深度限制阻断
		hasValidReading, err := m.checkHasValidReading(ctx, userID, now)
		if err != nil {
			log.Printf("[error] Worker-%d 用户[%d]阅读数据校验异常: %v", workerID, userID, err)
			continue
		}

		fmt.Printf("[Worker-%d] - 是否在注册保护范围：%s\n", workerID, boolToStr(isRegistrationPeriodMet))
		fmt.Printf("[Worker-%d] - 是否存在订单记录：%s\n", workerID, boolToStr(hasOrderRecord))
		fmt.Printf("[Worker-%d] - 是否存在有效阅读记录：%s\n", workerID, boolToStr(hasValidReading))

		// 如果满足任意一项排除保护条件，跳过归档
		if isRegistrationPeriodMet || hasOrderRecord || hasValidReading {
			fmt.Printf("[Worker-%d] 🔔 用户[%d]不符合迁移存档条件，跳过！\n\n", workerID, userID)
			continue
		}

		fmt.Printf("[Worker-%d] 🎯 用户[%d]符合迁移存档条件，进行迁移处理...\n\n", workerID, userID)

		// 符合全部归档条件，启动分布式单人事务迁移
		if err := m.executeMigrationTransaction(ctx, userID); err != nil {
			log.Printf("❌ [Worker-%d] 用户[%d] 迁移失败: %v\n", workerID, userID, err)
		} else {
			fmt.Printf("✅ [Worker-%d] 用户[%d]的相关数据，迁移成功！\n", workerID, userID)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// checkHasOrderRecord 检查是否存在订单记录
func (m *Migrator) checkHasOrderRecord(ctx context.Context, userID int64) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM `nb_order` WHERE `user_id` = ?"
	err := m.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// checkHasValidReading 检查是否存在有效阅读记录
func (m *Migrator) checkHasValidReading(ctx context.Context, userID int64, now int64) (bool, error) {
	query := "SELECT COUNT(`book_id`) as total_books, MAX(`last_read_time`) as max_last_read_time FROM `nb_read_book_history` WHERE `user_id` = ?"

	var totalBooks sql.NullInt64
	var maxLastReadTime sql.NullInt64

	err := m.db.QueryRowContext(ctx, query, userID).Scan(&totalBooks, &maxLastReadTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if !totalBooks.Valid || !maxLastReadTime.Valid {
		return false, nil
	}

	tBooks := totalBooks.Int64
	mLastRead := maxLastReadTime.Int64

	if tBooks <= 0 || mLastRead <= 0 {
		return false, nil
	}

	isConditionMet := tBooks <= int64(ReadingBookThreshold) && (now-mLastRead) > int64(Day5InSeconds)

	return !isConditionMet, nil
}

// executeMigrationTransaction 执行迁移，开启底层 ACID 事务安全读写并销毁
func (m *Migrator) executeMigrationTransaction(ctx context.Context, userID int64) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	// 利用 defer 确保在 panic 或者任何逻辑意外退出时事务一定会被 Rollback
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // 重新抛出 panic 以便在上层被捕获或记录日志
		}
	}()

	// 迁移主体用户表
	userCopySql := "INSERT INTO `nb_user_archive` SELECT * FROM `nb_user` WHERE `id` = ?"
	if _, err := tx.ExecContext(ctx, userCopySql, userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user 归档失败: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM `nb_user` WHERE `id` = ?", userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user 物理删除失败: %w", err)
	}

	// 迁移用户扩展附加表
	extraCopySql := "INSERT INTO `nb_user_extra_archive` SELECT * FROM `nb_user_extra` WHERE `user_id` = ?"
	if _, err := tx.ExecContext(ctx, extraCopySql, userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user_extra 归档失败: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM `nb_user_extra` WHERE `user_id` = ?", userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user_extra 物理删除失败: %w", err)
	}

	// 迁移钱包资产数据表
	walletCopySql := "INSERT INTO `nb_user_wallet_archive` SELECT * FROM `nb_user_wallet` WHERE `user_id` = ?"
	if _, err := tx.ExecContext(ctx, walletCopySql, userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user_wallet 归档失败: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM `nb_user_wallet` WHERE `user_id` = ?", userID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("nb_user_wallet 物理删除失败: %w", err)
	}

	// 统一提交
	return tx.Commit()
}

// elapsedTime 输出消耗时长
func elapsedTime(startTime int64) {
	endTime := time.Now().Unix()
	log.Printf("[info] Elapsed time: %d s\n", endTime-startTime)
}

// boolToStr 输出布尔值转换字符串
func boolToStr(b bool) string {
	if b {
		return "[YES]"
	}
	return "[NO]"
}

func main() {
	// 创建支持取消的上下文，用于捕获系统退出信号实现优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统强制退出信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		<-sigChan

		log.Println("[warn] 接收到系统中止信号，正在优雅关闭未完成任务...")
		cancel()
	}()

	// 初始化数据库连接
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		log.Fatalf("[fatal] 数据库驱动初始化失败: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	// 设置连接池，高并发高性能标准配置
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("[fatal] 数据库断开，无法连通: %v", err)
	}

	fmt.Println("[worker] 启动 Worker 执行任务")

	migrator := &Migrator{db: db}
	startTime := time.Now().Unix()

	// 执行业务
	if err := migrator.Run(ctx); err != nil {
		log.Printf("[error] 脚本异常终止: %v", err)
	}

	elapsedTime(startTime)
}
