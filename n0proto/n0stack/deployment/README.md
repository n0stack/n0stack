# Deployment

- Provisionin を抽象化し、宣言的に定義できるようにする
  - よって非同期も可
- この層は immutable なリソースを扱う

## immutable であるということは

- 宣言的に定義する
- よってインターフェイスは ARD
  - Apply, Read, Delete
