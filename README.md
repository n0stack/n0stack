# n0stack

[![Build Status](https://travis-ci.org/n0stack/proto.svg?branch=master)](https://travis-ci.org/n0stack/proto)

## Motivation

現在のOpenstackやKubernetesは以下のような問題を抱えていると考えている。

- Openstackやkubernetesは汎用性を求めすぎてオプションの量がしんどい (k8sはまともだけど)
  - もちろん開発者が多く、いいプロダクトだと思うのですみわけをしていきたい
    - Openstackは構築コストが低い検証環境のIaaS基盤として
    - Kubernetesは...
- プロダクション環境は結局自社に合わせたインフラが必要になるが、たいていは完全に合致しないためある程度書き換えないといけない
  - 結局ソースコードを全部理解するような人が必要で、そこからいい感じに書き換えるか互換のあるものを作るというコストは結構高い
  - 設定を四苦八苦するよりもわかりやすくて短いソースを読んだほうが理解できるし、何なら自分で書いたほうがいい

以上を踏まえ我々は、 *最初から自分で最高に合致するものを書けばよい* という結論に至った。よって開発するにあたって以下のことを前提とする。

- コードは短く、シンプルに、汎用性の一切を排除
  - オプションも最小限に
- gRPCでインターフェイスはしっかり定義してやる
  - OSSの再利用性を確保するため
    - 例えばComputeの実装はそのまま使えるけどNetworkだけ機材が違うから書き換えたいなど
  - provisioning(リソースの仮想化)、middleware(SaaS)、application(アプリケーションの抽象化)、service(サービスの抽象化)など抽象化するレベルを分ける
- 我々が作るメインストリームはあくまでモデルケース
  - 変える必要があれば各自がフォークして良しなにやること

加えて、個人的に楽しいので可能な限り新しく面白い技術を採用していく。

## 原則

- あくまでprotoファイルにはどんな環境でも共有であることが自明な値のみを定義する
  - 実装に依存する値は `annotations` にいれる
  - 例えば、ボリュームのURLなどは実装に依存する

### 標準フィールド

- Metadata (1 ~ 9): 管理に必要な情報をいれる
- Spec   (10 ~ 49): 実装に依存しない汎用的なモデルの設計をいれる
- Status    (50 ~): 実装に依存しない汎用的なモデルの状態をいれる

```pb
  // Make unique!!
  string name = 1;
  // string namespace = 2;

  map<string, string> annotations = 3;
  // map<string, string> labels = 4;

  // Version is used for optimistic lock.
  // Start from 0.
  // Getした値とApplyする値が一致した場合、受け入れられる
  uint64 version = 5;
```

#### 書き込みの優先順位

|  | API | User | Component |
| -- | -- | -- | -- |
| Metadata | 2 | 1 | - |
| Spec | - | 1 | - |
| Status | 1 | - | - |
