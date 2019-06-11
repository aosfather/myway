package core

import (
	"time"
)

/**
  token
    key
    有效时间
    user
    role
*/

type TokenStore interface {
	//获取token
	GetToken(id string) *AccessToken
	//保存token
	SaveToken(id string, t *AccessToken, expire int64)
}

type TokenManager struct {
	Expire int64 //过期时间
	Store  TokenStore
}

func (this *TokenManager) GetToken(id string) *AccessToken {
	if this.Store != nil {
		token := this.Store.GetToken(id)
		//无效的和超期的会认为无效
		if token != nil && token.Validate(id) && !token.IsExpire(this.Expire) {
			return token
		}

	}

	return nil
}

func (this *TokenManager) CreateToken(user string, role string) *AccessToken {
	id := CreateUUID() //创建唯一的ID
	vcode := GetMd5str(id)
	t := AccessToken{vcode, user, time.Now().Unix(), role}
	if this.Store != nil {
		this.Store.SaveToken(id, &t, this.Expire)
	}

	return &t
}

//访问token
type AccessToken struct {
	Id         string
	User       string
	CreateTime int64
	Role       string
}

//校验token
func (this *AccessToken) Validate(value string) bool {
	if GetMd5str(value) == this.Id {
		return true
	}
	return false
}

func (this *AccessToken) IsExpire(expire int64) bool {
	now := time.Now().Unix()
	if now-this.CreateTime > expire {
		return true
	}

	return false
}
