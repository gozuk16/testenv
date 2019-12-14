# gret
Resource and Environment testing tool implemented with Go.

Go実装のリソース、環境テストツール

主にWindowsでファイルの構成と環境変数のテストをしたいけど、Serverspecを使うためにRuby入れたりしたくないし、GossはWindows対応してないっぽいので簡単なものを作り始めました。（いつ出来るものやら・・・）

## Memo
- spf13/cobra のお勉強。CLIツールの側を作成
- path/filepath でMacで開発しつつWindowsで使えるように
- 指定されたディレクトリ以下を走査してファイル情報を取得
- ファイルのハッシュ値を取得

## ToDo
- YAMLにファイル情報を書き出す

