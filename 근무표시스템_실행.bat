@echo off
chcp 65001 >nul
cd /d "%~dp0"

if not exist "index.html" (
  echo [오류] 이 배치파일과 같은 폴더에 index.html이 없습니다.
  pause
  exit /b 1
)

rem 참고: 이 방식(file:// 직접 열기)에서는 카카오맵 JS SDK가 도메인 인증을 할 수 없어
rem 고객지도/지도위치확인 화면이 동작하지 않습니다. 그 기능이 필요하면 windows-exe 폴더의
rem PatrolOps.exe(로컬 HTTP 서버로 여는 방식)를 사용해주세요.
echo 근무표 관리 통합관리 시스템을 기본 브라우저로 엽니다...
start "" "index.html"
