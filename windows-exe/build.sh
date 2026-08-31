#!/bin/bash
# index.html이 바뀐 뒤 PatrolOps.exe를 다시 만들 때 사용합니다.
# 필요 조건: Go (brew install go), macOS/Linux에서도 Windows용으로 크로스 빌드됩니다.
set -e
cd "$(dirname "$0")"
cp ../index.html ./index.html

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
