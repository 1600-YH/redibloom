package re_bloom

import (
	"github.com/spaolacci/murmur3"
	"math"
)

type Encryptor struct {
}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) Encrypt(origin string) int32 {
	// 创建一个32位的murmur3 hash函数实例
	hasher := murmur3.New32()
	// origin转为字节数组，并写入hash函数实例中
	_, _ = hasher.Write([]byte(origin))
	// 计算并返回写入数据的32位哈希值，通过取模将值限定在 int32的范围内
	return int32(hasher.Sum32() % math.MaxInt32)
}
