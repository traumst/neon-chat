package utils

type CtxKey string

const (
	ReqIdKey           CtxKey = "X-Request-Id"
	ActiveUser         CtxKey = "X-Active-User"
	AppState           CtxKey = "X-App-State"
	DBConn             CtxKey = "X-DB-Conn"
	DBTx               CtxKey = "X-DB-Transaction"
	SessionCookie      string = "session"
	LetterBytes        string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Byte                      = 1
	KB                        = 1024 * Byte
	MB                        = 1024 * KB
	GB                        = 1024 * MB
	MaxFileName        int    = 120
	MaxUploadBytesSize int64  = 50 * KB
	MaxCacheSize       int    = 1024
)
