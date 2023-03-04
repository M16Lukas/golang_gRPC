## gRPC

### gRPCとは
Googleによって2015年にオープンソース化されたRPC（Remote Procedure Call）のためのプロトコル

### RPC（Remote Procedure Call）とは
- ネットワーク上の他の端末と通信するための仕組み
- REST APIのようにパスやメソッドを指定する必要はなく、メソッド名と引数を指定する
- gRPC以外にJSON-RPCなどがあるが、今はgRPCがデフォルトスタンダード

### gRPCの特徴
- データフォーマットにProtocol Buffersを使用
  - バイナリにシリアライズすることで送信データ量が減り高速な通信を実現
  - 型付けされたデータ転送が可能
- IDL（Protocol Buffers）からサーバー側・クライアント側に必要なソースコードを生成
- 通信にはHTTP/2を使用
- 特定の言語やプラットフォームに依存しない

### gRPCが適したケース
- Microservice間の通信
  - 複数の言語やプラットフォームで構成される可能性がある
  - バックエンド間であれば、gRPCの恩恵が多く得られる
- モバイルユーザーが利用するサービス
  - 通信量が削減できるため、通信容量制限にかかりにくい
- 速度が求められる場合

### gRPCの開発の流れ
1. protoファイルの作成
1. protoファイルをコンパイルしてサーバー・クライアントの雛形コードを作成
1. 雛形コードを使用してサーバー・クライアントを実装

---

## HTTP/2

### HTTP/1.1の課題
- リクエストの多量化
  - １リクエストに対して１レスポンスという制約があり、大量のリソースで構成されているページを表示するには大きなネットになる
- プロトコルオーバヘッド
  - Cookieやトークンなどを毎回リクエストヘッダに付与してリクエストするため、オーバヘッドが大きくなる

### HTTP/1.1の特徴
- ストリームという概念を導入
  - 1つのTCP接続を用いて、複数のリクエスト・レスポンスのやり取りが可能
  - TCP接続を減らすことができるので、サーバーの負荷軽減
- ヘッダーの圧縮
  - ヘッダーをHPACKという圧縮方式で圧縮し、さらにキャッシュを行うことで、差分のみの送受信することで効率化
- サーバープッシュ
  - クライアントからのリクエスト無しにサーバーからデータを送信できる
  - 事前に必要と思われるリソースを送信しておくことで、ラウンドトリップの回数を削減し、リソース読み込みまでの時間を短縮

### Demo
http://www.http2demo.

---

## Service

### Serviceとは
- RPC(メソッド)の実装単位
  - サービス内に定義するメソッドがエンドポイントになる
  - １サービス内に複数のメソッドを定義できる
- サービス名、メソッド名、引数、戻り値を定義する必要がある
- コンパイルしてGoファイルに変換すると、インターフェースとなる
  - アプリケーション側でこのインターフェースを実装する

### Serviceのサンプル

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (SayHelloRequest) returns (SayHelloResponse);
}
```

---

## gRPCの通信方式

### ４種類の通信方式
- Unary RPC
- Server Streaming RPC
- Client Streaming RPC
- Bidirectional Streaming RPC

### Unary RPC
- 1リクエスト1レスポンスの方式
- 通信の関数コールのように扱うことができる
- 用途
  - APIなど
- Service定義

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (SayHelloRequest) returns (SayHelloResponse);
}
```

### Server Streaming RPC
- １リクエスト・複数レスポンスの方式
- クライアントはサーバーから送信完了の信号が送信されるまでストリームのメッセージを読み続ける
- 用途
  - サーバーからのプッシュ通知など
- Service定義

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (SayHelloRequest) returns (stream SayHelloResponse);
}
```

### Client Streaming RPC
- 複数リクエスト・１レスポンスの方式
- サーバーはクライアントからリクエスト完了の信号が送信されるまでストリームからメッセージを読み続け、レスポンスを返さない
- 用途
  - 大きなファイルのアップロードなど
- Service定義

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (stream SayHelloRequest) returns(SayHelloResponse);
}
```

### Bidirectional Streaming RPC
- 複数リクエスト・複数レスポンスの方式
- クライアントとサーバーのストリームが独立しており、リクエストとレスポンスはどのような順序でもよい
- 用途
  - チャットやオンライン対戦ゲームなど
- Service定義

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (stream SayHelloRequest) returns (stream SayHelloResponse);
}
```
