#!/bin/bash
# Import config
. ./config.cfg

# Global Variables
userid=$(id -u)
osinfo=$(cat /etc/issue|cut -d" " -f1|head -n1)
unameOut="$(uname -s)"
startDir=$(pwd)
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac
echo ${machine}

# Go source installation variables
TARGET_GO_DIR=$HOME/go/src/github.com/sea-erkin
TARGET_PROJECT_DIR=$TARGET_GO_DIR/that-shouldnt-be-there

# Make go directories
mkdir -p $HOME/goroot
mkdir -p $TARGET_GO_DIR

clear

echo '#######################################################################'
echo '#                       That-Shouldnt-Be-There                        #'
echo '#######################################################################'
echo

if [ "${userid}" != '0' ]; then
  echo '[Error]: You must run this setup script with root privileges.'
  echo
  exit 1
fi

echo $osinfo

case ${machine} in
  Mac)

    # install golang
    echo
    echo "[*] Installing golang to build source"
    #wget --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36" https://storage.googleapis.com/golang/go1.9.1.darwin-amd64.tar.gz
    tar -C /usr/local -xzf go1.9.1.darwin-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bash_profile
    rm go1.9.1.darwin-amd64.tar.gz

    cd .. && cp -r that-shouldnt-be-there/ $TARGET_PROJECT_DIR

    export GOROOT=$HOME/goroot
    export GOPATH=$HOME/go  
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

    cd $TARGET_PROJECT_DIR && go get && go build

    # install phantomjs
    echo
    echo "[*] Installing phantomjs binaries"
    wget --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36" https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-macosx.zip
    unzip phantomjs-2.1.1-macosx.zip
    mv phantomjs-2.1.1-macosx/bin/phantomjs ./other-tools/phantomjs/
    rm -rf phantomjs-2.1.1-macosx
    rm phantomjs-2.1.1-macosx.zip
    echo "[*] Installed phantomjs"

    # install sqlite3
    echo
    echo "[*] Installing sqlite3 binaries"
    wget --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36" https://www.sqlite.org/2017/sqlite-tools-osx-x86-3200100.zip
    unzip sqlite-tools-osx-x86-3200100.zip
    mkdir ./other-tools/sqlite3
    mv sqlite-tools-osx-x86-3200100/sqlite3 ./other-tools/sqlite3
    rm -rf sqlite-tools-osx-x86-3200100/
    rm sqlite-tools-osx-x86-3200100.zip
    cd $TARGET_PROJECT_DIR
    echo "[*] Installed sqlite3"

  ;;
  Linux)
    # Get the flavor of linux and install golang respectively of each distro. Arch, ubuntu/kali for now
    apt install golang-go	

    cd .. && cp that-shouldnt-be-there/* $TARGET_PROJECT_DIR

    export GOROOT=$HOME/goroot
    export GOPATH=$HOME/go  
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

    cd $TARGET_PROJECT_DIR && go get && go build
    ;;
  *)
  echo 'Machine not detected'
esac

# None machine specific tools
# install sublister
echo
echo "[*] Installing Sublist3r for passive subdomain identification"
git clone https://github.com/aboul3la/Sublist3r.git
mv Sublist3r ./other-tools
cd ./other-tools/Sublist3r
sudo pip install -r requirements.txt
cd $TARGET_PROJECT_DIR
echo "[*] Installed Sublist3r"

if [[ $CONFIG_SUBDOMAIN_ALTDNS = true ]]; then
  # install altdns
  echo
  echo "[*] Installing altdns for subdomain permutation"
  git clone https://github.com/infosec-au/altdns.git
  mv altdns ./other-tools
  cd ./other-tools/altdns
  sudo pip install -r requirements.txt
  echo "[*] Installed altdns"
  cd $TARGET_PROJECT_DIR
fi

# initialize sqlite database schema
echo
echo "[*] Initializing sqlite3 database to store results"
./other-tools/sqlite3/sqlite3 tsbt.db tsbt.db -init ./state/initDb.sql
mv tsbt.db ./state
echo "[*] Initialized sqlite3 database"

# if [ ${MACHINE_TYPE} == 'x86_64' ]; then
#     wget --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36" https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-x86_64.tar.bz2
#     tar -xvf phantomjs-2.1.1-linux-x86_64.tar.bz2
#     cd phantomjs-2.1.1-linux-x86_64/bin/
#     mv phantomjs ../../
#     cd ../..
#     rm -rf phantomjs-2.1.1-linux-x86_64
#     rm phantomjs-2.1.1-linux-x86_64.tar.bz2
# else
#     wget --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36" https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-i686.tar.bz2
#     tar -xvf phantomjs-2.1.1-linux-i686.tar.bz2
#     cd phantomjs-2.1.1-linux-i686/bin/
#     mv phantomjs ../../
#     cd ../..
#     rm -rf phantomjs-2.1.1-linux-i686
#     rm phantomjs-2.1.1-linux-i686.tar.bz2
# fi
  
# OS Specific Installation Statement
# case ${osinfo} in
#   # Kali 2 dependency Install
#   Kali2)   
#   ;;
#   # Kali Dependency Installation
#   Kali)
#   ;;
#   # Debian 7+ Dependency Installation
#   Debian)
#   ;;
#   # Ubuntu Dependency Installation
#   Ubuntu)
#   ;;
#   # Notify Manual Installation Requirement And Exit
#   *)
#     echo "[Error]: ${osinfo} is not supported by this setup script."
#     echo
#     exit 1
# esac

# Finish Message
echo '[*] Setup script completed successfully.'
echo
