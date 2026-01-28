# 変更計画: 端末幅に合わせた右端省略

## 目的
- 端末幅を考慮し、出力が折り返されず右端で省略されるようにする

## 方針
- 端末幅は `golang.org/x/term` の `GetSize` で取得し、失敗時のみ `COLUMNS` を参照する
- 行の表示幅が端末幅を超える場合は、remote を優先的に右端で chop して収める
- まず remote を短くし、それでも収まらない場合は commit message を右端で chop する
- 表示幅の計算は ASCII 1、非 ASCII 2 とする
- `COLUMNS` は数値に変換できた場合のみ採用し、0 以下や不正値は「省略なし」にする

## 手順
1. 端末幅の取得方法を整理し、フォールバックを決める
2. 表示幅計算と右端切り捨ての関数を追加する
3. 出力生成の直前で commit message を必要に応じて切る
4. 既存の整形と色付けが崩れないことを確認する

## 実装概要
- `golang.org/x/term` で端末幅を取得し、失敗時のみ `COLUMNS` を `strconv.Atoi` で取得
- 右端 chop 用の `chopRightByWidth` を追加
- 出力直前に表示幅を計算し、remote → commit message の順で短くする
- `golang.org/x/term` 追加により `go.mod` と `go.sum` を更新する

## 影響範囲
- main_gbranch.go

## 確認事項
- commit message 以外を切る条件の扱い
- 取得できない場合の幅の扱い
