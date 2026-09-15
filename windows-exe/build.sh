#!/bin/bash
# index.html이 바뀐 뒤 PatrolOps.exe를 다시 만들 때 사용합니다.
# 필요 조건: Go (brew install go), macOS/Linux에서도 Windows용으로 크로스 빌드됩니다.
#
# 개인정보 안내(2026-09-15 이후): vehicle-management.html은 이제 시드 데이터 자체가
# 항상 마스킹되어 있고(실명 → 이O호 등), 회사 실제 데이터는 앱 안의 "엑셀 업로드"로
# 그 자리에서 직접 올려 각자 브라우저의 localStorage에만 저장한다 — 이 파일(따라서 이
# exe)에는 실제 개인정보가 절대 포함되지 않는다. (예전엔 별도의 마스킹 안 된 실제
# 파일을 실수로 그대로 embed해 82개 커밋에 걸쳐 유출된 사고가 있었는데, 그 원인이었던
# "실제 파일 vs 데모 파일" 이원화 구조 자체를 없앤 것이 이번 수정이다.)
set -e
cd "$(dirname "$0")"
cp ../index.html ./index.html
cp ../vehicle-management.html ./vehicle-management.html
cp ../voc-management.html ./voc-management.html

# versioninfo.json에 적어둔 저작권/제품 정보를 exe 파일 속성(우클릭 > 속성 > 자세히)에
# 심는다. goversioninfo가 없으면 최초 1회 설치(go install)한 뒤 resource.syso를 새로
# 생성 — go build가 같은 디렉터리의 .syso 파일을 자동으로 링크해준다.
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

GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui -s -w" -o "PatrolOps.exe" .
echo "빌드 완료: windows-exe/PatrolOps.exe"
