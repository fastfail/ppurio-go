package ppuriogo

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const URI = "https://message.ppurio.com"
const URI_REQUEST_ACCESS_TOKEN = URI + "/v1/token"
const URI_MESSAGE = URI + "/v1/message"

type Ppurio struct {
	ppurioAccount string
	from          string

	encodedBasicAPIKey string

	accessToken       string
	accessTokenExpire time.Time
	client            *http.Client
}

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

var ppurioErrors map[int]error

type PpurioError struct {
	message string
}

func (p PpurioError) Error() string {
	return p.message
}

var ErrorInvalidRequest = PpurioError{message: "잘못된 요청"}
var ErrorInvalidUrl = errors.New("잘못된 URL")
var ErrorInvalidAuthorizationError = errors.New("토큰 발급 호출 시, Authorization 헤더가 유효하지 않음")
var ErrorInvalidToken3002 = errors.New("토큰이 유효하지 않음")
var ErrorInvalidIP = errors.New("아이피가 유효하지 않음")
var ErrorInvalidAccount = errors.New("계정이 유효하지 않음")
var ErrorInvalidToken3005 = errors.New("토큰이 유효하지 않음")
var ErrorInvalidAuthenticationHeader = errors.New("Authentication Header가 유효하지 않음")
var ErrorFailedToIssueAccessToken = errors.New("엑세스 토큰 발행 실패")
var ErrorTooManyRequests = errors.New("너무 많은 요청")
var ErrorAPIAccessDisabled = errors.New("api 접근 권한이 비활성화 상태")
var ErrorInvalidAuthKey = errors.New("인증키가 유효하지 않음")
var ErrorNeverIssueAuthKey = errors.New("인증키를 발행 받지 않음")
var ErrorInvalidMessageKey = errors.New("메시지키가 유효하지 않음")
var ErrorTooLateToCancelReservation = errors.New("예약 취소 가능 시간이 지남")
var ErrorAlreadySendingMessage = errors.New("메시지가 이미 발송중인 상태")
var ErrorCantCancelReservation = errors.New("예약을 취소할 수 없음")

func init() {
	ppurioErrors = map[int]error{
		2000: ErrorInvalidRequest,
		2001: ErrorInvalidUrl,
		3001: ErrorInvalidAuthorizationError,
		3002: ErrorInvalidToken3002,
		3003: ErrorInvalidIP,
		3004: ErrorInvalidAccount,
		3005: ErrorInvalidToken3005,
		3006: ErrorInvalidAuthenticationHeader,
		3007: ErrorFailedToIssueAccessToken,
		3008: ErrorTooManyRequests,
		4004: ErrorAPIAccessDisabled,
		4006: ErrorInvalidAuthKey,
		4007: ErrorNeverIssueAuthKey,
		4009: ErrorInvalidMessageKey,
		4010: ErrorTooLateToCancelReservation,
		4011: ErrorAlreadySendingMessage,
		4012: ErrorCantCancelReservation,
	}
}

func NewPpurio(ppurioAccount, accessKey, from string) (*Ppurio, error) {
	p := &Ppurio{
		ppurioAccount:      ppurioAccount,
		encodedBasicAPIKey: base64.RawURLEncoding.EncodeToString([]byte(strings.Join([]string{ppurioAccount, accessKey}, ":"))),
		from:               from,
		client:             &http.Client{},
	}
	err := p.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	return p, nil
}
func NewPpurioWithTimeout(ppurioAccount, accessKey, from string, timeout time.Duration) (*Ppurio, error) {
	p, err := NewPpurio(ppurioAccount, accessKey, from)
	if err != nil {
		return nil, err
	}
	p.client = &http.Client{Timeout: timeout}
	return p, nil
}

func (p *Ppurio) TextMessage(messageParam *MessageParams) (string, error) {
	param, err := json.Marshal(&messageParam)
	if err != nil {
		return "", fmt.Errorf("failed to marshal param: %w", err)
	}
	fmt.Println("param:", string(param))
	resp, err := p.doRequest("POST", URI_MESSAGE, param)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var e ErrorResponse
		err := json.Unmarshal(body, &e)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal error response: %s", string(body))
		}
		errCode, err := strconv.Atoi(e.Code)
		if err != nil {
			return "", fmt.Errorf("unexpected error code response: %s", e.Code)
		}
		err, ok := ppurioErrors[errCode]
		if !ok {
			return "", fmt.Errorf("unexpected error code response: %d", errCode)
		}
		return "", err
	}

	var r MessageResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Code != "1000" || r.Description != "ok" {
		return "", fmt.Errorf("unexpected result - code: %d, description: %s", r.Code, r.Description)
	}
	return r.MessageKey, nil
}
