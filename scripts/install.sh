#!/bin/bash -e
#
# Usage:
#   $ curl -fsSL https://raw.githubusercontent.com/denizgursoy/gotouch/main/scripts/install.sh | bash
# or
#   $ wget -q https://raw.githubusercontent.com/denizgursoy/gotouch/main/scripts/install.sh -O- | bash
#

case $(uname -sm) in
'Linux x86_64') ;;
'Linux arm64') ;;
'Linux i386') ;;
'Darwin x86_64') ;;
'Darwin arm64') ;;

*)
  echo "Sorry, you'll need to install gotouch manually."
  exit 1
  ;;
esac

os=$(uname -s)
arch=$(uname -m)
tag=$(basename $(curl -fs -o/dev/null -w %{redirect_url} https://github.com/denizgursoy/gotouch/releases/latest))
filename="gotouch_${tag#v}_${os}_${arch}.tar.gz"

curl -LO https://github.com/denizgursoy/gotouch/releases/download/${tag}/${filename}

tar -xzf ${filename} --overwrite -C /usr/local/bin/ --wildcards --no-anchored 'gotouch*'

rm ${filename}
