package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/halllllll/lget"
	"github.com/spf13/viper"
)

var cd string

const (
	userDataFolderName string = "user_csvs"
	userLogFolderName  string = "user_action_csvs"
)

type Conf struct {
	UserLogStartAtUnixTime int `mapstructure:"LGET_ALLUSER_ACTIONLOG_STARTATUNIXTIME"`
	UserLogEndAtUnixTime   int `mapstructure:"LGET_ALLUSER_ACTIONLOG_ENDATUNIXTIME"`
	UserLogBetweenMinutes  int `mapstructure:"LGET_ALLUSER_ACTIONLOG_BETWEEN_MINUTES"`
}

var cnf *Conf

func loadConf(c *Conf) (err error) {
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&c)
	return
}

func init() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// load env file
	viper.AddConfigPath(".")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.SetConfigFile(filepath.Join(cwd, ".env"))
	if err := loadConf(cnf); err != nil {
		panic(err)
	}
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// load credential file
	loginInfoJson, err := Env.ReadFile("secret.json")
	if err != nil {
		panic(err)
	}
	var lijs loginInfoJsonStruct
	err = json.Unmarshal(loginInfoJson, &lijs)
	if err != nil {
		panic(err)
	}
	loginInfo = &lget.LoginInfo{
		Host:    lijs.Host,
		AdminId: lijs.AdminId,
		AdminPw: lijs.AdminPw,
	}

	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cd = curDir
	usersDataCsvPath := filepath.Join(cd, userDataFolderName)
	if _, err := os.Stat(usersDataCsvPath); os.IsNotExist(err) {
		if err := os.MkdirAll(usersDataCsvPath, 0755); err != nil {
			panic(err)
		}
	}
	usersActionLogCsvPath := filepath.Join(cd, userLogFolderName)
	if _, err := os.Stat(usersActionLogCsvPath); os.IsNotExist(err) {
		if err := os.MkdirAll(usersActionLogCsvPath, 0755); err != nil {
			panic(err)
		}
	}
}

func someHeavyProcessMock(t chan time.Time) {
	// なんか重い処理のフリ
	count := 0
	for {
		fmt.Println("some heavy process mock start...")
		time.Sleep(5 * time.Minute)
		count++
		t <- time.Now()
		fmt.Printf("heavy process done!!! count: %d\n", count)
	}
}

func runGetUser(loginInfo *lget.LoginInfo, interval time.Duration, result chan []byte) {
	for {
		start := time.Now()
		// まずログインを済ませる
		l_get := lget.NewLget()

		opened_l_get, err := l_get.Login(loginInfo)
		if err != nil {
			panic(err)
		}

		// ユーザーデータを全部取得するAPIを叩く
		downloadFileUrl, err := opened_l_get.GetAllUser()
		if err != nil {
			panic(err)
		}
		fmt.Printf("download file url: %s\n", downloadFileUrl)

		rawData, err := opened_l_get.Download(downloadFileUrl)
		if err != nil {
			panic(err)
		}
		result <- rawData

		end := time.Now()

		execDif := end.Sub(start)
		if execDif < interval {
			dif := interval - execDif
			next := end.Add(time.Duration(dif.Seconds()) * time.Second).Format("2006/01/02/15:04:05")
			fmt.Printf("sleep until %s\n", next)
			time.Sleep(dif)
		} else {
			fmt.Printf("omg!! over time... duration: %d\n gogogo start now!\n", execDif-interval)
		}
	}
}

