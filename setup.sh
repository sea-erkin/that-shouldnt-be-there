. ./config.cfm

CONFIG_OS=$(gawk -F= '/^NAME/{print $2}' /etc/os-release)
START_DIR=$(pwd)
TARGET_GO_DIR=$HOME/go/src/github.com/sea-erkin/
TARGET_PROJECT_DIR=$TARGET_GO_DIR/that-shouldnt-be-there/

mkdir -p $HOME/goroot
mkdir -p $TARGET_PROJECT_DIR

if [[ $CONFIG_OS = '"Ubuntu"' ]]; then

  apt install golang-go	

  cd .. && mv that-shouldnt-be-there/* $TARGET_PROJECT_DIR

  export GOROOT=$HOME/goroot
  export GOPATH=$HOME/go  
  export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

  cd $TARGET_PROJECT_DIR && go get && go build

 # rm -rf $START_DIR

fi

if [[ $CONFIG_STORE_RESULTS = true ]]; then

  apt install sqlite3
  cd $TARGET_PROJECT_DIR && sqlite3 ./state/tsbt.db -init ./state/initDb.sql

fi

if [[ $CONFIG_SUBDOMAIN_SUBLISTER = true ]]; then

  apt install python-pip

  mkdir $TARGET_PROJECT_DIR/other-tools/
  cd $TARGET_PROJECT_DIR/other-tools/
 
  git clone https://github.com/aboul3la/Sublist3r.git
  
  cd Sublist3r && pip install -r requirements.txt

fi
