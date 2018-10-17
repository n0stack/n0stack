# proto

Protobuf definitions for all of n0stack services

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
