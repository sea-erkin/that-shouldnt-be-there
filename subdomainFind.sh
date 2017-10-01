# bin/sh

startDirectory=$(pwd)

for domain in $(cat domains.txt); do
    
    sleep 0.5

    cd $startDirectory

    echo "Starting in "$startDirectory" for domain: "$domain
    
    fileName=$domain"_sublister_$(date +%s)"

    # Run Sublister
    cd ./other-tools/Sublist3r/ && ./sublist3r.py -d $domain -o $fileName
    
    # Move file to subdomain todo directory
    mv $fileName ../../state/subdomains/todo/

    sleep 0.5

done;

# Run Next Script which parses and inserts the records into database.

echo "Starting parse and alert logic"

cd $startDirectory && ./subdomainParse.sh
