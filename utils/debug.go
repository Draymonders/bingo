package utils

var debugFlag bool

// OpenDebug 开启debug
func OpenDebug() {
	debugFlag = true
}

// IsDebug 是否打开debug
func IsDebug() bool {
	return debugFlag
}
