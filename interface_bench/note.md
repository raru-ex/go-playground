# メモ

## 考えたいこと

Factory(NewHoge)はInterfaceと実体のどちらを返すべきなのか

## 結論

まだしっかり言語化してないが、多分基本は以下になりそう
- 基本的には実体を返す - ”Accept Interfaces, Return Structs”
  - ただしGo公式もinterfaceを返すFactoryを作っている -> errors.New()
- Interfaceを返してはいけないということはないが、そうすることがメリットになるケースがほとんどない
  - Interfaceを返す場合、そのInterface定義に修正が入るとInterfaceを返している箇所すべてが壊れる(はず)
  - 複数のInterfaceを満たすようなstructではinterfaceで返しづらい
    - return interface{ iA; iB} みたいに気合で返せなくはないが、やる意味が薄い
  - interface経由でのメソッドアクセスのほうが若干パフォーマンスでオーバーヘッドがある
    - おそらく通常は無視できる範囲だが、あえて遅い方を選択するメリットがない
- interfaceのダウンキャストコストは存在するが、通常のシステムにおいてはほぼ無視できる
  - 関数呼び出し自体が倍ほどかかったりしているが、呼び出し時間よりもビジネスロジックのほうが圧倒的に遅いので無視できる誤差になる

## ベンチマークしてみた

### シンプルな構造体・シンプルな演算

プロパティが存在しないStructでSum(a, b int64) int64

#### 予想

interfaceをreturn型としているFactoryで生成されたインスタンスはダウンキャストを行って関数を実行しているのではないかと思うので以下の順番で速いと思う

Pointer = Entity > InterfacePointer

#### 実際に実行してみた結果

```sh
$ go test -bench=BenchmarkInterface -benchmem ./interface_bench/
goos: linux
goarch: amd64
pkg: go-playground/interface_bench
cpu: Intel(R) N95
BenchmarkInterface/InterfacePointer-4           1000000000               0.6240 ns/op          0 B/op          0 allocs/op
BenchmarkInterface/Pointer-4                    1000000000               0.3840 ns/op          0 B/op          0 allocs/op
BenchmarkInterface/Entity-4                     1000000000               0.3747 ns/op          0 B/op          0 allocs/op
PASS
ok      go-playground/interface_bench   1.568s
```

(Entity = Pointer) > InterfacePointer っぽい。
メソッドの定義を以下のようにPointerレシーバにしているので、呼び出しがポインタになっていてEntity, Pointerでのパフォーマンス差は当然なかった。
```go
func (*TestImpl) Sum(a, b int64) int64 {
	return a + b
}
```

### 考察

ダウンキャストによりパフォーマンスが落ちていると言っているが、Interface型のインスタンスはInterface内に実際の型の型情報を持っている(はず)なので、実際には間に挟まっている処理はかなり小さい。
ただ静的呼び出しではコンパイル時に呼び出すものが確定しているため、呼び出しに必要な情報がプリフェッチ出来る。
動的呼び出しでは、CPUなどの推論が上手く動作すれば静的なものと同程度の処理速度が出るが、推測が外れた場合に対象のデータを改めてフェッチするなどするためパフォーマンスに想定より差が出ることがある。(っぽい)
ついでに、動的dispatchはinline化されないという欠点もあるらしい(されてるが...)

## 参考

interfaceを返すことは必ずしも悪いことではない。 https://x.com/mattn_jp/status/1567382261980614663?s=20
あるオブジェクトのメソッド呼び出しはインタフェース経由だとインライン化されないです。https://x.com/mattn_jp/status/1567392612054892544
Goの実装者の言葉に ”Accept Interfaces, Return Structs” がある。
https://zenn.dev/spiegel/articles/20201129-interface-types-in-golang

## Memo

大事そうなワード
- インライン化
  - 関数呼び出し元に呼び出し先の関数をインライン展開して関数呼び出しのオーバーヘッドをなくす
  - プラス展開されたコードがあることで追加の最適化が行えるケースもあるらしい
  - むやみに行うとバイナリが大きくなったり、参照の局所性が下がってパフォーマンスが劣化するケースもある様子
- 動的ディスパッチ
  - 詳細すぎてむずいが、詳細。https://qiita.com/Akatsuki_py/items/e53a4c15513711570469
  - 動的に呼び出す関数をdispatchするという意味っぽいい？
- ポインタのメモリルックアップ
- コンパイラによる最適化

## 知らなかったこと

