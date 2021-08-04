package errmsg

const (
	SUCCESS       = 200
	UPDATE_ACCESS = 300
	ERROR         = 400

	CAN_NOT_GET_ACCESS_IP = 1001
	CAN_NOT_BIND_PARAMS   = 1002

	DB_OPTIONS_FAILED             = 2005
	REACH_MAX_SPRINKLE_TIMES      = 2006
	ETHERUM_ADDRESS_HAS_SPRINKLED = 2007
	INSERT_RECORD_FAILED          = 2008

	DECODE_SPRINKLED_ADDR_FAILED  = 3001
	REQUEST_TOKEN_TRANSFER_FAILED = 3002
	TRANSACTION_MINING_FAILED     = 3003
	// REQUEST_TOKEN_TRANSFER_FAILED =3002
)

var errMessage = map[int]string{
	SUCCESS: "请求成功",
	ERROR:   "请求失败",

	CAN_NOT_GET_ACCESS_IP: "无法获取访问者IP",
	CAN_NOT_BIND_PARAMS:   "无法获取绑定的参数",
	DB_OPTIONS_FAILED:     "数据库操作失败",

	REACH_MAX_SPRINKLE_TIMES:      "达到最大放水次数",
	ETHERUM_ADDRESS_HAS_SPRINKLED: "该地址已发过水",
	INSERT_RECORD_FAILED:          "创建记录失败",

	DECODE_SPRINKLED_ADDR_FAILED:  "解析出水地址错误",
	REQUEST_TOKEN_TRANSFER_FAILED: "请求token transfer失败",
	TRANSACTION_MINING_FAILED:     "",
}

func GetErrMsg(code int) string {
	return errMessage[code]
}