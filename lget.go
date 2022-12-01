package lget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type LoginInfo struct {
	Host    string
	AdminId string
	AdminPw string
}

type RapeHandler interface {
	Login(*LoginInfo) (*Lget, error)
}

func NewLget() RapeHandler {
	return &Lget{}
}

type apis struct {
	EntryPoint      url.URL
	manualUrl       url.URL
	loginUrl        url.URL
	dataUrl         url.URL
	jobStateUrl     url.URL
	downloadFileUrl url.URL
}

type Lget struct {
	LgetResp http.Response
	apis
}

func (info *Lget) Login(loginInfo *LoginInfo) (*Lget, error) {
	apiUrl := &url.URL{}
	apiUrl.Scheme = "https"
	apiUrl.Host = fmt.Sprintf("%s-api.l-gate.net", loginInfo.Host)
	// url確認用
	pseudoUrl := *apiUrl
	pseudoUrl.Path = "auth/check-logged-in"
	fmt.Println(pseudoUrl.String())
	// urlが存在しているかテスト
	resp, err := http.Get(pseudoUrl.String())

	if err != nil {
		err = fmt.Errorf("login '%s' error: %w", pseudoUrl.String(), err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("login error - GET '%s' reponse statuscode %d", pseudoUrl.String(), resp.StatusCode)
		return nil, err
	}

	lget := &Lget{}
	lget.LgetResp = *resp
	apiUrl.Path = "control"
	lget.EntryPoint = *apiUrl

	lget.loginUrl = lget.EntryPoint
	lget.loginUrl.Path = fmt.Sprintf("%s/auth/login", lget.EntryPoint.Path)

	// ログインしてみるテスト
	logined, errors := lget.knock(loginInfo)
	if errors != nil {
		panic(err)
	}
	fmt.Println(logined.StatusCode)
	// 1.cookieを取得
	// 2.ログイン n回チャレンジ
	// 3.確認
	return lget, nil
}

func (lget *Lget) knock(info *LoginInfo) (resp *http.Response, errors []error) {
	payload := &LgateLoginInfo{
		LoginId:  info.AdminId,
		Password: info.AdminPw,
	}
	// ログインはID/PWをjsonで投げる形なので
	loginInfoJson, err := json.Marshal(&payload)
	if err != nil {
		err = fmt.Errorf("json marshal error: %w", err)
		return nil, []error{err}
	}
	// 時間を置いて3回チャレンジ
	interval := 5
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Duration(interval) * time.Second)
		loginResp, err := http.Post(lget.loginUrl.String(), "application/json", bytes.NewBuffer(loginInfoJson))
		if err != nil || loginResp.StatusCode != 200 {
			err = fmt.Errorf("login error: %w, statuscode: %d", err, loginResp.StatusCode)
			errors = append(errors, err)
			continue
		}
		defer loginResp.Body.Close()
		respBody, err := io.ReadAll(loginResp.Body)
		if err != nil {
			err = fmt.Errorf("read response body error: %w", err)
			errors = append(errors, err)
			continue
		}
		var loginedResp LoginedResp
		if err := json.Unmarshal(respBody, &loginedResp); err != nil {
			err = fmt.Errorf("unmarshall response error: %w", err)
			errors = append(errors, err)
			continue
		}
		if loginedResp.Code == 200 {
			return loginResp, nil
		} else {
			err = fmt.Errorf("login status code not 200")
			errors = append(errors, err)
			continue
		}
	}
	return nil, errors
}
