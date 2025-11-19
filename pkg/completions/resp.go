package completions

import (
	"completion-agent/pkg/metrics"
	"completion-agent/pkg/model"
	"context"
	"fmt"
	"time"
)

/**
 * 补全使用情况统计结构体
 * @description
 * - 记录补全请求的token使用情况
 * - 包含输入token数、输出token数和总token数
 * - 用于监控和统计API使用情况
 * - 嵌入到CompletionResponse中返回给客户端
 * @example
 * usage := CompletionUsage{
 *     PromptTokens: 100,
 *     CompletionTokens: 50,
 *     TotalTokens: 150,
 * }
 */
type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

/**
 * 补全选择结构体
 * @description
 * - 表示补全请求的一个选择结果
 * - 包含生成的文本内容
 * - 支持多个选择结果，按优先级排序
 * - 用于向客户端返回补全建议
 * @example
 * choice := CompletionChoice{
 *     Text: "function example() {\n    return 'hello';\n}",
 * }
 */
type CompletionChoice struct {
	Text string `json:"text"`
}

/**
 * 补全性能统计结构体
 * @description
 * - 记录补全请求各阶段的性能数据
 * - 包含接收时间、上下文获取时间、排队时间、LLM处理时间和总时间
 * - 记录token使用统计信息
 * - 用于性能监控和优化分析
 * @example
 * perf := &CompletionPerformance{
 *     ReceiveTime: time.Now(),
 *     ContextDuration: 100 * time.Millisecond,
 *     LLMDuration: 500 * time.Millisecond,
 *     TotalDuration: 650 * time.Millisecond,
 * }
 */
type CompletionPerformance struct {
	ReceiveTime      time.Time     `json:"receive_time"`     //收到请求的时间
	EnqueueTime      time.Time     `json:"-"`                //开始排队时间
	ContextDuration  time.Duration `json:"context_duration"` //获取上下文的时长
	QueueDuration    time.Duration `json:"queue_duration"`   //排队时长
	LLMDuration      time.Duration `json:"llm_duration"`     //调用大语言模型耗用的时长
	TotalDuration    time.Duration `json:"total_duration"`   //总时长
	PromptTokens     int           `json:"prompt_tokens"`
	CompletionTokens int           `json:"completion_tokens"`
	TotalTokens      int           `json:"total_tokens"`
}

/**
 * 补全响应结构体
 * @description
 * - 表示补全请求的完整响应
 * - 包含响应ID、模型名称、补全选择列表、使用统计和状态
 * - 支持错误信息和详细输出
 * - 用于向客户端返回补全结果
 * @example
 * response := &CompletionResponse{
 *     ID: "req-123",
 *     Model: "gpt-3.5-turbo",
 *     Object: "text_completion",
 *     Choices: []CompletionChoice{{Text: "completion text"}},
 *     Status: model.StatusSuccess,
 * }
 */
type CompletionResponse struct {
	ID      string                   `json:"id"`
	Model   string                   `json:"model"`
	Object  string                   `json:"object"`
	Choices []CompletionChoice       `json:"choices"`
	Created int                      `json:"created"`
	Usage   CompletionPerformance    `json:"usage"`
	Status  model.CompletionStatus   `json:"status"`
	Error   string                   `json:"error"`
	Verbose *model.CompletionVerbose `json:"verbose,omitempty"`
}

/**
 * 记录补全性能指标
 * @param {string} modelName - 模型名称，用于指标分类
 * @param {string} status - 补全状态字符串，用于结果分类
 * @param {*CompletionPerformance} perf - 性能统计对象，包含各阶段耗时和token使用情况
 * @description
 * - 记录补全请求的各阶段耗时指标
 * - 记录补全请求计数指标
 * - 记录输入和输出token使用指标
 * - 使用metrics包进行指标上报
 * - 用于监控补全服务的性能和资源使用情况
 * @example
 * perf := &CompletionPerformance{
 *     QueueDuration: 50 * time.Millisecond,
 *     ContextDuration: 100 * time.Millisecond,
 *     LLMDuration: 500 * time.Millisecond,
 *     TotalDuration: 650 * time.Millisecond,
 *     PromptTokens: 100,
 *     CompletionTokens: 50,
 * }
 * Metrics("gpt-3.5-turbo", "success", perf)
 */
func Metrics(modelName string, status string, perf *CompletionPerformance) {
	metrics.RecordCompletionDuration(modelName, status,
		perf.QueueDuration, perf.ContextDuration, perf.LLMDuration, perf.TotalDuration)
	metrics.IncrementCompletionRequests(modelName, status)
	metrics.RecordCompletionTokens(modelName, metrics.TokenTypeInput, perf.PromptTokens)
	metrics.RecordCompletionTokens(modelName, metrics.TokenTypeOutput, perf.CompletionTokens)
}

/**
 * 创建错误响应
 * @param {*CompletionInput} input - 补全输入对象，包含请求信息
 * @param {model.CompletionStatus} status - 补全状态，表示错误类型
 * @param {*CompletionPerformance} perf - 性能统计对象，包含耗时和token信息
 * @param {error} err - 错误对象，包含错误详情
 * @returns {*CompletionResponse} 返回错误响应对象
 * @description
 * - 创建表示错误的补全响应
 * - 如果错误为nil，使用状态字符串作为错误信息
 * - 记录性能指标到监控系统
 * - 设置空的选择结果
 * - 包含错误详情和性能统计信息
 * @example
 * input := &CompletionInput{...}
 * perf := &CompletionPerformance{...}
 * response := ErrorResponse(input, model.StatusReqError, perf, fmt.Errorf("invalid request"))
 * // response.Status = "reqError"
 * // response.Error = "invalid request"
 * // response.Choices = [{}]
 */