func runGetAllLog(loginInfo *lget.LoginInfo, result chan actionLogData) {
	for {
		// .envからstartAtUnixTimeとendAtUnixTimeとbetweenminutesを読み込むのにviperを使う
		envRawVal := viper.GetInt64("LGET_ALLUSER_ACTIONLOG_STARTATUNIXTIME")
		startAtUnixTime := envRawVal

		envRawVal = viper.GetInt64("LGET_ALLUSER_ACTIONLOG_ENDATUNIXTIME")
		endAtUnixTime := envRawVal

		envRawVal = viper.GetInt64("LGET_ALLUSER_ACTIONLOG_BETWEEN_MINUTES")
		betweenInterval := envRawVal
		// endAtUnixTimeが今よりも先だった場合は待つ
		if time.Now().Before(time.Unix(endAtUnixTime, 0)) {
			fmt.Printf("suspend until %s\n", time.Unix(endAtUnixTime, 0))
			<-time.After(time.Until(time.Unix(endAtUnixTime, 0)))
			fmt.Printf("start!")
		}

		start := time.Now()
		// 全部取得するAPIを叩く
		// まずログインを済ませる
		l_get := lget.NewLget()

		opened_l_get, err := l_get.Login(loginInfo)
		if err != nil {
			panic(err)
		}

		// ex 2022-08-21 10:00:00 -> 1661043600
		// ex 2022-08-21 13:00:00 -> 1661054400
		downloadFileUrl, err := opened_l_get.GetLog(int(startAtUnixTime), int(endAtUnixTime))
		if err != nil {
			panic(err)
		}
		fmt.Printf("download file url: %s\n", downloadFileUrl)

		rawData, err := opened_l_get.Download(downloadFileUrl)
		if err != nil {
			panic(err)
		}
		ret := actionLogData{
			fileRawData: rawData,
			from:        int(startAtUnixTime),
			to:          int(endAtUnixTime),
		}
		result <- ret

		end := time.Now()
		fmt.Printf("execution time: %s\n", end.Sub(start))
		// 次回の時間を設定
		// envを上書きする
		startAtUnixTime = endAtUnixTime + 1
		endAtUnixTime = time.Unix(endAtUnixTime, 0).Add(time.Duration(betweenInterval) * time.Minute).Unix()

		viper.Set("LGET_ALLUSER_ACTIONLOG_STARTATUNIXTIME", startAtUnixTime)
		viper.Set("LGET_ALLUSER_ACTIONLOG_ENDATUNIXTIME", endAtUnixTime)
		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}

	}
}

func saveFile(data []byte, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

//go:embed secret.json
var Env embed.FS

var loginInfo *lget.LoginInfo

type loginInfoJsonStruct struct {
	Host    string `json:"host"`
	AdminId string `json:"login_id"`
	AdminPw string `json:"password"`
}

type actionLogData struct {
	fileRawData []byte
	from        int
	to          int
}

func main() {
	// 開始時間設定

	readyAction := time.Date(2022, time.December, 4, 14, 20, 0, 0, time.Local)

	<-time.After(time.Until(readyAction))

	// ユーザーデータ取得用ゴルーチン
	userResult := make(chan []byte)
	usersDataCsvPath := filepath.Join(cd, userDataFolderName)
	// インターバル指定して無限に取得開始
	go runGetUser(loginInfo, time.Minute*30, userResult)

	// ユーザー履歴取得用ゴルーチン
	userLogResult := make(chan actionLogData)
	userLogCsvPath := filepath.Join(cd, userLogFolderName)
	// startatunixtime, endatunixtimeは.envファイルから読み出すことにする
	go runGetAllLog(loginInfo, userLogResult)

	// mock (for feature api)
	pseudoResult := make(chan time.Time)
	go someHeavyProcessMock(pseudoResult)

	for {
		select {
		case userData := <-userResult:
			saveFileName := time.Now().Format("user_2006_01_02_150405.csv")
			err := saveFile(userData, filepath.Join(usersDataCsvPath, saveFileName))
			if err != nil {
				panic(err)
			}
		case userLogData := <-userLogResult:
			saveFileName := fmt.Sprintf("useractionlog_%s__%s.csv", time.Unix(int64(userLogData.from), 0).Format("2006_01_02_150405"), time.Unix(int64(userLogData.to), 0).Format("2006_01_02_150405"))
			err := saveFile(userLogData.fileRawData, filepath.Join(userLogCsvPath, saveFileName))
			if err != nil {
				panic(err)
			}

		case doneTime := <-pseudoResult:
			fmt.Printf("haha over %s, see you next Time!\n ", doneTime)
		}
	}
}
