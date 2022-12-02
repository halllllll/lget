package lget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/corpix/uarand"
)

const (
	CONTROLSESSID       = "CONTROLSESSID"
	LGATEORIGINALCOOKIE = "_ga_EP4PNHSVYP=TISATONISHIKIGI; _ga=TAKINAINOUE; "
)

func init() {

}

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
	loginUrl        url.URL
	helloUrl        url.URL
	dataUrl         url.URL
	jobStateUrl     url.URL
	downloadFileUrl url.URL
}

// 各種APIのエンドポイントをドメインから作成作成
func (a *apis) prepareApiUrls() {
	a.EntryPoint.Path = "control"
	a.loginUrl = a.EntryPoint
	a.loginUrl.Path = fmt.Sprintf("%s/auth/login", a.EntryPoint.Path)
	a.helloUrl = a.EntryPoint
	a.helloUrl.Path = fmt.Sprintf("%s/manual/get", a.EntryPoint.Path)
	a.dataUrl = a.EntryPoint
	a.dataUrl.Path = fmt.Sprintf("%s/action-log/download-csv-total", a.EntryPoint.Path)
	a.jobStateUrl = a.EntryPoint
	a.jobStateUrl.Path = fmt.Sprintf("%s/job-state/view", a.EntryPoint.Path)
	a.downloadFileUrl = a.EntryPoint
	a.downloadFileUrl.Path = fmt.Sprintf("%s/file/view", a.EntryPoint.Path)
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
	lget.EntryPoint = *apiUrl
	// assemble api urls
	lget.apis.prepareApiUrls()
	// ログインテスト
	logined, errors := lget.knock(loginInfo)
	if errors != nil {
		fmt.Printf("errors length: %d\n", len(errors))
		for _, err := range errors {
			fmt.Printf("%w\n", err)
		}
		return nil, fmt.Errorf("some error occured: %w", errors)
	}
	// 1. ログイン後のresponseからcookieを取得するだけ keyは決め打ち
	cookie, err := lget.doorBell(logined, CONTROLSESSID)
	if err != nil {
		panic(err)
	}
	fmt.Println(cookie)
	// 2.cookieを使って(おそらくセッション毎の)uuidを取得
	resultUuid, err := lget.chaim(cookie)
	if err != nil {
		panic(err)
	}
	fmt.Println(resultUuid)
	// 3.確認
	return lget, nil
}

// 所与のデータでログインできるか確認（数回チャレンジ）
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

// extract cookie by certainly cookie name
func (lget *Lget) doorBell(resp *http.Response, cookieName string) (string, error) {
	parser := &http.Request{Header: http.Header{"Cookie": resp.Header["Set-Cookie"]}}
	cookie, err := parser.Cookie(cookieName)
	if err != nil {
		err = fmt.Errorf("parse cookie error: %w", err)
		return "", err
	}
	return cookie.Value, nil
}

func (lget *Lget) chaim(cookie string) (resultUuid string, err error) {
	// なぜかこのurlで最初に飛ばないとresultを得られなかった（ブラウザでも最初にGETを飛ばしているっぽい）
	req, err := http.NewRequest(http.MethodGet, lget.helloUrl.String(), nil)
	if err != nil {
		err = fmt.Errorf("create request error: %w", err)
		return
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", CONTROLSESSID, cookie))
	client := &http.Client{}
	dataResp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("request error: %w", err)
		return
	}
	defer dataResp.Body.Close()

	data, err := io.ReadAll(dataResp.Body)
	if err != nil {
		err = fmt.Errorf("read response body erro: %w", err)
		return
	}

	var getDataResp GetDataResp
	if err = json.Unmarshal(data, &getDataResp); err != nil {
		err = fmt.Errorf("unmarhal response error: %w", err)
		return
	}
	if getDataResp.Code != 200 {
		err = fmt.Errorf("statuscode: %d", getDataResp.Code)
		return
	}
	resultUuid = getDataResp.Result.UUID
	return
}
