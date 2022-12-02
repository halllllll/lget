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
	"github.com/halllllll/golog"
)

const (
	CONTROLSESSID = "CONTROLSESSID"
	LGATEGACOOKIE = "_ga_EP4PNHSVYP=TISATONISHIKIGI; _ga=TAKINAINOUE; "
)

func init() {
	golog.LoggingSetting("lget.log")
}

type LoginInfo struct {
	Host    string
	AdminId string
	AdminPw string
}

// ログイン通る前
type LgetHandler interface {
	Login(*LoginInfo) (OpenedLgetHandler, error)
}

// ログイン後
type OpenedLgetHandler interface {
	GetLog(startUnixTime int, endUnixTime int) (string, error)
}

func NewLget() LgetHandler {
	return &Lget{}
}

type apis struct {
	EntryPoint      url.URL
	loginUrl        url.URL
	helloUrl        url.URL
	dataUrl         url.URL // first contact
	jobStateUrl     url.URL // toritate style
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
	cookie     string
	resultUuid string
}

func (info *Lget) Login(loginInfo *LoginInfo) (OpenedLgetHandler, error) {
	golog.InfoLog.Println("login challenge...")
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
		golog.ErrLog.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("login error - GET '%s' reponse statuscode %d", pseudoUrl.String(), resp.StatusCode)
		golog.ErrLog.Println(err)
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
		return nil, fmt.Errorf("some error occured: please refer to logfile")
	}
	// 1. ログイン後のresponseからcookieを取得するだけ keyは決め打ち
	cookie, err := lget.doorBell(logined, CONTROLSESSID)
	if err != nil {
		golog.ErrLog.Println(err)
		return nil, err
	}
	lget.cookie = cookie
	// 2.cookieを使って(おそらくセッション毎の)uuidを取得
	resultUuid, err := lget.chaim(cookie)
	if err != nil {
		golog.ErrLog.Println(err)
		return nil, err
	}
	lget.resultUuid = resultUuid

	golog.InfoLog.Println("login successed!")
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
		golog.ErrLog.Println(err)
		return nil, []error{err}
	}
	// 時間を置いて3回チャレンジ
	interval := 5
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Duration(interval) * time.Second)
		loginResp, err := http.Post(lget.loginUrl.String(), "application/json", bytes.NewBuffer(loginInfoJson))
		if err != nil || loginResp.StatusCode != 200 {
			err = fmt.Errorf("login error: %w, statuscode: %d", err, loginResp.StatusCode)
			golog.ErrLog.Println(err)
			errors = append(errors, err)
			continue
		}
		defer loginResp.Body.Close()
		respBody, err := io.ReadAll(loginResp.Body)
		if err != nil {
			err = fmt.Errorf("read response body error: %w", err)
			golog.ErrLog.Println(err)
			errors = append(errors, err)
			continue
		}
		var loginedResp LoginedResp
		if err := json.Unmarshal(respBody, &loginedResp); err != nil {
			err = fmt.Errorf("unmarshall response error: %w", err)
			golog.ErrLog.Println(err)
			errors = append(errors, err)
			continue
		}
		if loginedResp.Code == 200 {
			return loginResp, nil
		} else {
			err = fmt.Errorf("login status code not 200")
			golog.ErrLog.Println(err)
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
		golog.ErrLog.Println(err)
		return "", err
	}
	return cookie.Value, nil
}

