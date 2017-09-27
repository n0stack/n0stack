# CONTRIBUTING

## 開発方針

- pull request駆動
  - 2人以上のレビューをもらうこと
- ユニットテストはとりあえず作らない
- annotationは100%を目指す

## About documents

## document

`/docs` 配下にsphinx形式で可能な限り書くこと

### annotations

- python2のmypy形式を用いる
  - http://mypy.readthedocs.io/en/latest/cheat_sheet.html
- 理由としては実績があり、コメント形式であるため仕様が変わっても柔軟に対応できるため

### docstring

- docstringを[Google形式](http://google.github.io/styleguide/pyguide.html#Comments)で書く
  - private methodについてはdocstringを書かなくても良い
  - typeについてはmypyとかぶるので書かない

#### Example

```python
def create_bridge(cls, bridge_name, interface_name):
    # type: (str, str) -> Tuple[int, int]
    """Create bridge mastering interface selected on args

    1. Create a bridge.
        When the bridge already exists, ignoring.
        in command: `ip link add dev $bridge_name type bridge`
    2. Set interface master to the bridge.
        in command: `ip link set dev $interface_name master $bridge_name`
    3. Set up the interface and the bridge.
        in command: `ip link set up dev $names`

    Args:
        bridge_name: Name of creating bridge.
        interface_name: Name of bridge slave interface.

    Returns:
        Tuple
        - Index of bridge interface.
        - Index of interface interface.
    """
```

## CI

### tools

- flake8
- mypy
- (pylint)

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
