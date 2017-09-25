# CONTRIBUTING

## 開発方針

- pull request駆動
- ユニットテストはとりあえず作らない
- annotationは全てにつける

## About documents

### docstring

- docstringを[Google形式](http://google.github.io/styleguide/pyguide.html#Comments)で書く
  - private methodについてはdocstringを書かなくても良い

## CI

### tools

- pylint
- flake8
- mypy

### docs

- sphinxで自動生成する

## ブランチの扱い

### ブランチ名

`[process名]/[任意の単語]`

- `[任意の単語]` には可能な限り名詞を使うこと
- スネークケースを使うこと

### master

- 動作する状態

### development

- ファーストリリースまではメインリポジトリに
