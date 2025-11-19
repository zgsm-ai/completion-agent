package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
 * 全局 logger 实例
 * @description
 * - 全局日志记录器实例
 * - 使用zap高性能日志库
 * - 在init函数中初始化
 * - 提供应用程序的统一日志记录入口
 * @example
 * Logger.Info("应用启动")
 * Logger.Error("发生错误", zap.Error(err))
 */
var Logger *zap.Logger

/**
 * 初始化日志系统
 * @description
 * - 创建生产级别的日志配置
 * - 设置日志级别为Info
 * - 配置时间戳格式为本地时区
 * - 构建日志记录器实例
 * - 替换zap的全局logger
 * @throws
 * - 如果日志构建失败，会导致程序panic
 * @example
 * // 包初始化时自动调用
 * // 不需要手动调用
 */
func init() {
	// 使用 NewProductionConfig 并可选调整日志级别
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // 设置日志级别为 Info

	// 配置时间戳格式为本地时区
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
	}

	// 注意：禁止将日志输出到文件
	// config.OutputPaths = []string{"./app.log"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	// 替换 zap 的全局 logger
	zap.ReplaceGlobals(Logger)
}

/**
 * 设置日志级别
 * @param {string} level - 日志级别字符串，如"debug", "info", "warn", "error"
 * @description
 * - 解析输入的日志级别字符串
 * - 如果解析失败，记录警告日志并使用默认级别
 * - 更新全局logger的核心日志级别
 * - 支持的标准级别：debug, info, warn, error, dpanic, panic, fatal
 * @example
 * SetLevel("debug")
 * // 设置日志级别为debug，将显示更详细的日志
 *
 * SetLevel("invalid")
 * // 输出警告: Invalid log level, using default level (info)
 */
func SetLevel(level string) {
	levelValue, err := zapcore.ParseLevel(level)
	if err != nil {
		Logger.Warn("Invalid log level, using default level (info)")
		return
	}
	Logger.Core().Enabled(levelValue)
}

/**
 * 设置日志模式
 * @param {string} mode - 运行模式字符串，"debug"或"release"
 * @description
 * - 根据运行模式创建不同的日志配置
 * - debug模式使用开发配置，提供更友好的输出格式
 * - release模式使用生产配置，优化性能
 * - 两种模式都使用本地时区显示时间戳
 * - 替换zap的全局logger
 * @throws
 * - 如果日志构建失败，会导致程序panic
 * @example
 * SetMode("debug")
 * // 设置为调试模式，日志输出更易读
 *
 * SetMode("release")
 * // 设置为发布模式，日志输出性能优化
 */
func SetMode(mode string) {
	var l *zap.Logger
	var err error
	if mode == "debug" {
		// 开发模式也使用本地时区
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
		}
		l, err = config.Build()
	} else {
		// 生产模式也使用本地时区
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
		}
		l, err = config.Build()
	}
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
}

/**
 * 刷新所有日志到输出
 * @description
 * - 调用全局logger的Sync方法
 * - 确保所有缓冲的日志都被写入输出
 * - 通常在应用程序退出前调用
 * - 用于保证日志数据不丢失
 * @example
 * defer Sync()
 * // 在main函数结束时调用，确保所有日志都被写入
 */
func Sync() {
	Logger.Sync()
}

// 便捷函数，直接调用全局 logger 的方法

/**
 * 记录信息级别日志
 * @param {string} msg - 日志消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Info级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Info方法
 * - 用于记录常规的应用程序运行信息
 * @example
 * Info("用户登录", zap.String("username", "john"), zap.Int("userId", 123))
 * // 输出: {"level":"info","msg":"用户登录","username":"john","userId":123}
 */
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

/**
 * 记录错误级别日志
 * @param {string} msg - 错误消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Error级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Error方法
 * - 用于记录应用程序中的错误情况
 * @example
 * Error("数据库连接失败", zap.Error(err), zap.String("host", "localhost"))
 * // 输出: {"level":"error","msg":"数据库连接失败","error":"connection refused","host":"localhost"}
 */
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

/**
 * 记录调试级别日志
 * @param {string} msg - 调试消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Debug级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Debug方法
 * - 用于记录详细的调试信息，仅在调试模式下输出
 * @example
 * Debug("处理请求", zap.String("method", "GET"), zap.String("path", "/api/users"))
 * // 在debug模式下输出: {"level":"debug","msg":"处理请求","method":"GET","path":"/api/users"}
 */
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

/**
 * 记录警告级别日志
 * @param {string} msg - 警告消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Warn级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Warn方法
 * - 用于记录可能需要注意但不会导致程序错误的情况
 * @example
 * Warn("缓存即将过期", zap.Time("expireTime", time.Now().Add(time.Hour)))
 * // 输出: {"level":"warn","msg":"缓存即将过期","expireTime":"2023-01-01T12:00:00Z"}
 */
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

/**
 * 记录致命错误级别日志
 * @param {string} msg - 致命错误消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Fatal级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Fatal方法
 * - 记录后会导致程序调用os.Exit(1)退出
 * - 用于记录无法恢复的致命错误
 * @example
 * Fatal("配置文件读取失败", zap.Error(err), zap.String("configPath", "/etc/app/config.json"))
 * // 输出错误信息后程序退出
 */
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

/**
 * 记录恐慌级别日志
 * @param {string} msg - 恐慌消息内容
 * @param {...zap.Field} fields - 可变参数，额外的日志字段
 * @description
 * - 记录Panic级别的日志消息
 * - 支持添加结构化的字段信息
 * - 委托给全局Logger的Panic方法
 * - 记录后会导致程序panic
 * - 用于记录严重错误并触发panic机制
 * @example
 * Panic("系统状态异常", zap.String("state", "critical"), zap.Int("errorCode", 500))
 * // 输出错误信息后触发panic
 */
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

/**
 * 创建带有额外字段的 logger
 * @param {...zap.Field} fields - 可变参数，要添加到logger的字段
 * @returns {*zap.Logger} 返回带有额外字段的新logger实例
 * @description
 * - 基于全局logger创建新的logger实例
 * - 新logger包含所有指定的额外字段
 * - 每次调用都创建新的logger实例
 * - 用于在特定上下文中添加固定的日志字段
 * @example
 * userLogger := With(zap.String("userId", "123"), zap.String("sessionId", "abc"))
 * userLogger.Info("用户操作")
 * // 输出: {"level":"info","msg":"用户操作","userId":"123","sessionId":"abc"}
 */
func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}
