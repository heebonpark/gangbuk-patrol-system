// PatrolOps 로컬 실행기
//
// index.html 등을 빌드 시점에 실행 파일 안에 통째로 넣고(go:embed), 실행할 때마다
// 이 PC의 임시 폴더에 같은 경로로 풀어놓은 뒤 기본 브라우저로 엽니다. 매번 같은
// 경로에 쓰기 때문에 브라우저 저장소(localStorage)에 저장된 근무표/차량관리
// 데이터는 실행할 때마다 그대로 유지됩니다.
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
// kakao-map.html(고객지도)/kakao-map-admin.html(지도 관리자 대시보드)/map-check.html
// (지도위치확인)도 같은 이유로 시드 데이터가 없다 — 주소도 CSV/엑셀을 그 자리에서
// 업로드해 지오코딩하고, 카카오 지도 JS 키도 화면에서 직접 입력한다(또는 기본 키를 그대로
// 쓴다). 카카오 지도 SDK를 CDN(dapi.kakao.com)에서 불러오므로, 이 화면들도 실행 PC에
// 인터넷 연결이 있어야 한다.
//
// 2026-09-18부터: 카카오맵 JS SDK는 "요청을 보낸 도메인"이 카카오 개발자센터에 등록된
// 도메인과 정확히 일치해야 동작하는데, 예전 방식(브라우저로 file:///C:/... 경로를 직접
// 여는 방식)은 주소창이 file:// 이라 도메인 매칭 자체가 성립하지 않아 카카오맵 관련
// 화면(고객지도/지도위치확인)이 구조적으로 절대 동작할 수 없었다. 그래서 이제 exe가
// 파일을 직접 열지 않고, 고정 포트(47291)로 이 PC 안에서만 도는 로컬 HTTP 서버를 띄운
// 뒤 http://localhost:47291/index.html 을 연다. 이 정확한 주소를 카카오 개발자센터의
// JavaScript SDK 도메인 목록에 등록해두면 카카오맵이 정상 동작한다.
//
// 로컬 서버는 이 exe 프로세스가 떠 있는 동안만 응답하므로(select{}로 계속 대기),
// 브라우저 창을 다 닫아도 프로세스 자체는 백그라운드에 남는다 — 작업 관리자에서
// "PatrolOps.exe"로 보이며, 필요하면 거기서 종료할 수 있다. 앱을 새로 빌드해 다시
// 실행하면, 시작 시 먼저 이전에 떠 있는 서버에 종료 신호(/__shutdown)를 보내고 최신
// 파일로 새로 띄우므로 항상 최신 버전이 보인다(오래된 서버가 포트를 붙들고 있어
// 업데이트가 안 보이는 문제를 자동으로 해결).
//
// 주의: 이 실행 파일은 빌드 시점의 스냅샷을 담고 있습니다.
// 앱이 업데이트되면 이 실행 파일도 새로 빌드해서 다시 배포해야 최신 화면이 보입니다.
package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

//go:embed index.html
var indexHTML []byte

//go:embed vehicle-management.html
var vehicleManagementHTML []byte

//go:embed voc-management.html
var vocManagementHTML []byte

//go:embed kakao-map.html
var kakaoMapHTML []byte

//go:embed kakao-map-admin.html
var kakaoMapAdminHTML []byte

//go:embed subscription-dashboard.html
var subscriptionDashboardHTML []byte

//go:embed map-check.html
var mapCheckHTML []byte

//go:embed dong-boundaries.json
var dongBoundariesJSON []byte

// 카카오 개발자센터 JavaScript SDK 도메인 목록에 http://localhost:47291 을
// 등록해둬야 카카오맵(고객지도/지도위치확인)이 동작한다. 이 포트를 바꾸면
// 그쪽 등록값도 같이 바꿔야 한다.
const serverPort = "47291"

func messageBox(title, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	proc.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}

// 이전 실행에서 띄워둔 로컬 서버가 아직 떠 있으면(오래된 버전일 수 있음) 내려달라고
// 요청한다. 응답이 오지 않아도(애초에 떠 있지 않았던 정상적인 경우) 무시한다.
func shutdownPreviousServer() {
	client := http.Client{Timeout: 500 * time.Millisecond}
	client.Post("http://127.0.0.1:"+serverPort+"/__shutdown", "text/plain", nil)
	time.Sleep(200 * time.Millisecond)
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

	files := map[string][]byte{
		"index.html":                  indexHTML,
		"vehicle-management.html":     vehicleManagementHTML,
		"voc-management.html":         vocManagementHTML,
		"kakao-map.html":              kakaoMapHTML,
		"kakao-map-admin.html":        kakaoMapAdminHTML,
		"subscription-dashboard.html": subscriptionDashboardHTML,
		"map-check.html":              mapCheckHTML,
		"dong-boundaries.json":        dongBoundariesJSON,
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), data, 0644); err != nil {
			messageBox("PatrolOps 실행 오류", "실행에 필요한 파일을 준비하지 못했습니다.\n"+err.Error())
			return
		}
	}

	shutdownPreviousServer()

	ln, err := net.Listen("tcp", "127.0.0.1:"+serverPort)
	if err != nil {
		messageBox("PatrolOps 실행 오류", "로컬 서버를 시작하지 못했습니다(포트 "+serverPort+" 사용 중).\n"+err.Error())
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/__shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		go func() {
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}()
	})
	mux.Handle("/", http.FileServer(http.Dir(tmpDir)))
	go http.Serve(ln, mux)

	url := "http://localhost:" + serverPort + "/index.html"
	cmd := exec.Command("cmd", "/c", "start", "", url)
	if err := cmd.Start(); err != nil {
		messageBox("PatrolOps 실행 오류", "브라우저를 여는 데 실패했습니다.\n"+err.Error())
		return
	}

	// 브라우저가 계속 이 로컬 서버로 요청을 보내야 하므로, 이 프로세스도 계속 떠 있는다.
	select {}
}
