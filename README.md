# これは何？

Azure の Scheduled Event をチェックして、Pixela にポストするツールです。

# コンパイル

```console
make build
```

# 使い方(Pixela と連携する例)

## 1. Pixela のユーザーとグラフの作成

### ユーザーの作成

https://docs.pixe.la/entry/post-user

### グラフの作成

https://docs.pixe.la/entry/post-graph

## 2. 設定ファイルの作成

Pixela の設定を `config.json` に記述します。

例:
```json
{
 "command": [
         "curl -X PUT https://pixe.la/v1/users/tsunomurtest/graphs/test-graph/increment -H 'X-USER-TOKEN:thisissecret' -H 'Content-Length:0'"
 ]
}
```

## 3. ツールの実行

```bash
./sechecker
```

## 4. cron への登録例

イベント情報は標準出力、ログは標準エラー出力に出力します。

```bash
* * * * * cd /home/tsunomur/sechecker && ./sechecker >> log 2>&1
```

# 出力例

## 再起動イベント発生時
```bash
$ ./sechecker
2020/07/16 15:37:37 EventState=1
2020/07/16 15:37:37 curl -X PUT https://pixe.la/v1/users/tsunomur/graphs/graph1/increment -H 'X-USER-TOKEN:secret' -H 'Content-Length:0'
2020/07/16 15:37:37 {"message":"Success.","isSuccess":true}
{"DocumentIncarnation":10,"Events":[{"EventId":"CFA22193-3606-45E1-B7EA-976666B6629F","EventStatus":"Scheduled","EventType":"Reboot","ResourceType":"VirtualMachine","Resources":["VM"],"NotBefore":"Thu, 16 Jul 2020 15:48:46 GMT"}]}
```

# プラグイン

plugin ディレクトリにプラグインがあります。

# Kubernetes 向けのサンプル

k8s-sample にDaemonSet で起動し、cron で本ツールを動かすサンプルがあります。