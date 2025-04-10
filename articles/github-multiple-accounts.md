---
title: "複数のGitHubアカウントを自動で切り替える方法"
emoji: "🔀"
type: "tech"
topics: ["git", "github", "ssh", "zsh"]
published: true
---

## 前提

複数のGitHubのアカウントを運用したいとき、以下の2つを考慮する必要があります。

- `git config` の `user.name` と `user.email`: これはコミットの署名者情報です。
- GitHubと通信を実施するためのSSHキー：リポジトリのアクセス権限でありSSH接続をするために必要です。

複数アカウントで楽に切り替えて、SSH認証情報はあっているけど、`user.name` と `user.email`が別アカウントのものだった！を回避するために記載します。

実際SSH認証情報が合っていればコミット等の操作は出来てしまいます。

- ディレクトリごとに自動で`git config`を切り替えるように設定する
- コマンドを1回入力して、SSH認証情報の切り替えを実施する

を順番に解説します。

## ディレクトリごとに自動で`git config`を切り替える

個人、会社のプロジェクトがあるディレクトリがある状態を作ります。

- ~/company（会社）
- ~/personal（個人）

**1. ~/.gitconfig(グローバルな設定ファイル)** を編集

`~/.gitconfig` に `includeIf` を設定することで、特定のディレクトリ以下でのみ別の設定ファイルを読み込むようにできます。

```ini
[includeIf "gitdir:~/company/"]
    path = ~/.gitconfig_company

[includeIf "gitdir:~/personal/"]
    path = ~/.gitconfig_personal

```

**2. ~/.gitconfig_personal を作成**

```ini
[user]
    name = 個人のGitHubのユーザー名
    email = 個人のGitHubのメールアドレス
```

**3. ~/.gitconfig_company を作成**

```ini
[user]
    name = 会社のGitHubのユーザー名
    email = 会社のGitHubのメールアドレス
```

## コマンドを1回入力して、SSH認証情報の切り替えを実施する

複数のSSHキーで切り替える

1. ~/.ssh/config の設定を確認

```sshconfig
# 会社用
Host company
  HostName github.com
  IdentityFile ~/.ssh/company
  User git
  Port 22
  TCPKeepAlive yes
  IdentitiesOnly yes

# 個人用
Host personal
  HostName github.com
  IdentityFile ~/.ssh/personal
  User git
  Port 22
  TCPKeepAlive yes
  IdentitiesOnly yes
```

2. ssh-agentを自動切り替えをするようにaliasを作成してコマンドで切り替える

ssh-agent：SSHの秘密鍵を一時的にメモリに保持してくれるプログラム

Macであれば、`.zshrc`に以下を追加します。

```shell
alias company = 'ssh-add -D && ssh-add ~/.ssh/company && echo "Switched to company"'
alias personal = 'ssh-add -D && ssh-add ~/.ssh/personal && echo "Switched to personal"'
```

設定後シェルの再起動を実施

```shell
source ~/.zshrc
```

## 実施に動かす

こちらでターミナルを使用する際に、
`company`と入力すれば、SSH認証情報が会社用に
`personal`と入力すれば、SSH認証情報が個人用に
なります。

SSH認証情報の切り替え忘れが発生しても、リポジトリとはそもそも通信が出来ないので問題ないです。

たとえば、`personal`と入力すると

```shell
All identities removed.
Identity added: /Users/sho/.ssh/personal (個人のGitHubのメールアドレス)
Switched to personal
```

このようなログが表示されます。
