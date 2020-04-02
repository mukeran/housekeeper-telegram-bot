package methods

import (
	"HouseKeeperBot/database"
	"HouseKeeperBot/modules/cccat/models"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	cccatBaseUrl                       = "https://cccat.io"
	cccatIndexUrl                      = cccatBaseUrl + "/user/index.php"
	cccatLoginUrl                      = cccatBaseUrl + "/user/_login.php"
	cccatCheckinUrl                    = cccatBaseUrl + "/user/_checkin.php"
	cccatCodeLoginEmailOrPasswordWrong = "0"
	cccatCodeLoginSuccessful           = "1"
)

var (
	ErrSigned                      = errors.New("the account has been signed today")
	ErrNoSuchAccount               = errors.New("no such account")
	ErrInsufficientCookie          = errors.New("insufficient cookie")
	ErrWrongAccountEmailOrPassword = errors.New("wrong account email or password")
	ErrInvalidCookie               = errors.New("invalid cookie")
	ErrLoginFailed                 = errors.New("login failed")
	ErrUnknown                     = errors.New("unknown error")
	regexpSigned                   = regexp.MustCompile(`You have already checked in`)
	regexpSuccessful               = regexp.MustCompile(`Get (0|[1-9][0-9]*)MB transfer`)
	regexpSuccessfulDouble         = regexp.MustCompile(`Get (0|[1-9][0-9]*)MB \+(0|[1-9][0-9]*)MB transfer`)
	regexpRemaining                = regexp.MustCompile(
		`Remaining Transfer: (0|[1-9][0-9]*(?:\.(?:0|[1-9][0-9]*)?))GB`)
)

type cccatLoginResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type cccatSignResult struct {
	Msg string `json:"msg"`
}

func SignWithAccountID(accountID uint) (got uint, err error) {
	account := GetAccountByID(accountID)
	if account == nil {
		return 0, ErrNoSuchAccount
	}
	return Sign(account)
}

func QueryRemainingTransferWithAccountID(accountID uint) (remaining float64, err error) {
	account := GetAccountByID(accountID)
	if account == nil {
		return 0, ErrNoSuchAccount
	}
	return QueryRemainingTransfer(account)
}

func getCookie(email, password string) (uid, userPwd string, err error) {
	form := url.Values{}
	form.Set("email", email)
	form.Set("passwd", password)
	form.Set("remember_me", "week")
	req, err := http.NewRequest(http.MethodPost, cccatLoginUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result cccatLoginResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}
	if result.Code == cccatCodeLoginEmailOrPasswordWrong {
		return "", "", ErrWrongAccountEmailOrPassword
	} else if result.Code != cccatCodeLoginSuccessful {
		log.Printf("Unknown login code. Response: %v", data)
		return "", "", ErrLoginFailed
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "uid" {
			uid = cookie.Value
		} else if cookie.Name == "user_pwd" {
			userPwd = cookie.Value
		}
	}
	if uid == "" || userPwd == "" {
		return "", "", ErrInsufficientCookie
	}
	return uid, userPwd, nil
}

func Sign(account *models.Account) (got uint, err error) {
	uid, userPwd := account.CookieUID, account.CookieUserPwd
	if account.HasLoginCredentials {
		uid, userPwd, err = getCookie(account.Email, account.Password)
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(http.MethodGet, cccatCheckinUrl, nil)
	if err != nil {
		return
	}
	req.AddCookie(&http.Cookie{Name: "uid", Value: uid})
	req.AddCookie(&http.Cookie{Name: "user_pwd", Value: userPwd})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if bytes.Contains(data, []byte("login-page")) {
		return 0, ErrInvalidCookie
	}
	var result cccatSignResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}
	defer func(accountID uint, raw string, got *uint, err *error) {
		signLog := models.SignLog{
			AccountID: accountID,
			Status: func() uint {
				switch *err {
				case nil:
					return models.SignStatusSuccessful
				case ErrSigned:
					return models.SignStatusSigned
				default:
					return models.SignStatusFailed
				}
			}(),
			GotTransfer: *got,
			Raw:         raw,
		}
		tx := database.Db.Begin()
		defer tx.RollbackUnlessCommitted()
		if v := tx.Create(&signLog); v.Error != nil {
			log.Panic(v.Error)
		}
		if v := tx.Commit(); v.Error != nil {
			log.Panic(v.Error)
		}
	}(account.ID, result.Msg, &got, &err)
	if regexpSigned.MatchString(result.Msg) {
		return 0, ErrSigned
	} else if regexpSuccessful.MatchString(result.Msg) {
		match := regexpSuccessful.FindStringSubmatch(result.Msg)
		_got, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return 0, err
		}
		got = uint(_got)
	} else if regexpSuccessfulDouble.MatchString(result.Msg) {
		match := regexpSuccessfulDouble.FindStringSubmatch(result.Msg)
		_got1, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return 0, err
		}
		_got2, err := strconv.ParseUint(match[2], 10, 32)
		if err != nil {
			return 0, err
		}
		got = uint(_got1 + _got2)
	} else {
		log.Printf("Unknown sign response: %v", result.Msg)
		return 0, ErrUnknown
	}
	return
}

func QueryRemainingTransfer(account *models.Account) (remaining float64, err error) {
	uid, userPwd := account.CookieUID, account.CookieUserPwd
	if account.HasLoginCredentials {
		uid, userPwd, err = getCookie(account.Email, account.Password)
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(http.MethodGet, cccatIndexUrl, nil)
	if err != nil {
		return
	}
	req.AddCookie(&http.Cookie{Name: "uid", Value: uid})
	req.AddCookie(&http.Cookie{Name: "user_pwd", Value: userPwd})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	match := regexpRemaining.FindSubmatch(data)
	if match == nil {
		log.Printf("Unknown query remaining transfer response: %v", data)
		return 0, ErrUnknown
	}
	remaining, err = strconv.ParseFloat(string(match[1]), 64)
	return
}
