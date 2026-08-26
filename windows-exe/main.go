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
	_ "embed"
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
	tmpDir := filepath.Join(os.TempDir(), "PatrolOps")
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
