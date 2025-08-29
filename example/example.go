package main

import (
	"fmt"

	ppuriogo "github.com/fastfail/ppurio-go"
)

func main() {
	p, err := ppuriogo.NewPpurio("뿌리오 ID", "API 키", "발신인번호")
	if err != nil {
		panic(fmt.Errorf("failed to initialize ppurio-go: %w", err))
	}
	receiver1 := "수신인번호1"

	// SMS 전송 (~90Byte)
	messageKey, err := p.TextMessage(p.NewTextMessageParams("일이삼사오륙칠팔구",
		ppuriogo.DuplicateFlagN,
		[]ppuriogo.Target{{To: receiver1}}, false, nil, nil))
	if err != nil {
		fmt.Println("전송오류: ", err.Error())
	} else {
		fmt.Println("messageKey:", messageKey)
	}

	// LMS 전송 (91Byte~, 한글 2Byte 취급)
	messageKey, err = p.TextMessage(p.NewTextMessageParams("일이삼사오륙칠팔구일이삼사오륙칠팔구일이삼사오륙칠팔구일이삼사오륙칠팔구일이삼사오륙칠팔구",
		ppuriogo.DuplicateFlagN,
		[]ppuriogo.Target{{To: receiver1}}, false, nil, nil))
	if err != nil {
		fmt.Println("전송오류: ", err.Error())
	} else {
		fmt.Println("messageKey:", messageKey)
	}
}
