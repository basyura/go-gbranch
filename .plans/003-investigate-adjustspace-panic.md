# 変更計画: adjustSpace の panic 調査とテスト追加

## 目的
- strings.Repeat の panic 原因を特定し、再発防止を保証する
- adjustSpace の境界条件をテストで固定化する

## 方針
- 再現条件を明文化し、テストで回帰防止する
- ユニットテストで ASCII / 非 ASCII 混在の表示幅を検証する
- 既存関数の公開範囲は変えず、同一パッケージ内で直接呼び出す

## 手順
1. panic 再現例をテストケースに落とし込む
2. main_gbranch_test.go を追加して adjustSpace を検証する
3. go test ./... でローカル実行できる状態にする
4. go vet で失敗しないよう、色付け出力のフォーマット処理を見直す

## 実装概要
- テーブル駆動で ASCII / 非 ASCII / 境界値のケースを用意する
- panic 再現例（例: "あい" と maxLen=3 ）で panic しないことを確認する
- 必要に応じて strLen を利用し、期待表示幅を検証する
- 色付け出力は Sprintf ではなく Sprint を用いてフォーマット解釈を回避する

## 原因候補と再現例
- adjustedLength の判定が加算前のみのため、最後の文字で maxLen を超えると return されず末尾に落ちる
- 例: adjustSpace("あい", 3) は strLen("あい") が 4 のためループ後に strings.Repeat(" ", -1) となり panic

## 修正案の説明（案 1）
- 調整済みの幅で列揃えを維持するため、短い文字列は右パディングで幅を合わせる
- 長い文字列は chopRightByWidth で切り詰め、短い文字列は padRightByWidth で埋める

## 影響範囲
- main_gbranch_test.go

## 確認事項
- 期待値は「表示幅基準」で良いか
- commitMsgLength=70 の上限をテストに含めるか
