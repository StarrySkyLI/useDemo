package consts

type IntBool int
type LogType int

const (
	HeaderToken = "AuthorizationJwt"
)

const (
	Trance    = "trance"
	ClientIp  = "client_ip"
	UserAgent = "user_agent"
)

const (
	IntBoolTrue  IntBool = 1
	IntBoolFalse IntBool = 0
)

const (
	IntBoolDbTrue  IntBool = 1 //成功状态
	IntBoolDbFalse IntBool = 2 //失败状态
)

const (
	FrontOperationLogs  LogType = 11 //前台操作日志
	FrontLoginLogs      LogType = 12 //前台登录日志
	SystemOperationLogs LogType = 21 //系统操作日志
	SystemLoginLogs     LogType = 22 //系统登录日志
)

const (
	AppName         = "app_name" // 项目名称
	Business        = "business" // 业务组
	Token           = "token"
	TokenUid        = "token_uid"
	TokenUidRole    = "token_uid_role"
	TranceId        = Trance
	DeviceId        = "device_id" // 设备ID
	Version         = "version"
	Source          = "source"
	IMEI            = "imei"
	PackageName     = "package_name" // 版本号
	ReqHost         = "req_host"
	ReqPath         = "req_path"
	ContentLanguage = "content_language"
	BusinessCode    = "business_code"
	OriginHost      = "origin_host"
)

const BusinessCodeDefaultValue = "10000000"
