// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     gaowenbin(gaowenbinmarr@gmail.com), chanxuehong(chanxuehong@gmail.com)

package card

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

type MeetingTicketUpdateUserParameters struct {
	Code   string `json:"code"`              // 用户的门票唯一序列号
	CardId string `json:"card_id,omitempty"` // 要更新门票序列号所述的card_id ， 生成券时use_custom_code 填写true 时必填。

	Zone       string `json:"zone,omitempty"`        // 区域
	Entrance   string `json:"entrance,omitempty"`    // 入口
	SeatNumber string `json:"seat_number,omitempty"` // 座位号
}

// 更新电影票.
//  领取电影票后通过调用“更新电影票”接口update 电影信息及用户选座信息
func (clt *Client) MeetingTicketUpdateUser(para *MeetingTicketUpdateUserParameters) (err error) {
	if para == nil {
		return errors.New("nil MeetingTicketUpdateUserParameters")
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/meetingticket/updateuser?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
