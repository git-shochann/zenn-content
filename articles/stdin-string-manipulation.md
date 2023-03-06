---
title: "コーディング試験で解いた文字列操作を細かく解説してみる"
emoji: "😽"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go", "golang"] # タグ名
published: false # 公開設定
---

## 初めに

とあるコーディング試験で以下のような問題が出ました。

Q: 文章中にアルファベット大文字または数字から始まる単語が何種類あるか数えてください。

条件: 単語の個数なので重複はカウントしません。

入力は標準入力から取得して、最後に出力します。

ex) input

`Favorite food is yakiniku. Age is 25 years old. Favorite hobby is coding.`

この問題を解く際に学んだことが沢山あったので、まとめました。

## 最初に考えたこと

1. まずは空白とピリオドで文字を区切って新しくスライスを作成する
2. 重複を判定する関数を用意して、重複を取り除いた新しいスライスを作成する
3. スライスを for ~ range で回していき、文字を取得して、大文字か小文字か判定する

このように最初に整理して考えました。

最初に全てのコードを提示しても分かり辛いので、main 関数から最後まで切り分けて解説します。

## エントリーポイント (main 関数)

```go
func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    text := scanner.Text()
    fmt.Printf("Solution(text): %v\n", Solution(text))
}
```

`func NewScanner(r io.Reader) *Scanner`

Scanner 型は、標準入力やファイルからのデータの読み込みや、文字列を指定の区切りで分割する機能を提供してくれます。
引数には io.Reader 型 => io.Reader interface は、読み込み可能なデータを表現するためのインターフェースを提供します。
なので読み取り可能なデータを引数にとってあげます。

ここで設定する os.Stdin は、標準入力を表す os.File 型のオブジェクトであり、
`Stdin = NewFile(uintptr(syscall.Stdin), "/dev/stdin")`
`func NewFile(fd uintptr, name string) *File`です。

bufio.NewScanner()に渡すことで、標準入力から読み込みを行うよと最初に設定します。

- オブジェクト
  「Scanner 型のオブジェクトを作成する」という表現において、この「オブジェクト」とは変数や関数、構造体などの実体のことです。
  要するに、構造体のフィールド等を埋める作業を行い、実体を作成するということになります。

`scanner.Scan()`

scanner.Scan()は、実際の標準入力からデータを読み取り、そのデータを空白文字で区切って 1 つのトークンとします。

そのトークンは scanner オブジェクトの token フィールドに一時的に保持します。

構造体のフィールドにトークンとして保存するだけなので、特に戻り値の利用はない処理になります。

`text := scanner.Text()`

scanner.Scan()で scanner 構造体の token フィールドを読み取り、読み込んだトークンを文字列として返します。

## 実際に処理する部分

```go
func Solution(str string) int {

    slice := strings.FieldsFunc(str, func(r rune) bool {
        return unicode.IsSpace(r) || r == '.'
    })
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除
    for _, v := range uniqueSlice {
        initial := v[0:1]
    ok := CheckRegex(initial)
    if ok {
        count++
    }
    }
    return count
}
```

```go
    slice := strings.FieldsFunc(str, func(r rune) bool {
        return unicode.IsSpace(r) || r == '.'
    })
```

`func FieldsFunc(s string, f func(rune) bool) []string {}`

ここではまず 1 つ 1 つの単語にアクセスして、空白とピリオドで分割するように設定します。
無名関数が引数の場合はその通りにシグニチャを用意すれば OK です。
戻り値として空白とピリオドを除いたスライスを返します。

その結果の変数 slice は以下のような値を持っています。

`["Favorite" "food" "is" "yakiniku" "Age" "is" "25" "years" "old" "Favorite" "hobby" "is" "coding"]`

その後、重複を削除する関数が呼び出されます。

```go
func DeleteDuplicate(strings []string) []string {
    var unique []string
    m := map[string]bool{}
    for _, v := range strings {
        if len(word) == 0 {
            continue
        }
        if _, ok := m[v]; !ok {
            m[v] = true
            unique = append(unique, v)
        }
    }
    return unique
}
```

まず重複を除いたスライスを新しく作成したいので`var unique []string`と宣言します。

