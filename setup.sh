mkdir other-tools

. ./config.cfm

CONFIG_OS=$(gawk -F= '/^NAME/{print $2}' /etc/os-release)
START_DIR=$(pwd)
TARGET_GO_DIR=$HOME/go/src/github.com/sea-erkin/

if [[ $CONFIG_SUBDOMAIN_SUBLISTER = true ]]; then

  git clone https://github.com/aboul3la/Sublist3r.git
  mv Sublist3r /other-tools

fi

echo $CONFIG_OS

if [[ $CONFIG_OS = '"Ubuntu"' ]]; then

  mkdir -p $TARGET_GO_DIR

  export GOROOT=$HOME/go
  export GOPATH=$HOME/go  
  export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

  apt install golang-go	

  git clone https://github.com/sea-erkin/that-shouldnt-be-there.git
  
  mv that-shouldnt-be-there/ $TARGET_GO_DIR

  cd $TARGET_GO_DIR/that-shouldnt-be-there && go get && go build

  rm -rf $START_DIR

fi
