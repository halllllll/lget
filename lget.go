package lget

import (
	"fmt"
	"net/http"
	"net/url"
)

type loginInfo struct {
	host    string
	adminId string
	adminPw string
}

type RapeHandler interface {
	Login() (*Lget, error)
	SetHost(string) error
	SetAdminId(string) error
	SetAdminPw(string) error
}

func NewLget() RapeHandler {
	return &loginInfo{}
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

func (info *loginInfo) SetHost(host string) error {
	info.host = host
	return nil
}

func (info *loginInfo) SetAdminId(adminId string) error {
	info.adminId = adminId
	return nil
}
func (info *loginInfo) SetAdminPw(adminPw string) error {
	info.adminPw = adminPw
	return nil
}

func (info *loginInfo) Login() (*Lget, error) {
	apiUrl := &url.URL{}
	apiUrl.Scheme = "https"
	// apiUrl.Host = fmt.Sprintf("%s.l-gate.net", info.Host)
	apiUrl.Host = fmt.Sprintf("%s-api.l-gate.net", info.host)
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

	// ログインしてみるテスト
	// 1.cookieを取得
	// 2.ログイン n回チャレンジ
	// 3.確認
	return lget, nil
}

// func (lget *Lget) knock(info *LoginInfo) (cookie string, err error) {
// 	type Payload struct {
// 		loggin_id string
// 		password  string
// 	}
// 	// ログインはID/PWをjsonで投げる形なので
// 	paylaod := Payload{
// 		loggin_id: info.AdminId,
// 		password:  info.AdminPw,
// 	}
// 	loginInfoJson, err := json.Marshal(&paylaod)
// 	if err != nil {
// 		err = fmt.Errorf("json marshal error: %w", err)
// 		return "", err
// 	}
// 	// 時間を置いて3回チャレンジ
// 	interval := 5
// 	for i := 1; i <= 3; i++ {
// 		time.Sleep(time.Duration(interval) * time.Second)
// 		loginResp, err := http.Post(lget.manualUrl.String(), "application/json", bytes.NewBuffer(loginInfoJson))
// 		if err != nil || loginResp.StatusCode != 200 {
// 			// err = fmt.Errorf("login error: %w, statuscode: %d", err, loginResp.StatusCode)
// 			// utils.ErrLog.Println(err)
// 			continue
// 		}
// 		// defer loginResp.Body.Close()
// 		respBody, err := io.ReadAll(loginResp.Body)
// 		if err != nil {
// 			// err = fmt.Errorf("read response body error: %w", err)
// 			// utils.ErrLog.Println(err)
// 			continue
// 		}
// 		var loginedResp *typefile.LoginedResp
// 		if err := json.Unmarshal(respBody, &loginedResp); err != nil {
// 			// err = fmt.Errorf("unmarshall response error: %w", err)
// 			// utils.ErrLog.Println(err)
// 			continue
// 		}
// 		if loginedResp.Code == 200 {
// 			// return loginResp, nil
// 		} else {
// 			// err = fmt.Errorf("login status code not 200")
// 			// utils.ErrLog.Println(err)
// 			continue
// 		}
// 	}
// 	return nil, err
// }
