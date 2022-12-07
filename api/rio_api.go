package api

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

const domain = "https://api.appmeta.cn"

type UserInfo struct {
	Email string `json:"email"`
}

type FirApi struct {
	ApiToken string
	Email    string
}

func (f *FirApi) Login(token string) error {
	url := domain + "/user"
	client := resty.New()
	// body := `{"api_token":` + token + `}`

	resp, err := client.R().SetQueryParam("api_token", token).SetHeader("Content-Type", "application/json").Get(url)

	if err != nil || resp.StatusCode() != 200 {
		return err
	}
	var userInfo UserInfo

	json.Unmarshal(resp.Body(), &userInfo)
	f.Email = userInfo.Email
	f.ApiToken = token
	return nil
}
