package management

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Authing/authing-go-sdk/lib/constant"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/Authing/authing-go-sdk/lib/util/cacheutil"
	"github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// Client is a client for interacting with the GraphQL API of `Authing`
type Client struct {
	HttpClient *http.Client
	userPoolId string
	secret     string
	Host       string

	// Log is called with various debug information.
	// To log to standard out, use:
	//  client.Log = func(s string) { log.Println(s) }
	Log func(s string)
}

func NewClient(userPoolId string, secret string, host ...string) *Client {
	var clientHost string
	if len(host) == 0 {
		clientHost = constant.CoreAuthingDefaultUrl
	} else {
		clientHost = host[0]
	}
	c := &Client{
		userPoolId: userPoolId,
		secret:     secret,
		Host:       clientHost,
	}
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
		_, err := GetAccessToken(c)
		if err != nil {
			return nil
		}
	}
	return c
}

func NewClientWithError(userPoolId string, secret string, host ...string) (*Client, error) {
	if userPoolId == "" {
		return nil, errors.New("请填写 userPoolId 参数")
	}
	if secret == "" {
		return nil, errors.New("请填写 secret 参数")
	}
	var clientHost string
	if len(host) == 0 {
		clientHost = constant.CoreAuthingDefaultUrl
	} else {
		clientHost = host[0]
	}
	c := &Client{
		userPoolId: userPoolId,
		secret:     secret,
		Host:       clientHost,
	}
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{Timeout: 30 * time.Second}
		_, err := GetAccessToken(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) SendHttpRequest(url string, method string, query string, variables map[string]interface{}) ([]byte, error) {
	var req *http.Request
	if method == constant.HttpMethodGet {
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		if len(variables) > 0 {
			q := req.URL.Query()
			for key, value := range variables {
				q.Add(key, fmt.Sprintf("%v", value))
			}
			req.URL.RawQuery = q.Encode()
		}
	} else {
		in := struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables,omitempty"`
		}{
			Query:     query,
			Variables: variables,
		}
		var buf bytes.Buffer
		var err error
		if query == constant.StringEmpty {
			err = json.NewEncoder(&buf).Encode(variables)
		} else {
			err = json.NewEncoder(&buf).Encode(in)
		}
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, &buf)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	}

	// 增加header选项
	if !strings.HasPrefix(query, "query accessToken") {
		token, err := GetAccessToken(c)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("x-authing-userpool-id", ""+c.userPoolId)
	req.Header.Add("x-authing-request-from", constant.SdkType)
	req.Header.Add("x-authing-sdk-version", constant.SdkVersion)
	req.Header.Add("x-authing-app-id", ""+constant.AppId)
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) SendHttpRestRequest(url string, method string, variables map[string]interface{}) ([]byte, error) {
	var req *http.Request
	if method == constant.HttpMethodGet {
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		if variables != nil && len(variables) > 0 {
			q := req.URL.Query()
			for key, value := range variables {
				q.Add(key, fmt.Sprintf("%v", value))
			}
			req.URL.RawQuery = q.Encode()
		}

	} else {

		var buf bytes.Buffer
		var err error
		if variables != nil {
			err = json.NewEncoder(&buf).Encode(variables)

		}
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, &buf)
		req.Header.Add("Content-Type", "application/json")
	}

	token, _ := GetAccessToken(c)
	req.Header.Add("Authorization", "Bearer "+token)

	req.Header.Add("x-authing-userpool-id", ""+c.userPoolId)
	req.Header.Add("x-authing-request-from", constant.SdkType)
	req.Header.Add("x-authing-sdk-version", constant.SdkVersion)
	req.Header.Add("x-authing-app-id", ""+constant.AppId)
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, nil
}

func (c *Client) httpGet(url string, client *http.Client) (string, error) {
	reqest, err := http.NewRequest(constant.HttpMethodGet, c.Host+url, nil)
	if err != nil {
		return "", err
	}

	// 增加header选项
	token, _ := GetAccessToken(c)
	reqest.Header.Add("Authorization", "Bearer "+token)
	reqest.Header.Add("x-authing-userpool-id", ""+c.userPoolId)
	reqest.Header.Add("x-authing-request-from", constant.SdkType)
	reqest.Header.Add("x-authing-sdk-version", constant.SdkVersion)
	reqest.Header.Add("x-authing-app-id", ""+constant.AppId)

	resp, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := string(body)
	return result, nil
}

