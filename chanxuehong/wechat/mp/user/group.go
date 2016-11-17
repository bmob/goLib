// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package user

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

const GroupCountLimit = 100 // 一个公众账号，最多支持创建100个分组。

// 用户分组
type Group struct {
	Id        int64  `json:"id"`    // 分组id, 由微信分配
	Name      string `json:"name"`  // 分组名字, UTF8编码
	UserCount int    `json:"count"` // 分组内用户数量
}

// 创建分组.
//  name: 分组名字（30个字符以内）.
func (clt *Client) CreateGroup(name string) (group *Group, err error) {
	if name == "" {
		err = errors.New(`name == ""`)
		return
	}

	var request struct {
		Group struct {
			Name string `json:"name"`
		} `json:"group"`
	}
	request.Group.Name = name

	var result struct {
		mp.Error
		Group `json:"group"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.Group.UserCount = 0 //
	group = &result.Group
	return
}

// 查询所有分组.
func (clt *Client) ListGroup() (groups []Group, err error) {
	var result = struct {
		mp.Error
		Groups []Group `json:"groups"`
	}{
		Groups: make([]Group, 0, 16),
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/get?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	groups = result.Groups
	return
}

// 修改分组名.
//  name: 分组名字（30个字符以内）.
func (clt *Client) GroupRename(groupId int64, newName string) (err error) {
	if newName == "" {
		return errors.New(`newName == ""`)
	}

	var request struct {
		Group struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"group"`
	}
	request.Group.Id = groupId
	request.Group.Name = newName

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 查询用户所在分组.
func (clt *Client) UserInWhichGroup(openId string) (groupId int64, err error) {
	var request = struct {
		OpenId string `json:"openid"`
	}{
		OpenId: openId,
	}

	var result struct {
		mp.Error
		GroupId int64 `json:"groupid"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/getid?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	groupId = result.GroupId
	return
}

// 移动用户分组.
func (clt *Client) MoveUserToGroup(openId string, toGroupId int64) (err error) {
	var request = struct {
		OpenId    string `json:"openid"`
		ToGroupId int64  `json:"to_groupid"`
	}{
		OpenId:    openId,
		ToGroupId: toGroupId,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/members/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 批量移动用户分组.
func (clt *Client) MoveUsersToGroup(openIdList []string, toGroupId int64) (err error) {
	if len(openIdList) <= 0 {
		return
	}

	var request = struct {
		OpenIdList []string `json:"openid_list,omitempty"`
		ToGroupId  int64    `json:"to_groupid"`
	}{
		OpenIdList: openIdList,
		ToGroupId:  toGroupId,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/groups/members/batchupdate?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