ここでのポイントはどう重複しているかどうかを判定するかということです。
`m := map[string]bool{}`という変数を用意して、こちらを使います。
まず map 型の性質として、同じキーが保存されることはありません。

参考例

```go
    m := make(map[string]int)
    m.key1 = 1
    m.key2 = 2
    m.key1 = 5 // 上書きする

    fmt.Println(m)
    // -> [key1: 5, key2: 2]
```

また空の文字列が含まれる場合があるので、以下はそれをスルーする記述です。
例えば最初の文字列が空白だったり、ピリオドで終わった文章で空の文字列が入ってしまうことがあるからです。

```go
    if len(word) == 0 {
        continue
    }
```

その後、区切られた文字列のスライスをそれぞれアクセスしていって、重複がどうかを判定します。
`m[v] -> m["Favorite"]`とキーにアクセスし、map のバリューが true か false で判定します。
要するに true であれば既にキーは存在しているということになります。
このように単語分繰り返していき、重複を除いたスライスを 戻り値として返します。

```go
func Solution(str string) int {

    slice := strings.Split(str, " ")
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除

    for _, v := range uniqueSlice {
        initial := v[0:1]
    ok := CheckRegex(initial)
    if ok {
        count++
    }
    }
    return count
}
```

さらに重複を除いたスライスを for range で回して、v[0:1]として頭文字にアクセスします。
そして CheckRegex()として関数を用意します。

```go
regex := regexp.MustCompile(`[A-Z0-9]`) // 準備

func CheckRegex(s string) bool {
    ok := regex.MatchString(s)  // 判定
    return ok
}
```

ここで判定を行います。尚 `regex := regexp.MustCompile(`[A-Z0-9]`)` は一度のみの実行でいいので CheckRegex()外で宣言しています。
regex.MatchString()で true であれば大文字 or 数字 としてカウントして最後に出力して終了です。

### 全体のコード

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
)

var count int

func Solution(str string) int {
    slice := strings.FieldsFunc(str, func(r rune) bool {
        return unicode.IsSpace(r) || r == '.'
    })
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除
    for _, v := range uniqueSlice {
        initial := v[0:1]
        ok := CheckRegex(initial)
        if ok {
            count++
        }
    }
    return count
}

func DeleteDuplicate(strings []string) []string {
    // m[""]false が初期化した際の型
    m := make(map[string]bool)

    var unique []string
    for _, v := range strings {
        if len(word) == 0 {
            continue
        }
        // m[v]がtrueでなければ = まだそのキーはないということ
        if _, ok := m[v]; !ok {
            m[v] = true
            unique = append(unique, v)
        }
    }
    return unique

}

var regex = regexp.MustCompile(`[A-Z0-9]`) // 準備

func CheckRegex(s string) bool {
    ok := regex.MatchString(s) // 判定
    return ok
}
```

### 上記の改善点

一応上記でのコードでも動くのですが、

`strings.FieldsFunc()`でも全てを読み込むとなると、テキストファイルのサイズが大きいと変数の容量も増加し、
変数のサイズが大きくなるとメモリ確保やコピーなどに時間がかかります。

今回全てのテキスト情報を保持しなくても、文字列をスペース等が来るまで順番に読み取って処理を行い、
単語の重複チェックを行なっていくことで既に処理した文字列を破棄することが出来るようなロジックでも書くことが出来ました。

### 方法 2

```go
    scanner := bufio.NewScanner(os.Stdin)

    // どのように区切るかの設定を行う
    scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
        // 空白で区切る
        for i := 0; i < len(data); i++ {
            if data[i] == ' ' || data[i] == '.' {
                return i + 1, data[:i], nil
            }
        }
        return 0, data, bufio.ErrFinalToken
    })
