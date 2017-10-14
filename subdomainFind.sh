# bin/bash

# Import config
. ./config.cfg

startDirectory=$(pwd)

for domain in $(cat domains.txt); do

    cd $startDirectory

    echo "Starting in "$startDirectory" for domain: "$domain
    
    fileName=$domain"_sublister_$(date +%s)"

    # Run Sublister
    cd ./other-tools/Sublist3r/ && ./sublist3r.py -d $domain -o $fileName
    
    # Move file to the subdomain-more directory to get more subdomains based off the identified using altdns
    cp $fileName ../../state/subdomains-more/todo/

    # Move file to subdomain todo directory
    mv $fileName ../../state/subdomains/todo/

done;

if [[ $CONFIG_SUBDOMAIN_ALTDNS = true ]]; then
    echo "Running alt-dns to identify more subdomains"
    cd $startDirectory && ./subdomainFindAltdns.sh
else
    echo "Saving subdomain results"
    cd $startDirectory && ./subdomainParse.sh
fi


