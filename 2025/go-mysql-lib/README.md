# 封装MySQL类库
基于 `github.com/go-sql-driver/mysql` 包封装一个强大、高性能、高并发且具有良好扩展性的生产级 MySQL 类库。

> 需要考虑以下几个核心方面：

- 连接池管理 (Connection Pool Management): 充分利用 database/sql 内置的连接池功能。

- 预处理语句 (Prepared Statements): 这是性能的关键，尤其是对于重复执行的语句。

- 事务管理 (Transaction Management): 保证数据一致性，并能有效提升批量操作性能。

- 上下文 (Context) 支持: 允许取消操作和设置超时，在并发场景下至关重要。

- SQL 构建/抽象: 提供更便捷的 SQL 构建方式，减少手写 SQL 的错误率和复杂度。

- 错误处理与日志 (Error Handling & Logging): 完善的错误封装和日志记录，便于问题排查。

- 扩展性 (Extensibility): 方便未来添加中间件、监控或切换驱动。

# Go MySQL 生产级类库项目结构

为了实现高性能、高并发、良好扩展性和易维护性，我们将类库分解为以下核心文件：

```text
mysql_repo/
├── client.go               # 底层数据库客户端接口及其MySQL实现
├── executor.go             # 核心数据库操作执行器，提供CRUD和事务方法
├── options.go              # 连接池及类库配置选项
├── statement_cache.go      # SQL预处理语句缓存，提高重复查询性能
└── errors.go               # 自定义错误类型，便于错误识别和处理
```