func ErrorResponse(input *CompletionInput, status model.CompletionStatus,
	perf *CompletionPerformance, err error) *CompletionResponse {
	if err == nil {
		err = fmt.Errorf("%s", string(status))
	}
	Metrics(input.SelectedModel, string(status), perf)
	return &CompletionResponse{
		ID:      input.CompletionID,
		Model:   input.SelectedModel,
		Object:  "text_completion",
		Choices: []CompletionChoice{{Text: ""}}, // 使用后置处理后的补全结果
		Created: int(perf.ReceiveTime.Unix()),
		Usage:   *perf,
		Status:  status,
		Error:   err.Error(),
	}
}

/**
 * 创建成功响应
 * @param {*CompletionInput} input - 补全输入对象，包含请求信息
 * @param {string} completionText - 补全文本内容，表示生成的代码
 * @param {*CompletionPerformance} perf - 性能统计对象，包含耗时和token信息
 * @returns {*CompletionResponse} 返回成功响应对象
 * @description
 * - 创建表示成功的补全响应
 * - 设置状态为成功
 * - 记录性能指标到监控系统
 * - 包含补全文本和性能统计信息
 * - 不包含错误信息
 * @example
 * input := &CompletionInput{...}
 * perf := &CompletionPerformance{...}
 * response := SuccessResponse(input, "function test() {\n    return;\n}", perf)
 * // response.Status = "sucess"
 * // response.Choices[0].Text = "function test() {\n    return;\n}"
 * // response.Error = ""
 */
func SuccessResponse(input *CompletionInput, completionText string, perf *CompletionPerformance) *CompletionResponse {
	Metrics(input.SelectedModel, string(model.StatusSuccess), perf)
	return &CompletionResponse{
		ID:      input.CompletionID,
		Model:   input.SelectedModel,
		Object:  "text_completion",
		Choices: []CompletionChoice{{Text: completionText}}, // 使用后置处理后的补全结果
		Created: int(perf.ReceiveTime.Unix()),
		Usage:   *perf,
		Status:  model.StatusSuccess,
	}
}

/**
 * 创建取消请求响应
 * @param {*CompletionInput} input - 补全输入对象，包含请求信息
 * @param {*CompletionPerformance} perf - 性能统计对象，包含耗时和token信息
 * @param {error} err - 错误对象，包含取消原因
 * @returns {*CompletionResponse} 返回取消请求响应对象
 * @description
 * - 创建表示请求取消的补全响应
 * - 根据错误类型判断是超时还是主动取消
 * - 计算总耗时并记录性能指标
 * - 设置空的选择结果
 * - 包含错误详情和性能统计信息
 * @example
 * input := &CompletionInput{...}
 * perf := &CompletionPerformance{...}
 * response := CancelRequest(input, perf, context.Canceled)
 * // response.Status = "canceled"
 * // response.Error = "context canceled"
 *
 * response = CancelRequest(input, perf, fmt.Errorf("timeout"))
 * // response.Status = "timeout"
 * // response.Error = "timeout"
 */
func CancelRequest(input *CompletionInput, perf *CompletionPerformance, err error) *CompletionResponse {
	status := model.StatusTimeout
	if err.Error() == context.Canceled.Error() {
		status = model.StatusCanceled
	}
	perf.TotalDuration = time.Since(perf.ReceiveTime)
	Metrics(input.SelectedModel, string(status), perf)
	return &CompletionResponse{
		ID:      input.CompletionID,
		Model:   input.SelectedModel,
		Object:  "text_completion",
		Choices: []CompletionChoice{{Text: ""}},
		Created: int(perf.ReceiveTime.Unix()),
		Usage:   *perf,
		Status:  status,
		Error:   err.Error(),
	}
}

/**
 * 创建拒绝请求响应
 * @param {*CompletionInput} input - 补全输入对象，包含请求信息
 * @param {*CompletionPerformance} perf - 性能统计对象，包含耗时和token信息
 * @param {model.CompletionStatus} status - 补全状态，表示拒绝类型
 * @param {error} err - 错误对象，包含拒绝原因
 * @returns {*CompletionResponse} 返回拒绝请求响应对象
 * @description
 * - 创建表示请求被拒绝的补全响应
 * - 记录性能指标到监控系统
 * - 设置空的选择结果
 * - 包含拒绝原因和性能统计信息
 * - 用于过滤器链拒绝请求的情况
 * @example
 * input := &CompletionInput{...}
 * perf := &CompletionPerformance{...}
 * response := RejectRequest(input, perf, model.StatusRejected, fmt.Errorf("request rejected by filter"))
 * // response.Status = "reject"
 * // response.Error = "request rejected by filter"
 * // response.Choices = [{}]
 */
func RejectRequest(input *CompletionInput, perf *CompletionPerformance, status model.CompletionStatus, err error) *CompletionResponse {
	Metrics(input.SelectedModel, string(status), perf)
	return &CompletionResponse{
		ID:      input.CompletionID,
		Model:   input.SelectedModel,
		Object:  "text_completion",
		Choices: []CompletionChoice{{Text: ""}},
		Created: int(perf.ReceiveTime.Unix()),
		Usage:   *perf,
		Status:  status,
		Error:   err.Error(),
	}
}
