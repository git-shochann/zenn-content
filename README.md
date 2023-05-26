# Zenn CLI

* [📘 How to use](https://zenn.dev/zenn/articles/zenn-cli-guide)

```npx zenn new:article```
articles/ランダムなslug.mdというファイルが作成される。slug（スラッグ）はその記事のユニークな ID のようなもの。

```
---
title: "" # 記事のタイトル
emoji: "😸" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # tech: 技術記事 / idea: アイデア記事
topics: [] # タグ。["markdown", "rust", "aws"]のように指定する
published: true # 公開設定（falseにすると下書き）
---
ここから本文を書く
```

```$ npx zenn preview```
プレビュー開始