# 使用示例
```go
// User 结构体用于映射数据库表中的用户数据
type User struct {
	ID        int64     // 用户ID，通常是自增主键
	Name      string    // 用户名
	Email     string    // 用户邮箱，通常是唯一索引
	CreatedAt time.Time // 记录创建时间
}

func main() {
	// 替换为你的 MySQL DSN (Data Source Name)。
	// 例如: "username:password@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4"
	// parseTime=true 是非常重要的，它会将 MySQL 的 DATE/DATETIME 类型自动解析为 Go 的 time.Time。
	dsn := "root:123456@tcp(127.0.0.1:3306)/testdb?parseTime=true&charset=utf8mb4"

	// 创建一个 Options 实例，可以根据需要调整参数。
	// 这里使用默认值，但你可以自定义，例如：
	// opts := db_executor.DefaultOptions()
	// opts.MaxOpenConns = 50
	// opts.Logger = log.New(os.Stdout, "MY_APP_DB: ", log.LstdFlags|log.Lshortfile)
	opts := db_executor.DefaultOptions()

	// 初始化数据库执行器。这是你与数据库交互的主要对象。
	executor, err := db_executor.NewMySQLExecutor(dsn, opts)
	if err != nil {
		log.Fatalf("初始化数据库执行器失败: %v", err) // 致命错误，应用程序无法启动
	}
	defer executor.Close() // 确保在 main 函数退出时关闭所有数据库资源

	// 创建一个带有超时的上下文，用于控制所有数据库操作的执行时间。
	// 这是生产级应用中处理并发和请求超时的标准做法。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10秒超时
	defer cancel() // 确保在 ctx 作用域结束时取消上下文

	fmt.Println("--- 数据库操作开始 ---")

	// --- 1. 创建表 (如果不存在) ---
	// 使用 Exec 方法执行 DDL (Data Definition Language) 语句
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err = executor.Exec(ctx, createTableSQL)
	if err != nil {
		log.Fatalf("创建表 'users' 失败: %v", err)
	}
	fmt.Println("表 'users' 创建或已存在。")

	// --- 2. 插入数据 ---
	fmt.Println("\n--- 插入数据 ---")
	insertSQL := "INSERT INTO users(name, email, created_at) VALUES(?, ?, ?)"

	// 插入第一条数据
	userID1, err := executor.Insert(ctx, insertSQL, "张三", "zhangsan@example.com", time.Now())
	if err != nil {
		fmt.Printf("插入用户张三失败: %v\n", err)
	} else {
		fmt.Printf("插入用户张三成功，ID: %d\n", userID1)
	}

	// 插入第二条数据
	userID2, err := executor.Insert(ctx, insertSQL, "李四", "lisi@example.com", time.Now())
	if err != nil {
		fmt.Printf("插入用户李四失败: %v\n", err)
	} else {
		fmt.Printf("插入用户李四成功，ID: %d\n", userID2)
	}

	// 尝试插入重复 email 的数据（预期会失败，因为 email 字段有 UNIQUE 约束）
	_, err = executor.Insert(ctx, insertSQL, "张三的克隆", "zhangsan@example.com", time.Now())
	if err != nil {
		// 打印错误信息，但不会中断程序，因为这是预期行为
		fmt.Printf("尝试插入重复 email 失败 (预期错误): %v\n", err)
	}

	// --- 3. 查询单条数据 ---
	fmt.Println("\n--- 查询单条数据 ---")
	var user1 User
	querySingleSQL := "SELECT id, name, email, created_at FROM users WHERE email = ?"

	// 使用 QueryRow 方法查询单条数据，并传入一个匿名函数来处理扫描结果
	err = executor.QueryRow(ctx, querySingleSQL, func(row *sql.Row) error {
		// row.Scan() 将数据库行的数据扫描到 Go 结构体的字段中
		return row.Scan(&user1.ID, &user1.Name, &user1.Email, &user1.CreatedAt)
	}, "zhangsan@example.com")

	if err != nil {
		if db_executor.IsNoRowsError(err) { // 使用自定义的错误判断函数
			fmt.Println("未找到 email 为 'zhangsan@example.com' 的用户。")
		} else {
			fmt.Printf("查询用户张三失败: %v\n", err)
		}
	} else {
		fmt.Printf("查询到用户张三: %+v\n", user1) // %+v 会打印结构体字段名和值
	}

	// 尝试查询一个不存在的用户
	var userNotFound User
	err = executor.QueryRow(ctx, querySingleSQL, func(row *sql.Row) error {
		return row.Scan(&userNotFound.ID, &userNotFound.Name, &userNotFound.Email, &userNotFound.CreatedAt)
	}, "nonexistent@example.com")
	if db_executor.IsNoRowsError(err) {
		fmt.Println("未找到 email 为 'nonexistent@example.com' 的用户 (预期)。")
	}


	// --- 4. 查询多条数据 ---
	fmt.Println("\n--- 查询多条数据 ---")
	queryMultiSQL := "SELECT id, name, email, created_at FROM users WHERE name LIKE ? ORDER BY id DESC"
	// 使用 Query 方法查询多条数据
	rows, err := executor.Query(ctx, queryMultiSQL, "%%") // "%%" 匹配所有用户，类似 SQL 的 '%'
	if err != nil {
		log.Fatalf("查询所有用户失败: %v", err)
	}
	defer rows.Close() // **重要：** 确保关闭 rows，释放数据库资源

	users := []User{} // 用于存储查询到的用户列表
	for rows.Next() { // 遍历每一行结果
		var u User
		// 将当前行的数据扫描到 User 结构体中
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			log.Fatalf("扫描用户数据失败: %v", err)
		}
		users = append(users, u)
	}
	// 检查遍历过程中是否有其他错误发生（例如网络中断）
	if err := rows.Err(); err != nil {
		log.Fatalf("遍历用户数据错误: %v", err)
	}

	fmt.Printf("查询到 %d 位用户:\n", len(users))
	for _, u := range users {
		fmt.Printf("  %+v\n", u)
	}

	// --- 5. 更新数据 ---
	fmt.Println("\n--- 更新数据 ---")
	updateSQL := "UPDATE users SET name = ?, created_at = ? WHERE id = ?"
	// 更新用户张三的名字和创建时间
	rowsAffected, err := executor.Update(ctx, updateSQL, "更新后的张三", time.Now(), userID1)
	if err != nil {
		fmt.Printf("更新用户张三失败: %v\n", err)
	} else {
		fmt.Printf("更新用户张三成功，影响行数: %d\n", rowsAffected)
	}

	// --- 6. 事务操作 ---
	fmt.Println("\n--- 事务操作 ---")
	// 使用 WithTransaction 方法执行一个原子性操作块
	txErr := executor.WithTransaction(ctx, func(tx *sql.Tx) error {
		// 在事务中执行插入操作
		_, err := tx.ExecContext(ctx, "INSERT INTO users(name, email, created_at) VALUES(?, ?, ?)", "王五", "wangwu@example.com", time.Now())
		if err != nil {
			fmt.Println("事务内插入王五失败 (将回滚):", err)
			return err // 返回错误将导致整个事务回滚
		}
		fmt.Println("事务内插入王五成功")

		// 假设这里发生了一个导致事务失败的错误，例如尝试插入一个已经存在的 email
		// 因为 email 有唯一约束，这将导致错误并触发回滚
		_, err = tx.ExecContext(ctx, "INSERT INTO users(name, email, created_at) VALUES(?, ?, ?)", "赵六", "zhangsan@example.com", time.Now())
		if err != nil {
			fmt.Println("事务内插入赵六失败 (预期回滚):", err)
			return err // 返回错误，事务将回滚，王五的插入也会被撤销
		}
		fmt.Println("事务内插入赵六成功") // 这行代码通常不会被执行到，因为前面会返回错误

		return nil // 返回 nil 表示事务中的所有操作都成功
	}, nil) // sql.TxOptions 可以为 nil，表示使用数据库默认的隔离级别

	if txErr != nil {
		fmt.Printf("事务执行失败，已回滚: %v\n", txErr)
	} else {
		fmt.Println("事务执行成功，已提交。")
	}

	// --- 7. 删除数据 ---
	fmt.Println("\n--- 删除数据 ---")
	deleteSQL := "DELETE FROM users WHERE email = ?"
	// 删除用户李四
	rowsAffected, err = executor.Delete(ctx, deleteSQL, "lisi@example.com")
	if err != nil {
		fmt.Printf("删除用户李四失败: %v\n", err)
	} else {
		fmt.Printf("删除用户李四成功，影响行数: %d\n", rowsAffected)
	}

	fmt.Println("\n--- 数据库操作结束 ---")
}
```
