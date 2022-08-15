package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const RestBaseURL = "https://api.kraken.com"

type RestClient struct {
	apiKey     string
	privateKey string
	decodedKey []byte
	baseUrl    string
	httpClient *http.Client
}

func NewRestClient(apiKey, privateKey string) *RestClient {
	decoded, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Printf("cannot decode privateKey: %v", err)
	}
	return &RestClient{
		apiKey:     apiKey,
		privateKey: privateKey,
		decodedKey: decoded,
		baseUrl:    RestBaseURL,
		httpClient: &http.Client{Timeout: time.Second * 30},
	}
}

type wsTokenResponse struct {
	Result WsAuthToken `json:"result"`
	Error  []string    `json:"error"`
}

type WsAuthToken struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

func (r *RestClient) WsToken() (*WsAuthToken, error) {
	payload := make(url.Values)
	nonce := time.Now().UnixMilli()
	payload.Set("nonce", fmt.Sprintf("%d", nonce))
	resp, err := r.post("/0/private/GetWebSocketsToken", payload)
	if err != nil {
		return nil, fmt.Errorf("cannot request auth token for WS: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	reader := json.NewDecoder(resp.Body)
	var data wsTokenResponse
	if err = reader.Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot decode response: %v", err)
	}
	if len(data.Error) > 0 {
		return nil, fmt.Errorf("remote server returned an error: %v", data.Error)
	}
	return &data.Result, nil
}

type balanceResponse struct {
	Result map[string]string `json:"result"`
	Error  []string          `json:"error"`
}

type Balances map[string]float64

func (r *RestClient) Balances() (Balances, error) {
	payload := make(url.Values)
	nonce := time.Now().UnixMilli()
	payload.Set("nonce", fmt.Sprintf("%d", nonce))
	resp, err := r.post("/0/private/Balance", payload)
	if err != nil {
		return nil, fmt.Errorf("cannot get balances: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	reader := json.NewDecoder(resp.Body)
	var respData balanceResponse
	if err = reader.Decode(&respData); err != nil {
		return nil, fmt.Errorf("cannot decode response: %w", err)
	}
	if len(respData.Error) > 0 {
		return nil, fmt.Errorf("remote server returned an error: %v", respData.Error)
	}
	log.Printf("balance response: %#v", respData.Result)
	balances := make(Balances)
	for c, v := range respData.Result {
		balances[c], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse balance volume: %w; response: %v", err, respData)
		}
	}
	return balances, nil
}

type OrderResp struct {
	Description struct {
		Order string `json:"order"`
		Close string `json:"close"`
	} `json:"descr"`
	TxId []string `json:"txid"`
}

func (r *RestClient) post(uri string, data url.Values) (*http.Response, error) {
	fullUrl := r.baseUrl + uri
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("API-Key", r.apiKey)
	req.Header.Set("API-Sign", r.sign(uri, data))
	return r.httpClient.Do(req)
}

func (r *RestClient) sign(uriPath string, data url.Values) string {
	sha := sha256.New()
	sha.Write([]byte(data.Get("nonce") + data.Encode()))
	shaSum := sha.Sum(nil)
	mac := hmac.New(sha512.New, r.decodedKey)
	mac.Write(append([]byte(uriPath), shaSum...))
	macSum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(macSum)
}
