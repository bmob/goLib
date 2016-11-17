// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"

	"github.com/chanxuehong/wechat/mp"
)

// 下载多媒体到文件.
func (clt *Client) DownloadMedia(mediaId, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.downloadMediaToWriter(mediaId, file)
}

// 下载多媒体到 io.Writer.
func (clt *Client) DownloadMediaToWriter(mediaId string, writer io.Writer) error {
	if writer == nil {
		return errors.New("nil writer")
	}
	return clt.downloadMediaToWriter(mediaId, writer)
}

// 下载多媒体到 io.Writer.
func (clt *Client) downloadMediaToWriter(mediaId string, writer io.Writer) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := "http://file.api.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	ContentType, _, _ := mime.ParseMediaType(httpResp.Header.Get("Content-Type"))
	if ContentType != "text/plain" && ContentType != "application/json" {
		// 返回的是媒体流
		_, err = io.Copy(writer, httpResp.Body)
		return
	}

	// 返回的是错误信息
	var result mp.Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	switch result.ErrCode {
	case mp.ErrCodeOK:
		return // 基本不会出现
	case mp.ErrCodeInvalidCredential, mp.ErrCodeTimeout: // 失效(过期)重试一次
		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			goto RETRY
		}
		fallthrough
	default:
		err = &result
		return
	}
}

// 根据上传的缩略图媒体创建图文消息素材.
//  articles 的长度不能大于 NewsArticleCountLimit.
func (clt *Client) CreateNews(articles []Article) (info *MediaInfo, err error) {
	if len(articles) == 0 {
		err = errors.New("图文消息是空的")
		return
	}
	if len(articles) > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, len(articles))
		return
	}

	var request = struct {
		Articles []Article `json:"articles,omitempty"`
	}{
		Articles: articles,
	}

	var result struct {
		mp.Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 根据上传的视频文件 media_id 创建视频媒体, 群发视频消息应该用这个函数得到的 media_id.
//  NOTE: title, description 可以为空.
func (clt *Client) CreateVideo(mediaId, title, description string) (info *MediaInfo, err error) {
	var request = struct {
		MediaId     string `json:"media_id"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		MediaId:     mediaId,
		Title:       title,
		Description: description,
	}

	var result struct {
		mp.Error
		MediaInfo
	}

	incompleteURL := "https://file.api.weixin.qq.com/cgi-bin/media/uploadvideo?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}
