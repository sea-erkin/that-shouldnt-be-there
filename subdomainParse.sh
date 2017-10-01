# bin/sh

startDirectory=`pwd`

# Run Parse Logic
cd $startDirectory && ./that-shouldnt-be-there -c=./state/config.json -parseSubdomain

# Run Alert Logic
cd $startDirectory && ./that-shouldnt-be-there -c=./state/config.json -alertSubdomain