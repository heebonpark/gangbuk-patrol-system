#!/bin/bash
# index.html이 바뀐 뒤 PatrolOps.exe를 다시 만들 때 사용합니다.
# 필요 조건: Go (brew install go), macOS/Linux에서도 Windows용으로 크로스 빌드됩니다.
#
# 중요(2026-09-15 사고 이후): 반드시 vehicle-management-demo.html(마스킹판)만 embed한다.
# ../vehicle-management.html은 실제 회사 개인정보(실명·연락처·계좌번호)가 그대로 든
# 파일이라 이 저장소는 물론 exe(공개 배포물)에도 절대 들어가면 안 된다 — 예전에 이
# 스크립트가 그 파일을 그대로 복사해 exe에 실제로 embed했다가 Public 저장소에 82개
# 커밋에 걸쳐 유출된 사고가 있었다. 개인적으로 실데이터가 든 exe가 필요하면 이 스크립트를
# 쓰지 말고 그 자리에서 직접 vehicle-management.html을 수동으로 넣어 빌드한 뒤, 그 결과물은
# 커밋/공유하지 말 것.
set -e
cd "$(dirname "$0")"
cp ../index.html ./index.html
cp ../vehicle-management-demo.html ./vehicle-management.html
# 안전장치: 지금 embed하려는 파일이 실제 개인정보 파일과 우연히도(혹은 실수로) 똑같으면
# 빌드를 바로 중단한다 — 사람이 이 스크립트를 잘못 고쳐도 실제 데이터가 exe에 섞여
# 들어가는 걸 한 번 더 막기 위함.
if [ -f ../vehicle-management.html ] && cmp -s ../vehicle-management.html ./vehicle-management.html; then
  echo "오류: 실제 개인정보가 담긴 vehicle-management.html이 그대로 embed될 뻔했습니다. 빌드를 중단합니다." >&2
  exit 1
fi

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