```

まず区切る設定を行います。Split()には無名関数を渡します。

```go
    words := map[string]struct{}{}
    for scanner.Scan() {
        word := scanner.Text()
        if len(word) == 0 {
            continue
        }
        capital := rune(word[0]) // runeで比較したいので変換する
        if unicode.IsUpper(capital) || unicode.IsDigit(capital) {
            // 一致したらキーにその単語を、値は空で設定する
            words[word] = struct{}{}
        }
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
    fmt.Printf("%d words\n\n", len(words))
    // mapをfor range で回す
    for word := range words {
        fmt.Println(word)
    }
```

最初に`words := map[string]struct{}{}`とありますが、こちらも重複をチェックするために用意しています。後ほど説明します。

まず、scanner.Scan()は標準入力空白で区切って(空白以外の区切りにも出来るがデフォルトは空白)、データを空白文字で区切って 1 つのトークンとします。
そのトークンは scanner オブジェクトの token フィールドに一時的に保持します。なのでこれを利用して、1 つ 1 つ文字を見ていけばいいだけです。

`word[0]`とアクセスすると、1 文字が byte 型で返ってくるため、rune 型に変換します。

rune 型に関してですが、
コンピューターは文字を直接扱うことが出来ません。実行するときも 0 と 1 の 2 進数の機械語を用いて、実行します。
なので "あ" という文字に対応した数字を用意して(Unicode コードポイント)、2 進数に変換しようという仕組みが必要です。
それが Unicode といった文字コードになります。

例えば、Unicode コードポイントとして 16 進数の "0041" が割り当てられた文字 "A" を扱う場合、コンピューターは以下のように処理します。

1. Unicode コードポイント "0041" を 2 進数表現に変換する
2. 2 進数表現をコンピューターが扱える形式に変換する

コンピューターは、16 進数や 10 進数などの数値表現を内部的に 2 進数に変換して扱います。そのため、2 進数表現 "0000 0000 0100 0001" は、コンピューターが扱える形式であり、機械が理解できる形式です。

`A -> 0041(16進数/Unicode/rune型) → 0000 0000 0100 0001(2進数)`

要するに上記のような流れになります。なので Unicode は必要です。

`unicode.isUpper()`と`unicode.IsDegit()`は rune を引数に取って比較が出来るので、

そして条件が true であれば、以下のようになります。

`words[word] = struct{}{}`

まず解説すると、

`words := map[string]struct{}{}`

まず map のキーに string、バリューには構造体リテラルを用いて、宣言しています。

まず空の構造体型として`struct{}`という型があると認識すれば OK です。そして値も一緒に構造体リテラルという形式で
初期化を行っています。

そもそも構造体リテラルというのは、初期化(宣言された変数に初めて値を代入すること)をする際に構造体型の値を一緒に設定するための方法です。

そしてこれがどう重複判定に利用できるかという部分を解説していきます。
特に今回の場合バリューにアクセスする必要が別にないのでキーだけで判断しようということです。
これにより、map のキーの存在を調べるときにメモリ使用量が大幅に削減されます。

言葉ではイメージつかないと思うのでコードで示します。

```go
    words := map[string]struct{}{}

    words["key1"] = struct{}{}
    words["key2"] = struct{}{}
    words["key1"] = struct{}{}

    fmt.Println(words)
    // -> map[key1:{} key2:{}]
```

map のキーは重複することがありません。なので同じキー名が来ても上書きされるだけなのでその仕組みを利用しているだけです。
`words := map[string]struct{}{}`

`words[word] = struct{}{}`
一致したらキーにその単語を、値は空で設定します。
words という map の集合体に word というキーを設定し、バリューは空です。

最後に map の集合体を for range で回して出力します。

### 全体のコード

```go

func otherSolution() {
    scanner := bufio.NewScanner(os.Stdin)

    // どのように区切るかの設定を行う -> 空白とピリオドで区切る
    scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
        for i := 0; i < len(data); i++ {
            if data[i] == ' ' || data[i] == '.' {
                return i + 1, data[:i], nil
            }
        }
        return 0, data, bufio.ErrFinalToken
    })

    words := map[string]struct{}{}
    for scanner.Scan() {
        word := scanner.Text()
        if len(word) == 0 {
            continue
        }
        capital := rune(word[0]) // runeで比較したいので変換する
        if unicode.IsUpper(capital) || unicode.IsDigit(capital) {
        // 一致したらキーにその単語を、値は空で設定する
            words[word] = struct{}{}
        }
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
    fmt.Printf("%d words\n\n", len(words))
    // mapをfor range で回す
    for word := range words {
        fmt.Println(word)
    }
}
```
