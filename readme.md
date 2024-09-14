# GTFS-RT 可視化ツール

このツールは、複数のGTFS-RTフィードから位置情報データを抽出し、[kepler.gl](https://kepler.gl/demo) で可視化するための GeoJSON および Parquet ファイルを生成します。

[![](https://img.youtube.com/vi/LccELuPGdV8/0.jpg)](https://www.youtube.com/watch?v=LccELuPGdV8)
※クリックすると[Youtubeで動画再生](https://youtu.be/LccELuPGdV8)されます。

※上図データは[東京公共交通オープンデータセンター](https://www.odpt.org/)にて公開されている横浜市営バスGTFS-RTデータを使用

## 特徴

* 複数のGTFS-RTフィードをまとめて処理
* kepler.gl で可視化しやすい GeoJSON 形式と、分析に適した Parquet 形式の出力に対応
* 車両位置情報、混雑状況、タイムスタンプを出力
* ZIP 圧縮された GTFS-RT ファイルの自動展開

## インストール

### Windows の場合

1. リポジトリをクローンまたはダウンロードします。
2. リポジトリ内の `main.exe` を実行します。

### Golang が使用できる環境の場合

1. リポジトリをクローンします。
2. ターミナルで以下のコマンドを実行します。

```
go run main.go
```

## 使い方

1. **GTFS-RT ファイルの準備**: 処理したい GTFS-RT ファイルを `GTFS-RTs` ディレクトリに配置します。拡張子は `.gitignore` 以外であれば何でも構いません。
2. **ツールの実行**: 上記の「インストール」セクションの手順に従ってツールを実行します。
3. **出力ファイルの確認**: ツールの実行が完了すると、`GTFS-RT.json` (GeoJSON) と `GTFS-RT.parquet` (Parquet) が生成されます。これらのファイルを kepler.gl で読み込んで可視化することができます。

## ZIP 圧縮された GTFS-RT ファイルの処理

複数の GTFS-RT ファイルが ZIP 圧縮されている場合は、以下の手順で処理します。

1. **ZIP ファイルの配置**: ZIP 圧縮された GTFS-RT ファイルを `zip` フォルダに配置します。
2. **ZIP ファイルの展開**: `unzip.exe` (Windows の場合) または `go run unzip.go` (Golang の場合) を実行して、ZIP ファイルを展開します。展開されたファイルは `GTFS-RTs` ディレクトリに格納されます。
3. **ツールの実行**: 上記の「使い方」セクションの手順に従ってツールを実行します。

## ファイルツリーとプロジェクト構成

```
├── go.mod
├── go.sum
├── gtfs.py
├── main.go
├── matching.go
├── parquet-test.py
├── pb
│   ├── gtfs-realtime.pb.go
│   ├── gtfs-realtime.proto
│   └── readme.md
├── pkg
│   └── pkg.go
├── readme.md
└── unzip.go

```

### 各ファイル・ディレクトリの説明

* **`go.mod`, `go.sum`**: Go の依存関係管理ファイル
* **`gtfs.py`**: GTFS 静的データの読み込みと Parquet ファイルからのデータの抽出を行う Python スクリプト
* **`main.go`**: GTFS-RT ファイルを読み込み、GeoJSON および Parquet ファイルを生成する Go プログラム
* **`matching.go`**: GTFS 静的データと GTFS-RT データをマッチングする Go プログラム
* **`parquet-test.py`**: Parquet ファイルを読み込み、CSV ファイルに変換する Python スクリプト
* **`pb`**: GTFS-RT の Protocol Buffers 定義ファイルと、Go で使用するためのコードを含むディレクトリ
    * **`gtfs-realtime.pb.go`**: `gtfs-realtime.proto` から生成された Go コード
    * **`gtfs-realtime.proto`**: GTFS-RT の Protocol Buffers 定義ファイル
    * **`readme.md`**: `gtfs-realtime.pb.go` を生成するためのコマンド
* **`pkg`**: 共通関数を含むパッケージ
    * **`pkg.go`**: ファイル操作、データ構造定義、距離計算などの関数
* **`readme.md`**: この README ファイル
* **`unzip.go`**: `zip` フォルダ内の ZIP ファイルを展開する Go プログラム

## コマンド実行例

### ZIP ファイルの展開

```
go run unzip.go
```

### GTFS-RT ファイルの処理

```
go run main.go
```

## ライセンス

このプロジェクトは MIT ライセンスで公開されています。 
