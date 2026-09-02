// 배포용(빈 데이터) index.html을 만드는 스크립트.
// 개발/테스트용 원본(../index.html)의 시드 데이터(경비지도사·지사·순회일정·직원)만
// 전부 빈 배열로 바꾼 사본을 만든다 — 로직/디자인/기능은 그대로, 데이터만 비움.
// 정확한 함수 경계를 regex로 억지로 맞추면 실수하기 쉬워서, 각 시드 함수 바로 다음에
// 나오는 "다음 선언부"를 기준점(anchor)으로 삼아 그 사이 구간을 통째로 교체한다.
// 원본 구조(seedGuides → seedBranches → seedRecords → RAW_EMPLOYEES_2026_08 →
// seedEmployees)가 바뀌면 이 스크립트도 같이 손봐야 한다(아래에서 anchor를 못 찾으면
// 바로 에러를 내서 조용히 잘못된 파일이 만들어지는 걸 막는다).

const fs = require('fs');
const path = require('path');

const srcPath = process.argv[2] || path.join(__dirname, '..', 'index.html');
const outPath = process.argv[3] || path.join(__dirname, 'index.html');

let html = fs.readFileSync(srcPath, 'utf8');

function sliceBetween(html, startMarker, endMarker, label){
  const startIdx = html.indexOf(startMarker);
  if(startIdx < 0) throw new Error(`[make-blank-html] "${label}" 시작 지점(${JSON.stringify(startMarker)})을 못 찾음 — 원본 구조가 바뀐 것 같습니다.`);
  const endIdx = html.indexOf(endMarker, startIdx);
  if(endIdx < 0) throw new Error(`[make-blank-html] "${label}" 끝 지점(${JSON.stringify(endMarker)})을 못 찾음 — 원본 구조가 바뀐 것 같습니다.`);
  return { startIdx, endIdx, text: html.slice(startIdx, endIdx) };
}

function replaceBetween(html, startMarker, endMarker, replacement, label){
  const { startIdx, endIdx } = sliceBetween(html, startMarker, endMarker, label);
  return html.slice(0, startIdx) + replacement + html.slice(endIdx);
}

let replacedCount = 0;

// 1) 경비지도사 목록 비우기
html = replaceBetween(
  html,
  'function seedGuides(){',
  'function seedBranches(){',
  'function seedGuides(){\n  return [];\n}\n',
  'seedGuides'
);
replacedCount++;

// 2) 지사 목록 비우기(한 줄짜리 함수)
{
  const before = html;
  html = html.replace(
    /function seedBranches\(\)\{ return \[[^\]]*\]; \}/,
    'function seedBranches(){ return []; }'
  );
  if(html === before) throw new Error('[make-blank-html] "seedBranches" 한 줄짜리 패턴을 못 찾음.');
  replacedCount++;
}

// 3) 순회일정 시드 레코드 비우기
html = replaceBetween(
  html,
  'function seedRecords(){',
  'const RAW_EMPLOYEES_2026_08',
  'function seedRecords(){\n  return [];\n}\n',
  'seedRecords'
);
replacedCount++;

// 4) 직원 원본 데이터(대용량 한 줄짜리 배열) 비우기
html = replaceBetween(
  html,
  'const RAW_EMPLOYEES_2026_08',
  'function seedEmployees(){',
  'const RAW_EMPLOYEES_2026_08 = [];\n',
  'RAW_EMPLOYEES_2026_08'
);
replacedCount++;

fs.writeFileSync(outPath, html);
console.log(`빈 데이터 버전 생성 완료(${replacedCount}개 항목 비움): ${outPath}`);
