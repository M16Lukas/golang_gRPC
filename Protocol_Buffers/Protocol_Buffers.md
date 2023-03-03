## Protocol Buffersとは

Googleによって2008年にオープンソース化されたスキーマ言語

> スキーマ言語
> 
> 何かしらの処理をさせるのではなく、
> 
> 要素や属性などの構造を定義するための言語

### スキーマ言語がなぜ重要か
多くのシステムやデータが様々な技術で複数のサービスやストレージに分割されるようになってきている（モノリス→マイクロサービス）

クライアント側もWEBだけではなくiOSやAndroidなどのモバイル対応が必須

スキーマ言語によって事前にどういったデータをやり取りするかを宣言的に定義しておく

### Protocol Buffersの特徴

- **gRPCのデータフォーマット**として使用されている
- **プログラミング言語から独立**しており、様々な言語に変換可能
- **バイナリ形式にシリアライズ**するので、サイズが小さく高速な通信が可能
- **型安全**にデータのやり取りが可能
- **JSONに変換**することも可能

### Protocol Buffersを使用した開発の進め方

1. スキーマの定義（`.proto`）
1. スキーマから開発言語のオブジェクトを自動生成
1. データ通信時にはバイナリ形式へシリアライズ

---

# message とは
- 複数のフィールドを持つことができる型定義
  - それぞれのフィールドはスカラ型もしくはコンポジット型
- 各言語のコードとしてコンパイルした場合、構造体やクラスとして変換される
- 一つの proto ファイルに複数のmessage型を定義することも可能

## messageの例
```proto
message Person {
  string name = 1;
  int32 id = 2;
  string email = 3;
}
```

message キーワードとメッセージ名を定義
```proto
message Person {
```

フィールド
> フィールドの型 + フィールド名 + タグ番号 + ;
```proto
  string name = 1;
  int32 id = 2;
  string email = 3;
```

---

## スカラー型

https://protobuf.dev/programming-guides/proto3/#scalar

---

## Tag
- Protocol Buffersではフィールドはフィールド名ではなく、タグ番号によって識別される
- 重複は許されず、一意である必要がある
- タグの最小値は１、最大値は2^29 - 1(536,870,911)
- `19000 ~ 19999` はProtocol Buffersの予約番号のため使用不可
- `1~15`番までは`1byte`で表すことができるので、よく使うフィールドには`1~15`番を割り当てる
- タグは連番にする必要はないので、あまり使わないフィールドはあえて16番以降を割り当てることも可能
- タグ番号を予約するなど、安全にProtocol Buffersを使用する方法も用意されている

--- 

## デフォルト値
- 定義したメッセージでデータをやり取りする際に、定義したフィールドがセットされていない場合、そのフィールドのデフォルト値が設定される
- デフォルト値は型によって決められている

| type | default |
| :--: | :--: | 
| string | 空の文字列 |
| byte | 空のbyte |
| bool | false |
| 整数型 | 0 |
| 浮動小数点数 | 0 |
| 列挙型 | タグ番号0の値（例：UNKNOWN） |
| repeated | 空のリスト |

---

## protoファイルのコンパイル

```zsh
protoc
```

### `-IPATH, --proto_path=PATH`
- protoファイルのimport文のパスを特定する

```
ディレクトリ構造

Protocol_Buffers
└──proto
    ├──employee.proto
    └─date.proto

```

```proto
import "proto/date.proto";
```
↓
```zsh
protoc　-I.
```

```proto
import "date.proto";
```
↓
```zsh
protoc -I./proto
```

- 複数の箇所からprotoファイルをインポートする必要がある場合、コロン区切りで複数のパスを記述することが可能

```zsh
protoc -I./test:./dev
```

- `-I`オプションを省略した場合は、カレントディレクトリが設定される（= `-I.`）

### 各言語に変換するためのオプション
- オプションによって、どの言語に変換するかを決定する
- Go言語のオプションはプラグインで追加する必要がある

```zsh
--cpp_out=OUT_DIR
--csharp_out=OUT_DIR
--java_out=OUT_DIR
--js_out=OUT_DIR
--kotlin_out=OUT_DIR
--objc_out=OUT_DIR
--php_out=OUT_DIR
--python_out=OUT_DIR
--ruby_out=OUT_DIR
```

### コンパイルするファイルの指定

- 対象のファイルをすべて並べる

```zsh
protoc -I. --go_out=. proto/employee.proto proto/date.proto
```

- ワイルドカード使用

```zsh
protoc -I. --go_out=. proto/*.proto
```
