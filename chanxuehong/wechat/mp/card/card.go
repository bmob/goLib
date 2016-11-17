// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     gaowenbin(gaowenbinmarr@gmail.com), chanxuehong(chanxuehong@gmail.com)

package card

import (
	"errors"
	"fmt"

	"github.com/chanxuehong/wechat/mp"
)

// 创建卡券接口.
func (clt *Client) CardCreate(card *Card) (cardId string, err error) {
	if card == nil {
		err = errors.New("nil card")
		return
	}

	var request = struct {
		*Card `json:"card,omitempty"`
	}{
		Card: card,
	}

	var result struct {
		mp.Error
		CardId string `json:"card_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	cardId = result.CardId
	return
}

// 查询卡券详情.
func (clt *Client) CardGet(cardId string) (card *Card, err error) {
	var request = struct {
		CardId string `json:"card_id"`
	}{
		CardId: cardId,
	}

	var result struct {
		mp.Error
		Card `json:"card"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/get?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	card = &result.Card
	return
}

// 删除卡券
func (clt *Client) CardDelete(cardId string) (err error) {
	var request = struct {
		CardId string `json:"card_id"`
	}{
		CardId: cardId,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/delete?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 批量查询卡列表
//  offset: 查询卡列表的起始偏移量，从0 开始，即offset: 5 是指从从列表里的第六个开始读取。
//  count : 需要查询的卡片的数量（数量最大50）
func (clt *Client) CardBatchGet(offset, count int) (cardIdList []string, err error) {
	if offset < 0 {
		err = fmt.Errorf("invalid offset: %d", offset)
		return
	}
	if count < 0 {
		err = fmt.Errorf("invalid count: %d", count)
		return
	}

	var request = struct {
		Offset int `json:"offset"`
		Count  int `json:"count"`
	}{
		Offset: offset,
		Count:  count,
	}

	var result struct {
		mp.Error
		CardIdList []string `json:"card_id_list"`
		TotalNum   int      `json:"total_num"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/batchget?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	if result.TotalNum != len(result.CardIdList) {
		err = errors.New("the total_num and length of card_id_list does not match")
		return
	}
	cardIdList = result.CardIdList
	return
}

// 更改卡券信息接口.
//  支持更新部分通用字段及特殊卡券（会员卡、飞机票、电影票、红包）中特定字段的信息。
//  注：更改卡券的部分字段后会重新提交审核，详情见字段说明。
func (clt *Client) CardUpdate(cardId string, card *Card) (err error) {
	if card == nil {
		return errors.New("nil Card")
	}
	card.CardType = "" //

	var request = struct {
		CardId string `json:"card_id"`
		*Card
	}{
		CardId: cardId,
		Card:   card,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}
	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 库存修改接口.
// cardId:      卡券ID
// increaseNum: 增加库存数量, 可以为负数
func (clt *Client) CardModifyStock(cardId string, increaseNum int) (err error) {
	var request struct {
		CardId             string `json:"card_id"`
		IncreaseStockValue int    `json:"increase_stock_value,omitempty"`
		ReduceStockValue   int    `json:"reduce_stock_value,omitempty"`
	}
	request.CardId = cardId
	switch {
	case increaseNum > 0:
		request.IncreaseStockValue = increaseNum
	case increaseNum < 0:
		request.ReduceStockValue = -increaseNum
	default: // increaseNum == 0
		return
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/modifystock?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
