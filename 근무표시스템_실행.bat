@echo off
chcp 65001 >nul
cd /d "%~dp0"

if not exist "index.html" (
  echo [오류] 이 배치파일과 같은 폴더에 index.html이 없습니다.
  pause
  exit /b 1
)

echo 근무표 관리 통합관리 시스템을 기본 브라우저로 엽니다...
start "" "index.html"
