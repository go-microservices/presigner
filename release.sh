#!/bin/sh

go get -d github.com/laher/goxc
goxc -wlc default publish-github -apikey=$GITHUB_TOKEN
# goxc bump
goxc -bc="linux,windows,darwin"
# git config --global user.email "release@circleci.com"
# git config --global user.name "CircleCI"
# git commit ".goxc.json" -m "Bump version [ci skip]"
# git push https://${GITHUB_TOKEN}@github.com/minodisk/presigner.git release
