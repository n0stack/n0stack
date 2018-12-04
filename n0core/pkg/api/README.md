# API

- APIの実装を行う
- ドメイン固有のロジックも含める

## 楽観的なロック

- PENDING ステートを使うことで楽観的なロックを行う
- 副作用のような時間のかかる処理の場合は入れること

```
------------- timeline ------------->

Reqeust 1 --> set PENDING (optimistic lock) --> some process --> set new state --> response OK
              |
              | Request 2 --> get PENDING --> response FailedPrecondition
```
