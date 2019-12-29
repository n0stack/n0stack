# Asynchronous Messaging Queue System Architecture

|||
|--|--|
| Status | deprecated |

## Context

OpenStackの問題点を考えた際に、無駄に規模が大きいことが問題であると考えていた。
よって、OpenStackのスケールを小さくしたアーキテクチャを考えれば良いと考えた。

## Decision

以上の経緯から、OpenStackのシステムアーキテクチャを踏襲し、ユーザにはHTTPで非同期APIを提供したうえで、コンポーネント間の通信をMessaging Queueで通信することを考えた。
やろうと思ったことは [これ](https://www.slideshare.net/TakeshiTake/n0stack-81419210) を参照してほしい。

## Consequences

力不足も相まって形にすることができず、 [Synchronous RPC System Architecture](synchronous_rpc_system_architecture.md) に設計方針を変更した。

### Pros

実際に運用できたわけではないため、よくわからないが一般的に言われている以下のようなメリットあるだろう。

- 簡単にコンポーネントをスケールされることができる
- コンポーネント間を粗結合にすることができる

### Cons

まず、MQを運用と構築することが難しい。
KafkaやPulsarなどは非常に大規模なミドルウェアであり、真面目に運用しようとすると多くのコストがかかる。
つまり、n0stackを構築することや運用することが非常に難しくなることを示しており、少なくともすべての基盤となるIaaS部分で使うべきではないと判断した。

また、MQは多くのコンポーネントを運用するために利用ことで真価を発揮する。
しかし、n0stackは構築などを簡単にすることを目指しているため、コンポーネントは少なくしたいと考えている。
そのため、MQを採用しても得られるメリットが少ない。

くわえて、Exactly onceでキューが送られることが保証されるわけではないため、イベントがべき等になるように設計する必要がある。
Virtual Machineの再起動などクラウド基盤ではイベントを正しく処理できなければならない。
それらの設計方法は知見としてあまり共有されておらず、トピックなどの使い方もよくわからず学習コストが高かった。
というか、自分がいつまで立ってもいい感じに設計できなかったため諦めた。

## Reference

- https://www.slideshare.net/TakeshiTake/n0stack-81419210
