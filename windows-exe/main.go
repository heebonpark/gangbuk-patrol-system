// PatrolOps 로컬 실행기
//
// index.html을 빌드 시점에 실행 파일 안에 통째로 넣고(go:embed), 실행할 때마다
// 이 PC의 임시 폴더에 같은 경로로 풀어놓은 뒤 기본 브라우저로 엽니다.
// 매번 같은 경로에 쓰기 때문에 브라우저 저장소(localStorage)에 저장된
// 근무표 데이터는 실행할 때마다 그대로 유지됩니다.
//
// 주의: 이 실행 파일은 빌드 시점의 index.html 스냅샷을 담고 있습니다.
// 앱이 업데이트되면 이 실행 파일도 새로 빌드해서 다시 배포해야 최신 화면이 보입니다.
package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

//go:embed index.html
var indexHTML []byte

func messageBox(title, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	proc.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}

func main() {
	// 임시 폴더 이름에 내용물(index.html) 해시를 붙인다. 예전엔 모든 빌드가 똑같이
	// "%TEMP%\PatrolOps"를 썼는데, 그러면 데모 데이터판(PatrolOps.exe)과 빈 데이터판
	// (PatrolOps_Blank.exe)을 같은 PC에서 번갈아 실행할 때마다 서로의 임시 파일을
	// 덮어써서 — 한쪽을 실행한 직후 다른 쪽을 실행하면 그 시점에 브라우저나 백신이
	// 그 파일을 잠깐 잡고 있는 경우 덮어쓰기가 실패해 "실행이 안 되는" 것처럼 보일 수
	// 있었다. 내용이 다르면 폴더도 자동으로 달라지므로 서로 절대 충돌하지 않고, 같은
	// 빌드는 항상 같은 폴더를 쓰므로(해시가 고정) localStorage가 유지된다는 원래 의도도
	// 그대로 유지된다.
	hash := sha256.Sum256(indexHTML)
	tmpDir := filepath.Join(os.TempDir(), "PatrolOps_"+hex.EncodeToString(hash[:])[:8])
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		messageBox("PatrolOps 실행 오류", "임시 폴더를 만들지 못했습니다.\n"+err.Error())
		return
	}

	htmlPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(htmlPath, indexHTML, 0644); err != nil {
		messageBox("PatrolOps 실행 오류", "실행에 필요한 파일을 준비하지 못했습니다.\n"+err.Error())
		return
	}

	cmd := exec.Command("cmd", "/c", "start", "", htmlPath)
	if err := cmd.Start(); err != nil {
		messageBox("PatrolOps 실행 오류", "브라우저를 여는 데 실패했습니다.\n"+err.Error())
	}
}
