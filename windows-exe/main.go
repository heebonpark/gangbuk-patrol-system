// PatrolOps 로컬 실행기
//
// index.html과 vehicle-management.html을 빌드 시점에 실행 파일 안에 통째로 넣고
// (go:embed), 실행할 때마다 이 PC의 임시 폴더에 같은 경로로 풀어놓은 뒤 기본
// 브라우저로 엽니다. 매번 같은 경로에 쓰기 때문에 브라우저 저장소(localStorage)에
// 저장된 근무표/차량관리 데이터는 실행할 때마다 그대로 유지됩니다.
//
// 개인정보 안내(2026-09-15부터): vehicle-management.html의 시드 데이터는 index.html
// 시드처럼 항상 마스킹되어 있다(실명 → 이O호 등). 이전에는 이 파일이 실제 회사
// 개인정보를 그대로 담은 별도 파일이었고, 그 실제 파일을 실수로 이 exe에 그대로
// embed해 82개 커밋에 걸쳐 Public 저장소에 유출된 사고가 있었다 — "실제 파일 vs
// 마스킹 데모 파일" 이원화 구조 자체가 원인이었기에, 지금은 파일을 하나로 합치고
// 시드 자체를 마스킹해서 그 이원화를 아예 없앴다. 회사 실제 데이터는 앱 안의
// "엑셀 업로드"로 그 자리에서 직접 올려 쓰며, 그 값은 이 소스 파일이 아니라 실행
// 중인 브라우저의 localStorage에만 저장된다.
//
// voc-management.html(VOC관리)도 같은 방식으로 embed한다 — 시드 데이터 없이 CSV를 그
// 자리에서 업로드해 쓰는 구조라 처음부터 개인정보가 들어있지 않다. 다만 xlsx.js·Chart.js를
// CDN(cdnjs)에서 불러오므로, 이 화면만은 실행 PC에 인터넷 연결이 있어야 정상 동작한다.
//
// 주의: 이 실행 파일은 빌드 시점의 스냅샷을 담고 있습니다.
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

//go:embed vehicle-management.html
var vehicleManagementHTML []byte

//go:embed voc-management.html
var vocManagementHTML []byte

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
	if err := os.WriteFile(filepath.Join(tmpDir, "vehicle-management.html"), vehicleManagementHTML, 0644); err != nil {
		messageBox("PatrolOps 실행 오류", "실행에 필요한 파일을 준비하지 못했습니다.\n"+err.Error())
		return
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "voc-management.html"), vocManagementHTML, 0644); err != nil {
		messageBox("PatrolOps 실행 오류", "실행에 필요한 파일을 준비하지 못했습니다.\n"+err.Error())
		return
	}

	cmd := exec.Command("cmd", "/c", "start", "", htmlPath)
	if err := cmd.Start(); err != nil {
		messageBox("PatrolOps 실행 오류", "브라우저를 여는 데 실패했습니다.\n"+err.Error())
	}
}
