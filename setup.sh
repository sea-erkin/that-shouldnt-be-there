mkdir other-tools

. ./config.cfm

CONFIG_OS=$(gawk -F= '/^NAME/{print $2}' /etc/os-release)

if [[ $CONFIG_SUBDOMAIN_SUBLISTER = true ]]; then

  git clone https://github.com/aboul3la/Sublist3r.git
  mv Sublist3r /other-tools

fi

echo $CONFIG_OS

if [[ $CONFIG_OS = '"Ubuntu"' ]]; then

  apt install golang-go	
  go get
  go build

fi