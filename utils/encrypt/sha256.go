package encrypt

import (
	"encoding/hex"

	"crypto/sha256"
)

func Sha256(input string) string {
	// 创建一个新的SHA-256哈希对象
	hash := sha256.New()

	// 将数据写入哈希对象
	hash.Write([]byte(input))

	// 计算哈希值并以十六进制字符串形式输出
	// Sum(nil) 表示返回结果为一个字节切片，如果你想要直接打印哈希值，可以使用 `Sum(nil)`
	// 如果你需要更多操作，可以传递一个已有的字节切片，比如 `Sum([]byte{})`
	// 注意：这个过程不会修改哈希对象，可以进行多次调用

	return hex.EncodeToString(hash.Sum(nil))
}
