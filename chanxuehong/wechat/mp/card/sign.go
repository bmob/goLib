// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     gaowenbin(gaowenbinmarr@gmail.com), chanxuehong(chanxuehong@gmail.com)

package card

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

// 卡卷通用签名方法.
func Sign(strs []string) (signature string) {
	sort.Strings(strs)

	Hash := sha1.New()
	for _, str := range strs {
		Hash.Write([]byte(str))
	}
	return hex.EncodeToString(Hash.Sum(nil))
}
