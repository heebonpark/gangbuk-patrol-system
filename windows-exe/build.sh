#!/bin/bash
# index.html이 바뀐 뒤 PatrolOps.exe를 다시 만들 때 사용합니다.
# 필요 조건: Go (brew install go), macOS/Linux에서도 Windows용으로 크로스 빌드됩니다.
set -e
cd "$(dirname "$0")"
cp ../index.html ./index.html
GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui -s -w" -o "PatrolOps.exe" .
echo "빌드 완료: windows-exe/PatrolOps.exe"
