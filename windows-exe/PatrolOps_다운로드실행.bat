@echo off
chcp 65001 >nul
setlocal
title PatrolOps 다운로드 및 실행

set "URL=https://raw.githubusercontent.com/heebonpark/gangbuk-patrol-system/main/windows-exe/PatrolOps_Blank.exe"
set "OUT=%~dp0PatrolOps_Blank.exe"

echo ============================================
echo   PatrolOps (빈 데이터 배포용) 다운로드
echo ============================================
echo.
echo GitHub에서 최신 파일을 내려받는 중입니다...
echo (이 창을 실행할 때마다 항상 최신 버전을 새로 받습니다)
echo.

where curl >nul 2>nul
if %errorlevel%==0 (
    curl -L -f -o "%OUT%" "%URL%"
) else (
    powershell -NoProfile -Command "try { Invoke-WebRequest -Uri '%URL%' -OutFile '%OUT%' } catch { exit 1 }"
)

if not exist "%OUT%" (
    echo.
    echo [오류] 다운로드에 실패했습니다.
    echo  - 인터넷 연결을 확인해주세요.
    echo  - 계속 안 되면 관리자에게 문의하세요.
    echo.
    pause
    exit /b 1
)

echo.
echo 다운로드 완료. 프로그램을 실행합니다...
echo (처음 실행하면 Windows가 "알 수 없는 게시자" 경고를 띄울 수 있습니다.
echo  "추가 정보" -^> "실행"을 눌러주세요.)
echo.
start "" "%OUT%"

exit /b 0
