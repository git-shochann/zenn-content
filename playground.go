package main

// You can edit this code!
// Click here and start typing.

func myFunc(chan string) {

}

func main() {
	// 情報伝搬用のチャネルを作成
	info := make(chan string)
	go myFunc(info) // 新規ゴルーチン起動時に、infoチャネルを渡しておく
}
