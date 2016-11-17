// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package account

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/chanxuehong/wechat/mp"
)

const (
	TemporaryQRCodeExpireSecondsLimit = 1800   // 临时二维码 expire_seconds 限制
	PermanentQRCodeSceneIdLimit       = 100000 // 永久二维码 scene_id 限制
)

// 永久二维码
type PermanentQRCode struct {
	// 下面两个字段同时只有一个有效
	SceneId     uint32 `json:"scene_id"`  // 场景值
	SceneString string `json:"scene_str"` // 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64

	Ticket string `json:"ticket"` // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	URL    string `json:"url"`    // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

// 二维码图片的URL, 可以GET此URL下载二维码或者在线显示此二维码.
func (qrcode *PermanentQRCode) PicURL() string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qrcode.Ticket)
}

// 临时二维码
type TemporaryQRCode struct {
	PermanentQRCode
	ExpiresIn int `json:"expire_seconds"` // 二维码的有效时间，以秒为单位。
}

// 创建临时二维码
//  SceneId:       场景值ID，为32位非0整型
//  ExpireSeconds: 二维码的有效时间，以秒为单位。
func (clt *Client) CreateTemporaryQRCode(SceneId uint32, ExpireSeconds int) (qrcode *TemporaryQRCode, err error) {
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneId uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	var result struct {
		mp.Error
		TemporaryQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.TemporaryQRCode.SceneId = SceneId
	qrcode = &result.TemporaryQRCode
	return
}

// 创建永久二维码
//  SceneId: 场景值ID，最大值为100000（目前参数只支持1--100000）
func (clt *Client) CreatePermanentQRCode(SceneId uint32) (qrcode *PermanentQRCode, err error) {
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneId uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	var result struct {
		mp.Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneId = SceneId
	qrcode = &result.PermanentQRCode
	return
}

// 创建永久二维码
//  SceneString: 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
func (clt *Client) CreatePermanentQRCodeWithSceneString(SceneString string) (qrcode *PermanentQRCode, err error) {
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	var result struct {
		mp.Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneString = SceneString
	qrcode = &result.PermanentQRCode
	return
}

// 通过ticket换取二维码, 写入到 writer.
//  NOTE: 调用者保证所有参数有效.
func qrcodeDownloadToWriter(ticket string, writer io.Writer, httpClient *http.Client) (err error) {
	qrcodeURL := "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket)
	httpResp, err := httpClient.Get(qrcodeURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusOK {
		_, err = io.Copy(writer, httpResp.Body)
		return
	}
	return errors.New("下载二维码出错, ticket: " + ticket)
}

// 通过ticket换取二维码, 写入到 writer.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func QRCodeDownloadToWriter(ticket string, writer io.Writer, httpClient *http.Client) (err error) {
	if writer == nil {
		return errors.New("nil writer")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, httpClient)
}

// 通过ticket换取二维码, 写入到 writer.
func (clt *Client) QRCodeDownloadToWriter(ticket string, writer io.Writer) (err error) {
	if writer == nil {
		return errors.New("nil writer")
	}
	if clt.HttpClient == nil {
		clt.HttpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, clt.HttpClient)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
//  如果 httpClient == nil 则默认用 http.DefaultClient
func QRCodeDownload(ticket, filepath string, httpClient *http.Client) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, httpClient)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
func (clt *Client) QRCodeDownload(ticket, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	if clt.HttpClient == nil {
		clt.HttpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, clt.HttpClient)
}
