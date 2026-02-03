# 変更計画: worktree パス表示の誤解析対応

## 目的
- git branch -vv の worktree パス表示でブランチ情報が崩れる原因を整理する
- (worktree path) を含む行でも Remote を正しく抽出できるようにする

## 方針
- parseAndPrintLine の正規表現を worktree path を任意要素として扱える形に拡張する
- worktree path は末尾ディレクトリ名だけを抽出して表示する
- 既存出力フォーマットは大枠を維持し、抽出と表示位置のみを調整する

## 手順
1. git branch -vv の出力仕様（括弧付きパス）を前提にパーサの正規表現を整理する
2. 正規表現を更新し、(path) の有無にかかわらず Remote を抽出できるようにする
3. worktree path から末尾ディレクトリ名のみを取得する
4. 既存の出力整形ロジックへの影響を確認する

## 実装概要
- parseAndPrintLine の正規表現を (... ) を任意グループとして追加する
- worktree path は / と \\ の両方を区切りとして末尾ディレクトリ名に変換する
- Remote 抽出は従来通り [...] の内容から行う
- commit message には worktree パスが混ざらないことを確認する

## フォーマット補足
- git branch -vv は `<marker><branch> <commit> (<worktree-path>) [<upstream>: <ahead/behind>] <subject>` の順で出力される
- `(<worktree-path>)` は別 worktree で checkout 中のブランチにだけ付く
- `<marker>` は `*` (現在ブランチ) / `+` (他 worktree で checkout 中) / なし
- 例: `wpf/dev (C:/repos/dev) [origin/hoge/dev] Merged PR 26149: "パラメータ"`

## 変更後の例
- `wpf/dev (dev) Merged PR 26149: "パラメータ" - origin/hoge/dev`

## 影響範囲
- main_gbranch.go

## 確認事項
- worktree path の表示形式が他にも存在するか
- Remote がないブランチの表示で副作用がないか
