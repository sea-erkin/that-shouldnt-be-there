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
   
    # install wget, nmap
    echo
    echo "[*] Installing wget, nmap using brew"
    brew install wget
    brew install nmap

    # install golang
    echo
    echo "[*] Installing golang to build source"
    brew install go

    cd .. && cp -r that-shouldnt-be-there/ $TARGET_PROJECT_DIR

    export GOROOT=$HOME/goroot
    export GOPATH=$HOME/go  
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

    cd $TARGET_PROJECT_DIR && go get && go build

    # install phantomjs
    echo
    echo "[*] Installing phantomjs binaries"
    brew install phantomjs    
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
    
    # OS Specific Installation Statement
    case ${osinfo} in
    # Kali 2 dependency Install
    Kali2)   
    ;;
    # Kali Dependency Installation
    Kali)
    ;;
    # Debian 7+ Dependency Installation
    Debian)
    ;;
    # Ubuntu Dependency Installation
    Ubuntu)
   
      sudo apt-get update
      sudo apt-get install build-essential chrpath libssl-dev libxft-dev
 
      apt install golang-go	
      cd .. && cp -r that-shouldnt-be-there/ $TARGET_PROJECT_DIR

      export GOROOT=$HOME/goroot
      export GOPATH=$HOME/go  
      export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

      cd $TARGET_PROJECT_DIR && go get && go build
      
      # install phantomjs
      sudo apt-get install libfreetype6 libfreetype6-dev
      sudo apt-get install libfontconfig1 libfontconfig1-dev

      export PHANTOM_JS="phantomjs-1.9.8-linux-x86_64"
      wget https://bitbucket.org/ariya/phantomjs/downloads/$PHANTOM_JS.tar.bz2
      sudo tar xvjf $PHANTOM_JS.tar.bz2
      
      sudo mv $PHANTOM_JS /usr/local/share
      sudo ln -sf /usr/local/share/$PHANTOM_JS/bin/phantomjs /usr/local/bin

      ln -s $(which phantomjs) ./other-tools/phantomjs/phantomjs    
      rm phantomjs-1.9.8-linux-x86_64.tar.bz2
      rm -rf phantomjs-1.9.8-linux-x86_64

      # install sqlite3
      echo
      echo "[*] Installing sqlite3 binaries"
      sudo apt-get install sqlite3 libsqlite3-dev    
      
      mkdir ./other-tools/sqlite3

      ln -s $(which sqlite3) ./other-tools/sqlite3/sqlite3    

      cd $TARGET_PROJECT_DIR
      echo "[*] Installed sqlite3"
    ;;
    # Notify Manual Installation Requirement And Exit
    *)
      echo "[Error]: ${osinfo} is not supported by this setup script."
      echo
      exit 1
    esac
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

cd $TARGET_PROJECT_DIR && go get && go build

echo '[*] Setup script completed successfully.'
echo
