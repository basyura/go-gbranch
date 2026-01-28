# 変更計画: 現在ブランチ記号の表示と列揃えの修正

## 目的
- 現在のブランチ記号が空白にならず確実に表示されるようにする
- ブランチ名の開始位置が行によってズレないようにする
- 絵文字や合成文字を含む記号でも列揃えが崩れないようにする

## 方針
- 現在ブランチ記号は環境変数が空の場合に必ず `specialSymbol` を使う
- 記号の幅を表示幅で揃えるため、記号は「末尾スペースなし」で保持し、出力時にパディングする
- 列揃えは `len` ではなく表示幅（`strLen`）で計算する
- 表示幅計算で「幅0文字（Mn/Cf）」は幅に加算しない
- `specialSymbol` と `GBRANCH_SYMBOL` の記号は幅2前提でパディング計算する

## 手順
1. `parseAndPrintLine` で現在ブランチ記号の決定ロジックを修正する
2. 記号の保持方針（末尾スペースなし）に合わせて `print` の整形を更新する
3. 記号列の最大幅を算出し、各行の記号を表示幅で右パディングする
4. `strLen` と `chopRightByWidth` の幅計算で幅0文字を除外する
5. 記号列の幅計算を「記号は幅2」とみなすルールに変更する
6. 期待する表示結果を `git branch -vv` 出力例で確認する

## 実装概要
- `GBRANCH_SYMBOL` が空の場合は `specialSymbol` を使い、その後に必要なスペースを足す
- `Symbol` は末尾スペースを含めず保持
- `print` 内で `max_symbol_width` を `strLen` で計算し、`Symbol` をパディングして揃える
- ブランチ名の列揃えも `strLen` を用いる
- `unicode.Is(unicode.Mn, r)` と `unicode.Is(unicode.Cf, r)` の場合は幅を加算しない
- 記号列の表示幅は「記号は幅2」とみなす（`specialSymbol` / `GBRANCH_SYMBOL`）

## 影響範囲
- main_gbranch.go

## 確認事項
- `GBRANCH_SYMBOL` に 1 文字/複数文字/絵文字を指定した場合の表示幅
- `specialSymbol` の幅（非 ASCII）に対する揃えの期待値
- 絵文字が合成文字（例: variation selector）を含む場合の揃え
- 記号幅を2前提にした際の macOS での見た目
