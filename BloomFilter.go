package re_bloom

import (
	"context"
	"fmt"
	"github.com/demdxx/gocast"
)

// 布隆过滤器服务
type BloomService struct {
	m, k      int32
	encryptor Encryptor
	client    *RedisClient
}

func NewBloomService(m, k int32, encryptor Encryptor, client *RedisClient) *BloomService {
	return &BloomService{
		m:         m,
		k:         k,
		encryptor: encryptor,
		client:    client,
	}
}

// 检查元素是否存在
func (b *BloomService) Exist(ctx context.Context, key, val string) (bool, error) {
	keyAndArgs := make([]interface{}, 0, b.k+2)
	keyAndArgs = append(keyAndArgs, key, b.k)
	// 加密
	for _, encrypted := range b.getKEncrypted(val) {
		keyAndArgs = append(keyAndArgs, encrypted)
	}

	rawResp, err := b.client.Eval(ctx, LuaBloomBatchGetBits, 1, keyAndArgs)
	if err != nil {
		return false, err
	}

	// 对lua脚本返回的结果转int
	resp := gocast.ToInt(rawResp)
	// 如果为1，表示可能存在，否则就是一定不存在
	if resp == 1 {
		return true, nil
	}
	return false, nil
}

func (b *BloomService) Add(ctx context.Context, key, val string) error {
	keyAndArgs := make([]interface{}, 0, b.k+2)
	keyAndArgs = append(keyAndArgs, key, b.k)
	for _, encrypted := range b.getKEncrypted(val) {
		keyAndArgs = append(keyAndArgs, encrypted)
	}

	rawResp, err := b.client.Eval(ctx, LuaBloomBatchSetBits, 1, keyAndArgs)
	if err != nil {
		return err
	}
	resp := gocast.ToInt(rawResp)
	if resp == 1 {
		return fmt.Errorf("resp: %d", resp)
	}
	return nil
}

func (b *BloomService) getKEncrypted(val string) []int32 {
	encrypteds := make([]int32, 0, b.k)
	origin := val
	for i := 0; int32(i) < b.k; i++ {
		encrypted := b.encryptor.Encrypt(origin)
		encrypteds = append(encrypteds, encrypted%b.m)
		if int32(i) == b.k-1 {
			break
		}
		// 将每次hash映射之后的值拼接原来的val，保证每次输入给hash函数的输入值独一无二，减少哈希冲突的可能
		origin = val + gocast.ToString(encrypted)
	}
	return encrypteds
}