// challenge getting result uuid
func (lget *Lget) chaim(cookie string) (resultUuid string, err error) {
	// なぜかこのurlで最初に飛ばないとresultを得られなかった（ブラウザでも最初にGETを飛ばしているっぽい）
	req, err := http.NewRequest(http.MethodGet, lget.helloUrl.String(), nil)
	if err != nil {
		err = fmt.Errorf("create request error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", CONTROLSESSID, cookie))
	client := &http.Client{}
	dataResp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("request error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	defer dataResp.Body.Close()

	data, err := io.ReadAll(dataResp.Body)
	if err != nil {
		err = fmt.Errorf("read response body erro: %w", err)
		golog.ErrLog.Println(err)
		return
	}

	var getDataResp GetDataResp
	if err = json.Unmarshal(data, &getDataResp); err != nil {
		err = fmt.Errorf("unmarhal response error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	if getDataResp.Code != 200 {
		err = fmt.Errorf("statuscode: %d", getDataResp.Code)
		golog.ErrLog.Println(err)
		return
	}
	resultUuid = getDataResp.Result.UUID
	return
}

// 以下データ取得系API
func (lget *Lget) GetLog(startUnixTime, endUnixTime int) (string, error) {
	if startUnixTime >= endUnixTime {
		err := fmt.Errorf("end unixtime should be later than start unixtime")
		golog.ErrLog.Println(err)
		return "", err
	}
	golog.InfoLog.Println("start: GET LOGS FOR ALL KINDS.")
	//　全種類の履歴取得用URL構築
	// url.URLでちゃんと構築したほうが行儀がいいかもしれない
	startUrl := fmt.Sprintf("%s?start_at=%d&end_at=%d&time_unit=hour&scope=tenant&action=&response_all=1&encoding=utf8", lget.dataUrl.String(), startUnixTime, endUnixTime)

	golog.ErrLog.Printf("start url: %s\n", startUrl)
	jobUuid, err := rattlingKnob(startUrl, lget.cookie)
	if err != nil {
		err = fmt.Errorf("get data (firstcontact) error: %w", err)
		golog.ErrLog.Println(err)
		return "", err
	}
	jobUrl, err := url.JoinPath(lget.jobStateUrl.String(), jobUuid)
	if err != nil {
		err = fmt.Errorf("compose url path error: %w", err)
		return "", err
	}
	golog.InfoLog.Printf("job url: %s\n", jobUrl)

	downloadFileUuid, err := brokenBuzzer(jobUrl, lget.cookie)
	if err != nil {
		return "", err
	}
	golog.InfoLog.Printf("file download uuid: %s\n", downloadFileUuid)

	// ダウンロードリンクを構築
	downloadFileUrl, err := url.JoinPath(lget.downloadFileUrl.String(), downloadFileUuid)
	if err != nil {
		err = fmt.Errorf("compose download url link error: %w", err)
		return "", err
	}

	return downloadFileUrl, err
}

// first contact for get data (from url)
func rattlingKnob(targetUrl, cookie string) (jobUuid string, err error) {
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		err = fmt.Errorf("create new request error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Cookie", fmt.Sprintf("%s%s=%s", LGATEGACOOKIE, CONTROLSESSID, cookie))
	client := &http.Client{}
	firstContactResp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("request error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	defer firstContactResp.Body.Close()
	firstContactData, err := io.ReadAll(firstContactResp.Body)
	if err != nil {
		err = fmt.Errorf("response read error: %w", err)
		golog.ErrLog.Println(err)
		return
	}
	var firstContact GetDataResp
	if err = json.Unmarshal(firstContactData, &firstContact); err != nil {
		err = fmt.Errorf("unmasharl error: %w", err)
		golog.ErrLog.Println(err)
		return
	}

	if !firstContact.IsSuccessful {
		err = fmt.Errorf("first contact error: statuscode %d\nresult: %#v", firstContact.Code, firstContact)
		golog.ErrLog.Println(err)
		return
	}

	return firstContact.Result.UUID, nil
}

func brokenBuzzer(targetUrl, cookie string) (downloadFileUuid string, err error) {
	if targetUrl == "" || cookie == "" {
		err = fmt.Errorf("both url and cookie DONT shoud be empty")
		golog.ErrLog.Println(err)
		return "", err
	}
	count := 0
	startTime := time.Now().Format("2006-01-02 15:04:05")

	golog.InfoLog.Printf("GO! start at %s\n", startTime)
	for {
		count += 1
		req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
		if err != nil {
			err = fmt.Errorf("create request error: %w ", err)
			golog.ErrLog.Println(err)
			return "", err
		}
		req.Header.Set("Cookie", fmt.Sprintf("%s%s=%s", LGATEGACOOKIE, CONTROLSESSID, cookie))

		client := &http.Client{}
		respRow, err := client.Do(req)
		if err != nil {
			err = fmt.Errorf("request error: %w", err)
			golog.ErrLog.Println(err)
			return "", err
		}
		if respRow.StatusCode != 200 {
			err = fmt.Errorf("status code: %d", respRow.StatusCode)
			golog.ErrLog.Println(err)
			return "", err
		}
		defer respRow.Body.Close()

		// レスポンスを読める形(バイトデータ)にする
		data, err := io.ReadAll(respRow.Body)
		if err != nil {
			err = fmt.Errorf("read response error: %w", err)
			golog.ErrLog.Println(err)
			return "", err
		}
		// バイトデータはjsonなので
		var curData GetDataResp
		if err := json.Unmarshal(data, &curData); err != nil {
			err = fmt.Errorf("unmarshall error: %w", err)
			golog.ErrLog.Println(err)
			return "", err
		}
		result := curData.Result
		msg := result.Message
		downloadFileUuid = result.Result.FileUUID
		if result.IsSuccess {
			golog.InfoLog.Printf("Done! %s\n", msg)
			break
		} else if msg == "CSVエクスポートに失敗しました。" {
			err = fmt.Errorf("csv export error: %s", msg)
			golog.ErrLog.Println(err)
			golog.ErrLog.Printf("result: %#v\n", result)
			return "", err
		} else if msg != "学習ログCSVエクスポートキューが実行待ちです。" {
			// unknown error message(even 2022/12/02)
			err = fmt.Errorf("unkown error : %s", msg)
			golog.ErrLog.Println(err)
			return "", err
		}
		golog.InfoLog.Printf("%d - time: %s - %s\n", count, time.Now().Format("2006-01-02 15:04:05"), msg)

		time.Sleep(15 * time.Second) // 15 sec is official interval (at least on browser)
	}
	return
}
