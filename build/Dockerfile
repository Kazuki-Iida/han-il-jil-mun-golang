# goバージョン
FROM golang:1.19.1-alpine
# アップデートとgitのインストール
RUN apk update && apk add git
# boiler-plateディレクトリの作成(8行目のコマンドで、ディレクトリが存在しない場合に自動で作ってくれるはずなので6行目はコメントアウト)
# RUN mkdir -p /go/src/github.com/boiler-plate
# ワーキングディレクトリの設定
WORKDIR /go/src/github.com/boiler-plate
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/github.com/boiler-plate
# パッケージのインポート
# RUN go get -u golang.org/x/tools/cmd/goimports
RUN go install golang.org/x/tools/cmd/goimports@latest
#golang-migrateのインストール
RUN go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#あとでコメントアウト外しても問題ないかテスト
# RUN go get "github.com/go-sql-driver/mysql"

RUN apk update && apk add git