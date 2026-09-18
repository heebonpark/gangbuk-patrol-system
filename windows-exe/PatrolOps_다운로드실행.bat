@echo off
chcp 65001 >nul
setlocal
title PatrolOps 다운로드 및 실행

rem raw.githubusercontent.com은 백신/브라우저 보안 프로그램이 더 의심스럽게 취급하는 경우가
rem 많아(다운로드가 검사 단계에서 멈추거나 막힘), 정식 배포 채널인 GitHub Releases를 쓴다.
rem 그리고 exe를 직접 받으면 회사 방화벽/메일 필터가 ".exe" 확장자 자체를 기계적으로
rem 막는 경우가 많아서, zip으로 감싼 파일을 받아 그 자리에서 풀어(Expand-Archive, Windows
rem 10 이상 기본 내장) 실행한다.
set "URL=https://github.com/heebonpark/gangbuk-patrol-system/releases/download/latest/PatrolOps_Blank.zip"
set "ZIP=%~dp0PatrolOps_Blank.zip"
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
    curl -L -f -o "%ZIP%" "%URL%"
) else (
    powershell -NoProfile -Command "try { Invoke-WebRequest -Uri '%URL%' -OutFile '%ZIP%' } catch { exit 1 }"
)

if not exist "%ZIP%" (
    echo.
    echo [오류] 다운로드에 실패했습니다.
    echo  - 인터넷 연결을 확인해주세요.
    echo  - 계속 안 되면 관리자에게 문의하세요.
    echo.
    pause
    exit /b 1
)

echo.
echo 압축을 푸는 중입니다...
powershell -NoProfile -Command "Expand-Archive -Path '%ZIP%' -DestinationPath '%~dp0' -Force"
del "%ZIP%" >nul 2>nul

if not exist "%OUT%" (
    echo.
    echo [오류] 압축 해제에 실패했습니다.
    echo  - 이 PC에 PowerShell의 Expand-Archive가 없을 수 있습니다(Windows 10 미만).
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
