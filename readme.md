# GTFS-RT 可視化ツール

GTFS-RTを可視化する為、複数のGTFS-RTから[kepler.gl](https://kepler.gl/demo)で可視化する為のGeoJSON及びParquetを生成するコードです。

<iframe width="560" height="315" src="https://www.youtube.com/embed/LccELuPGdV8" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

※上図データは[東京公共交通オープンデータセンター](https://www.odpt.org/)にて公開されている横浜市営バスGTFS-RTデータを使用

## 実行方法

``gtfsrts``ディレクトリにGTFS-RTを複数配置します。
拡張子は``.gitignore``以外であれば何でも構いません。

### Windowsの場合
``main.exe``を実行

### Golangが使用できる環境の場合
Linuxなどでも可能です。

```
go run main.go
```

を実行します。

## 複数のGTFS-RTがZIP圧縮されている場合

複数のGTFS-RTをまとめてZIP圧縮したファイル（複数可）を``zip``フォルダに格納し、``unzip.exe``又は``go run unzip``を実行してください。


