package constant

// 错误msg
const (
	SysError          string = "系统错误,请联系统管理员"
	ValidateFailed    string = "参数无效"
	UnauthorizedError string = "未认证"
	IllegalTokenError string = "认证失败"
	ForbiddenError    string = "没有权限"
	TokenExpiredError string = "登录过期"
	LogOutError       string = "你还没有登录"

	OriPassError        string = "原密码输入错误"
	UsernameOrPassError string = "用户名或密码错误"
	NodeURLError        string = "URL生成失败"

	UsernameExist       string = "用户名已存在"
	NodeNameExist       string = "节点名称名已存在"
	NodeServerNameExist string = "节点器名称名已存在"
	NodeNotExist        string = "不存在该节点"
	NodeTypeNotExist    string = "不存在该节点类型"
	RoleNotExist        string = "不存在该角色"
	SystemNotExist      string = "不存在该系统设置"
	FileTaskNotExist    string = "不存在该文件任务"

	AccountRegisterClosed string = "用户注册功能已关闭"
	AccountDisabled       string = "该用户已被禁用"
	FileTaskNotSuccess    string = "该任务还没有执行完成"

	FileSizeTooBig  string = "文件太大了"
	FileFormatError string = "文件格式不支持"
	FileUploadError string = "文件上传失败"

	CsvRowNotEnough string = "数据为空"
	CsvTitleError   string = "表头不正确"

	SystemEmailError string = "系统邮箱未设置"

	BlackListError   string = "由于您近期异常操作过于频繁,已限制访问,如需取消限制,请联系管理员"
	RateLimiterError string = "点的太快啦"

	GrpcError        string = "远程服务连接失败,请检查远程服务配置"
	GrpcAddNodeError string = "远程服务添加节点失败,请稍后再试"
	LoadKeyPairError string = "加载本地密钥和证书失败"

	PortIsOccupied         string = "端口被占用,请检查该端口或选择其他端口"
	PortRangeError         string = "端口范围在100-30000之间"
	NodeServerDeletedError string = "该服务器下存在节点"
)
