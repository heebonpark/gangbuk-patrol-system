#!/bin/bash
# 근무표 관리 통합관리 시스템을 macOS 기본 브라우저로 엽니다.
# (Finder에서 이 파일을 더블클릭하면 실행됩니다. 처음 실행 시 "확인되지 않은 개발자"
#  경고가 뜨면 파일을 마우스 우클릭(또는 control+클릭) 후 "열기"를 선택하세요.)
cd "$(dirname "$0")"

if [ ! -f "index.html" ]; then
  echo "[오류] 이 파일과 같은 폴더에 index.html이 없습니다."
  read -p "엔터 키를 누르면 종료합니다..."
  exit 1
fi

echo "근무표 관리 통합관리 시스템을 기본 브라우저로 엽니다..."
open "index.html"
