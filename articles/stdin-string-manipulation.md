---
title: "Goで文字列操作を行う際の具体的な処理を細かいところまで追ってみる"
emoji: "😽"
type: "tech" # tech: 技術記事 / idea: アイデア
topics: ["go", "golang"] # タグ名
published: false # 公開設定
---

## 初めに

とある問題集で以下のような問題がありました。

文章中にアルファベット大文字または数字から始まる単語が何種類あるか数えてください。

ただし単語の個数なので重複はカウントしません。

入力は標準入力から取得して、最後に出力します。

ex) input

`Favorite food is yakiniku. Age is 25 years old. Favorite hobby is coding.`

この問題を解く際に学んだことが沢山あったので、まとめました。

## 最初に考えたこと

1. まずは空白で文字を区切って新しくスライスを作成する
2. 重複を判定する関数を用意して、重複を取り除いた新しいスライスを作成する
3. スライスを range で回していき、文字を取得して、大文字か小文字か判定する

このように最初に整理して、ロジックを考えました。

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
ここで設定する os.Stdin は、標準入力を表す os.File 型のオブジェクトであり、`Stdin = NewFile(uintptr(syscall.Stdin), "/dev/stdin")`

bufio.NewScanner()に渡すことで、標準入力から読み込みを行うよと最初に設定します。

- オブジェクト
  「Scanner 型のオブジェクトを作成する」という表現において、「オブジェクト」とは、プログラミングにおける概念の 1 つであり、変数や関数、構造体などの実体のこと
  要するに、構造体のフィールド等を埋める作業を行い、実体を作成するということ

`scanner.Scan()`

scanner.Scan()は、実際の標準入力からデータを読み取り、そのデータを空白文字で区切って 1 つのトークンにする。

そのトークンは scanner オブジェクトの token フィールドに一時的に保持します。
例えば、scanner.Scan()がまた次に呼ばれると、新しいトークンを読み取って、token フィールドを上書きします。

構造体のフィールドにトークンとして分割して保存するだけなので、特に戻り値の利用はない処理になります。

`text := scanner.Text()`

scanner.Scan()で scanner 構造体の token フィールドを読み取り、読み込んだトークンを文字列として返します。

## 実際に処理する部分

```go
func Solution(str string) int {

    slice := strings.Split(str, " ")
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除
    for _, v := range uniqueSlice {
        // [M3 2000 sho tsuboya]
        initial := v[0:1]
    ok := CheckRegex(initial)
    if ok {
        count++
    }
    }
    return count
}
```

`slice := strings.Split(str, " ")`
ここではまず、1 つ 1 つの単語にアクセスするため、長"い 1 つの文字列を空白単位で Go のスライスに変換します。

上記の変数 slice は以下のような値を持っています。
`["Favorite" "food" "is" "yakiniku." "Age" "is" "25" "years" "old." "Favorite" "hobby" "is" "coding."]`

その後、重複の関数が呼び出されます。

```go

func DeleteDuplicate(strings []string) []string {
    var unique []string
    m := map[string]bool{}
    for _, v := range strings {
    if _, ok := m[v]; !ok {
        m[v] = true
        unique = append(unique, v)
    }
    return unique
}

```

まず重複を除いたスライスを新しく作成したいので`var unique []string`と宣言します。

ここでのポイントはどう重複しているかどうかを判定するかということです。
`m := map[string]bool{}`という変数を用意して、こちらを使います。
まず map 型の性質として、同じキーが保存されることはありません。

```go
    m := make(map[string]int)
    m.key1 = 1
    m.key2 = 2
    m.key1 = 5 // 上書きする

    fmt.Println(m)
    // -> [key1: 5, key2: 2]
```

空白で区切られた文字列のスライスをそれぞれアクセスしていって、重複がどうかを判定します。
m[v] -> m["Favorite"] とキーにアクセスし、map のバリューが true か false で判定します。
要するに true であれば既にキーは存在しているということになります。
上記の操作を単語分繰り返し重複を除いたスライスを 戻り値として返します。

```go

func Solution(str string) int {

    slice := strings.Split(str, " ")
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除

    // ここまでOK

    for _, v := range uniqueSlice {
        // [M3 2000 sho tsuboya]
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

## 全体のコード

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

    slice := strings.Split(str, " ")
    fmt.Printf("slice: %v\n", slice)
    uniqueSlice := DeleteDuplicate(slice) // 重複の削除
    fmt.Printf("uniqueSlice: %v\n", uniqueSlice)
    for _, v := range uniqueSlice {
        // [M3 2000 sho tsuboya]
        initial := v[0:1]
        ok := CheckRegex(initial)
        if ok {
            count++
        }
    }
    return count
}

// [M3 2000 sho tsuboya M3]
func DeleteDuplicate(strings []string) []string {
    // m[""]false が初期化した際の型
    m := make(map[string]bool)

    // 一応map[string]struct{}{}の方がいいかもしれない 値にアクセスしなくてもいいので
    // 詳細は # other_solution.goを参照

    var unique []string
    for _, v := range strings {
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
