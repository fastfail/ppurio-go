package ppuriogo

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MessageParams struct {
	Account       string        `json:"account"`
	MessageType   MessageType   `json:"messageType"`
	Content       string        `json:"content"`
	From          string        `json:"from"`
	DuplicateFlag DuplicateFlag `json:"duplicateFlag"`
	TargetCount   int           `json:"targetCount"`
	Targets       []Target      `json:"targets"`
	RefKey        string        `json:"refKey"`
	RejectType    *RejectType   `json:"rejectType,omitempty"`
	SendTime      *time.Time    `json:"sendTime,omitempty"`
	Subject       *string       `json:"subject,omitempty"`
	Files         *[]File       `json:"files,omitempty"`
}
type MessageType string

const (
	MessageTypeSMS MessageType = "SMS"
	MessageTypeLMS MessageType = "LMS"
	MessageTypeMMS MessageType = "MMS"
)

type DuplicateFlag string

const (
	DuplicateFlagY DuplicateFlag = "Y"
	DuplicateFlagN DuplicateFlag = "N"
)

type RejectType string

const RejectTypeAD RejectType = "AD"

type ChangeWord struct {
	Var1 string `json:"var1,omitempty"`
	Var2 string `json:"var2,omitempty"`
	Var3 string `json:"var3,omitempty"`
	Var4 string `json:"var4,omitempty"`
	Var5 string `json:"var5,omitempty"`
	Var6 string `json:"var6,omitempty"`
	Var7 string `json:"var7,omitempty"`
}

type Target struct {
	To         string  `json:"to"`
	ChangeWord *string `json:"changeWord,omitempty"`
	Name       *string `json:"name,omitempty"`
}

type File struct {
	Name string `json:"name"` // 파일 명
	Size int    `json:"size"` // 파일 크기, byte 단위
	Data string `json:"data"` // 바이너리 형식인 이미지 파일을 Base64 인코딩한 텍스트"
}

func (p *Ppurio) NewTextMessageParams(content2000 string, duplicateFlag DuplicateFlag, targets []Target, useRejectAD bool, subject30 *string, files *[]File) *MessageParams {
	m := &MessageParams{
		Account:       p.ppurioAccount,
		Content:       content2000,
		From:          p.from,
		DuplicateFlag: duplicateFlag,
		TargetCount:   len(targets),
		Targets:       targets,
		RefKey:        strings.Replace(uuid.NewString(), "-", "", 4),
	}
	if messageContentBytes(content2000) > 90 {
		m.MessageType = MessageTypeLMS
	} else {
		m.MessageType = MessageTypeSMS
	}
	if useRejectAD {
		rejectType := RejectTypeAD
		m.RejectType = &rejectType
	}
	if subject30 != nil && *subject30 != "" {
		m.Subject = subject30
	}
	if files != nil && len(*files) > 0 {
		m.Files = files
	}
	return m
}

// 예약용
func (p *Ppurio) NewMessageParamsWithSendTime(content2000 string, duplicateFlag DuplicateFlag, targets []Target, rejectAD bool, subject30 *string, files *[]File, sendTime *time.Time) *MessageParams {
	m := p.NewTextMessageParams(content2000, duplicateFlag, targets, rejectAD, subject30, files)
	m.SendTime = sendTime
	return m
}

type MessageResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	RefKey      string `json:"refKey"`
	MessageKey  string `json:"messageKey"`
}

func messageContentBytes(s string) int {
	total := 0
	for _, r := range s {
		if r < 128 { // ASCII
			total += 1
		} else {
			total += 2
		}
	}
	return total
}
