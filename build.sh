#!/bin/sh

# アプリケーションのビルド
go build -tags netgo -ldflags '-s -w' -o app

# テンプレートディレクトリをコピー
cp -r templates/ . 