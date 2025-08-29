package ppuriogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (p *Ppurio) doRequest(method, url string, body []byte) (*http.Response, error) {
	if p.accessToken == "" || p.accessTokenExpire.After(time.Now().Add(-60*time.Second)) {
		err := p.getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get access token: %w", err)
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		err := p.getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get access token: %w", err)
		}
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+p.accessToken)
		req.Header.Set("Content-Type", "application/json")
		return p.client.Do(req)
	}
	return resp, nil
}

func (p *Ppurio) getAccessToken() error {
	type IssueAccessTokenResponse struct {
		Token   string `json:"token"`
		Type    string `json:"type"`
		Expired string `json:"expired"`
	}
	req, _ := http.NewRequest("POST", URI_REQUEST_ACCESS_TOKEN, nil)
	req.Header.Add("Authorization", "Basic "+p.encodedBasicAPIKey)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var e ErrorResponse
		err := json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal error response: %s", string(body))
		}
		errCode, err := strconv.Atoi(e.Code)
		if err != nil {
			return fmt.Errorf("unexpected error code response: %s", e.Code)
		}
		err, ok := ppurioErrors[errCode]
		if !ok {
			return fmt.Errorf("unexpected error code response: %d", errCode)
		}
		return err
	}

	var r IssueAccessTokenResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	if r.Token == "" {
		return errors.New("failed to issue token")
	}
	exp, err := time.ParseInLocation("20060102150405", r.Expired, time.FixedZone("Asia/seoul", 9))
	if err != nil {
		return err
	}
	p.accessToken = r.Token
	p.accessTokenExpire = exp
	return nil
}
