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
	NodeQRCodeError     string = "二维码生成失败"
	NodeURLError        string = "URL生成失败"

	UsernameExist    string = "用户名已存在"
	NodeNameExist    string = "节点名称名已存在"
	NodeNotExist     string = "不存在该节点"
	NodeTypeNotExist string = "不存在该节点类型"
	RoleNotExist     string = "不存在该角色"
	SystemNotExist   string = "不存在该系统设置"

	UserRegisterClosed string = "用户注册功能已关闭"
	UserDisabled       string = "该用户已被禁用"

	FileSizeTooBig  string = "文件大小不能超过10MB"
	FileFormatError string = "文件格式只支持.zip"
	FileUploadError string = "文件上传失败"

	GrpcError        string = "远程服务连接失败"
	LoadKeyPairError string = "加载本机密钥和证书失败"

	SendEmailError string = "邮件发送失败"

	RateLimiterError string = "点的太快啦"

	BlackListError string = "由于您近期异常请求过于频繁,已限制访问,如需取消限制,请联系管理员"
)
