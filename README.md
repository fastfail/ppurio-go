# ppurio-go
뿌리오(ppurio.com) 문자서비스 API 비공식 Go 패키지입니다.

### Installation
```sh
go get github.com/fastfail/ppurio-go
```

<!-- USAGE EXAMPLES -->
## Usage

```go
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
```

<!-- ROADMAP -->
## Features

- [X] AccessToken 자동 갱신
- [ ] 문자 전송
    - [X] SMS(~90Byte)
    - [X] LMS(91~2000Byte)
    - [ ] MMS
- [ ] 예약 전송
- [ ] 카카오톡 전송

## Disclaimer

이 프로젝트는 어떠한 보증도 없이 "있는 그대로" 제공됩니다.  
본 소프트웨어 사용으로 인한 손해, 데이터 손실, 문제 등에 대해 개발자는 책임을 지지 않습니다.  

본 프로젝트는 학습 및 연구 목적으로만 제공됩니다.  
상업적 사용 또는 운영 환경에서의 사용 시 발생하는 모든 책임은 사용자에게 있습니다.  
또한, 외부 API나 제3자 서비스 사용 시 해당 서비스의 약관을 반드시 확인하시기 바랍니다.