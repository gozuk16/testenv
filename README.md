# gret
Resource and Environment testing tool implemented with Go.

Go実装のファイルテストツール

環境テストもできるようにしたい

主にWindowsでファイルの構成と環境変数のテストをしたいけど、Serverspecを使うためにRuby入れたりしたくないし、GossはWindows対応してないっぽいので簡単なものを作り始めました。（いつ出来るものやら・・・）

## Usage
サブコマンドが2つあります。

### add
```
$ gret add directory
```

起点となるディレクトリを指定します。  
相対パスでも絶対パスでもどちらでも構いません。  
結果として実行可能なファイル一覧をJSONとして表示します。  
現状ではリダイレクトしてファイルに保存してください。

デフォルトでドットファイルは除外されます。  
含めたい場合は `-a` をつけて実行します。

```
$ gret add -a directory
```
 
### test
```
$ gret test testfile
```

現状ではファイルについてテスト可能です。  
test対象を記述したjsonファイルを指定してください。

exceptに指定できるキーワードは、`match`, `exist`, `not_exist`です。  

|キーワード|説明|
|--|--|
|match|完全一致。sha1のハッシュ値と一致する場合はテスト成功|
|exist|fullpathのファイルが存在していれば成功|
|not\_exist|fullpathのファイルが存在していなければ成功。存在してはいけないファイルをチェックします|

`warningOverlay`, `warningMaxPath` は記述されているすべてのファイルに掛かります。  
`warningOverlay` は指定された拡張子について類似するファイルがないか走査します。  
デフォルトでは、`jar`, `dll`が指定されています。
testを実行すると同じディレクトリの中にfilenameの後ろにセパレーター("\_", "-", ".")+数字のファイルがないかチェックします。  
`-e` をつけてtestを実行するとbaseDir以下のディレクトリを再帰的に全てチェックします。

`warningOverlay` は記述されているファイルの絶対パスの文字数をチェックします。  
Windowsでは259文字以上のパスでアプリによって動作がおかしくなる可能性があります。  
UNCを考慮してデフォルトでは220文字に設定してあります。

```
    "warningOverlay": [
        "jar",
        "dll"
    ],
    "warningMaxPath": 220,
```

デフォルトでは成功したケースの結果表示は行いません。  
表示を行いたい場合は `-v` をつけて実行します。
```
$ gret add -v directory
```

## ToDo
- UnitTestをもっと書く
- 環境変数テストの実装
