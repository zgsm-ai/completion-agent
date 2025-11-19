package completions

// 补全请求结构
type CompletionRequest struct {
	Model           string                 `json:"model,omitempty"`
	Prompt          string                 `json:"prompt"`
	LanguageID      string                 `json:"language_id,omitempty"`
	ClientID        string                 `json:"client_id,omitempty"`
	CompletionID    string                 `json:"completion_id,omitempty"`
	ProjectPath     string                 `json:"project_path,omitempty"`
	FileProjectPath string                 `json:"file_project_path,omitempty"`
	ImportContent   string                 `json:"import_content,omitempty"`
	Temperature     float64                `json:"temperature,omitempty"`
	TriggerMode     string                 `json:"trigger_mode,omitempty"`
	ParentID        string                 `json:"parent_id,omitempty"`
	Stop            []string               `json:"stop,omitempty"`
	BetaMode        bool                   `json:"beta_mode,omitempty"`
	Verbose         bool                   `json:"verbose,omitempty"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
	Prompts         *PromptOptions         `json:"prompt_options,omitempty"`
	HideScores      *HiddenScoreOptions    `json:"calculate_hide_score,omitempty"`
}

// 提示词选项
type PromptOptions struct {
	Prefix           string `json:"prefix,omitempty"`
	Suffix           string `json:"suffix,omitempty"`
	CursorLinePrefix string `json:"cursor_line_prefix,omitempty"`
	CursorLineSuffix string `json:"cursor_line_suffix,omitempty"`
	CodeContext      string `json:"code_context,omitempty"`
}

// 计算隐藏分数配置
type HiddenScoreOptions struct {
	IsWhitespaceAfterCursor bool   `json:"is_whitespace_after_cursor"` //光标之后该行是否没有内容(空白除外)
	Prefix                  string `json:"prefix,omitempty"`           //光标前的所有内容
	DocumentLength          int    `json:"document_length"`            //文档长度
	PromptEndPos            int    `json:"prompt_end_pos"`             //光标在文档中的偏移
	PreviousLabel           int    `json:"previous_label"`             //上个请求是否被接受
	PreviousLabelTimestamp  int64  `json:"previous_label_timestamp"`   //上个请求被接受的时间戳
}
