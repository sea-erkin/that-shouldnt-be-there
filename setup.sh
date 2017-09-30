. ./config.cfm

CONFIG_OS=$(gawk -F= '/^NAME/{print $2}' /etc/os-release)
START_DIR=$(pwd)
TARGET_GO_DIR=$HOME/go/src/github.com/sea-erkin
TARGET_PROJECT_DIR=$TARGET_GO_DIR/that-shouldnt-be-there

export GOROOT=$HOME/go
export GOPATH=$HOME/go  
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

mkdir -p $TARGET_PROJECT_DIR

if [[ $CONFIG_OS = '"Ubuntu"' ]]; then

  apt install golang-go	

  git clone https://github.com/sea-erkin/that-shouldnt-be-there.git
  
  mv that-shouldnt-be-there/ $TARGET_GO_DIR

  cd $TARGET_GO_DIR/that-shouldnt-be-there && go get && go build

  rm -rf $START_DIR

fi

if [[ $CONFIG_STORE_RESULTS = true ]]; then

  apt install sqlite3
  cd $TARGET_PROJECT_DIR && sqlite3 ./state/tsbt.db -init ./state/initDb.sql

fi

if [[ $CONFIG_SUBDOMAIN_SUBLISTER = true ]]; then

  mkdir $TARGET_PROJECT_DIR/other-tools
  cd $TARGET_PROJECT_DIR/other-tools
 
  git clone https://github.com/aboul3la/Sublist3r.git

fi
