# これは何？

Azure の Scheduled Event をチェックして、Pixela にポストするツールです。

# コンパイル

```console
make build
```

# 使い方

## 1. Pixela のユーザーとグラフの作成

### ユーザーの作成

https://docs.pixe.la/entry/post-user

### グラフの作成

https://docs.pixe.la/entry/post-graph

## 2. 設定ファイルの作成

Pixela の設定を `config.json` に記述します。

例:
```json
{"UserID":"tsunomur", "GraphID":"graph1", "Secret": "secret"}
```

## 3. ツールの実行

```bash
./sechecker
```