func (c *Client) SendHttpRequestV2(url string, method string, query string, variables map[string]interface{}) ([]byte, error) {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return nil, err
	}
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(url)
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("x-authing-userpool-id", ""+c.userPoolId)
	req.Header.Add("x-authing-request-from", constant.SdkType)
	req.Header.Add("x-authing-sdk-version", constant.SdkVersion)
	req.Header.Add("x-authing-app-id", ""+constant.AppId)
	req.Header.SetMethod(method)
	req.SetBody(buf.Bytes())

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)
	body := resp.Body()
	return body, err
}

func QueryAccessToken(client *Client) (*model.AccessTokenRes, error) {
	type Data struct {
		AccessToken model.AccessTokenRes `json:"accessToken"`
	}
	type Result struct {
		Data Data `json:"data"`
	}

	variables := map[string]interface{}{
		"userPoolId": client.userPoolId,
		"secret":     client.secret,
	}

	b, err := client.SendHttpRequest(client.Host+constant.CoreAuthingGraphqlPath, constant.HttpMethodPost, constant.AccessTokenDocument, variables)
	if err != nil {
		return nil, err
	}
	err = checkError(b)
	if err != nil {
		return nil, err
	}
	var r Result
	if b != nil {
		json.Unmarshal(b, &r)
	}
	return &r.Data.AccessToken, nil
}

func checkError(b []byte) error {
	json, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}
	repErrors, exist := json.CheckGet("errors")
	if !exist {
		return nil
	}
	result, err := repErrors.Array()
	if err != nil {
		return err
	}
	if result != nil && len(result) > 0 {
		reason, err := json.Get("errors").GetIndex(0).Get("message").Get("message").String()
		if err != nil {
			return err
		}
		return errors.New(reason)
	}
	return nil
}

func GetAccessToken(client *Client) (string, error) {
	// 从缓存获取token
	cacheToken, b := cacheutil.GetCache(constant.TokenCacheKeyPrefix + client.userPoolId)
	if b && cacheToken != nil {
		return cacheToken.(string), nil
	}
	token, err := QueryAccessToken(client)
	if err != nil {
		return "", err
	}
	var expire = 24 * time.Hour
	cacheutil.SetCache(constant.TokenCacheKeyPrefix+client.userPoolId, *token.AccessToken, expire)
	return *token.AccessToken, nil
}

// SendEmail
// 发送邮件
func (c *Client) SendEmail(email string, scene model.EnumEmailScene) (*model.CommonMessageAndCode, error) {

	b, err := c.SendHttpRequest(c.Host+constant.CoreAuthingGraphqlPath, http.MethodPost, constant.SendMailDocument,
		map[string]interface{}{"email": email, "scene": scene})
	if err != nil {
		return nil, err
	}
	var response = &struct {
		Data struct {
			SendMail model.CommonMessageAndCode `json:"sendEmail"`
		} `json:"data"`
		Errors []model.GqlCommonErrors `json:"errors"`
	}{}

	jsoniter.Unmarshal(b, &response)
	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message.Message)
	}
	return &response.Data.SendMail, nil
}

// CheckLoginStatusByToken
// 检测登录状态
func (c *Client) CheckLoginStatusByToken(token string) (*model.CheckLoginStatusResponse, error) {

	b, err := c.SendHttpRequest(c.Host+constant.CoreAuthingGraphqlPath, http.MethodPost, constant.CheckLoginStatusDocument,
		map[string]interface{}{"token": token})
	if err != nil {
		return nil, err
	}
	var response = &struct {
		Data struct {
			CheckLoginStatus model.CheckLoginStatusResponse `json:"checkLoginStatus"`
		} `json:"data"`
		Errors []model.GqlCommonErrors `json:"errors"`
	}{}

	jsoniter.Unmarshal(b, &response)
	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message.Message)
	}
	return &response.Data.CheckLoginStatus, nil
}

// IsPasswordValid
// 检测密码是否合法
func (c *Client) IsPasswordValid(password string) (*struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}, error) {

	url := fmt.Sprintf("%s/api/v2/password/check", c.Host)
	b, err := c.SendHttpRestRequest(url, http.MethodPost, map[string]interface{}{"password": password})
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Message string `json:"message"`
		Code    int64  `json:"code"`
		Data    struct {
			Valid   bool   `json:"valid"`
			Message string `json:"message"`
		} `json:"data"`
	}{}
	jsoniter.Unmarshal(b, &resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return &resp.Data, nil
}
