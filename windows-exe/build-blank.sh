#!/bin/bash
# 지사/직원/순회일정/경비지도사 데이터가 전부 빈 상태인 배포용 exe(PatrolOps_Blank.exe)를
# 만듭니다. 실제 고객(다른 회사)에게 처음 나눠줄 때는 이 파일을 쓰세요 — 로그인 화면에서
# "관리자로 로그인"(기본 비밀번호 admin1234+!)한 뒤, 지사/직원을 새로 등록하면 됩니다.
# (지사 로그인 화면의 "+ 새 지사 등록"으로 지사 자체도 바로 만들 수 있습니다.)
#
# 주의: windows-exe/index.html은 main.go가 go:embed로 그대로 담는 파일이라, 이 스크립트는
# 잠시 그 파일을 빈 데이터 버전으로 바꿔치기해서 빌드한 뒤, 평소 개발용 exe(PatrolOps.exe)
# 빌드에 영향이 없도록 다시 원래(데모 데이터) 버전으로 복원합니다.
set -e
cd "$(dirname "$0")"

node make-blank-html.js ../index.html ./index.html
cp ../vehicle-management-demo.html ./vehicle-management.html
# 반드시 vehicle-management-demo.html(마스킹판)만 써야 한다 — ../vehicle-management.html은
# 실제 회사 개인정보(실명·연락처·계좌번호)가 그대로 든 파일이라 배포용 exe에 절대 담으면
# 안 된다(2026-09-15 사고 참고). 데모판도 완전히 빈 상태는 아니므로, 새 고객 배포판에서는
# 앱 안의 "관리자 > 데이터 초기화"로 직접 비울 수 있다.
if [ -f ../vehicle-management.html ] && cmp -s ../vehicle-management.html ./vehicle-management.html; then
  echo "오류: 실제 개인정보가 담긴 vehicle-management.html이 그대로 embed될 뻔했습니다. 빌드를 중단합니다." >&2
  exit 1
fi

if ! command -v goversioninfo >/dev/null 2>&1; then
  GOVERSIONINFO_BIN="$(go env GOPATH)/bin/goversioninfo"
  if [ ! -x "$GOVERSIONINFO_BIN" ]; then
    echo "goversioninfo 설치 중..."
    go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
  fi
else
  GOVERSIONINFO_BIN="goversioninfo"
fi
"$GOVERSIONINFO_BIN" -o resource.syso versioninfo.json

GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui -s -w" -o "PatrolOps_Blank.exe" .

# 평소 개발용 exe(PatrolOps.exe)를 다시 빌드할 때 쓰는 index.html을 원래(데모 데이터)
# 버전으로 복원 — 안 그러면 다음에 build.sh를 돌릴 때 빈 데이터가 섞여 들어감.
cp ../index.html ./index.html

echo "빌드 완료: windows-exe/PatrolOps_Blank.exe (지사/직원/일정 데이터 없음, 관리자 비밀번호는 admin1234+!)"
