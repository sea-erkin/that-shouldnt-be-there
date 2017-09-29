# bin/sh

startDirectory=`pwd`

for domain in `cat ./state/domains.txt`; do
    
    sleep 0.5

    cd $startDirectory

    echo "Starting in "$startDirectory" for domain: "$domain
    
    # Run Sublister
    cd ./other-tools/Sublist3r/ && ./sublist3r.py -d $domain -o output_sublister.txt

    # Run Parse Logic
    cd $startDirectory && ./that-shouldnt-be-there -config=./state/config.json -parseSublister=$domain

    # Run Alert Logic
    ./that-shouldnt-be-there -config=./state/config.json -alertSubdomain=$domain

    # Rename and archive file
    cd ./other-tools/Sublist3r/
    now=$(date +"%m-%d-%Y")
    fileName="output_sublister_"$domain"_$(date +%s).txt"
    mv output_sublister.txt $fileName
    mv $fileName ../../archive/

    sleep 0.5

done;
