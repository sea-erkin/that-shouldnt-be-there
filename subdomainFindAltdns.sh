# bin/sh

startDirectory=$(pwd)
moreSubdomainsTodoDirectory=$startDirectory'/state/subdomains-more/todo'
moreSubdomainsDoneDirectory=$startDirectory'/state/subdomains-more/done'
targetSubdomainTodoDirectory=$startDirectory'/state/subdomains/todo'

for subdomainFile in $(ls $moreSubdomainsTodoDirectory); do
    
    domainName=$(echo $subdomainFile| cut -f 1 -d "_")
    fileName=$domainName"_altdns_$(date +%s)"

    # Run altdns
    cd $startDirectory/other-tools/altdns/ && ./altdns.py -i $moreSubdomainsTodoDirectory/$subdomainFile -o data_output -w words.txt -r -s $fileName.tmp
    # ./altdns.py -i subdomains.txt -o data_output -w words.txt -r -s results_output.txt

    # Resultant file will look like hostname:ip
    # Echo split the file and put the hostname results in subdomain todo
    cat $fileName.tmp | cut -f 1 -d ":" > $targetSubdomainTodoDirectory/$fileName

    # Remove left over file
    rm $fileName.tmp

    # Move current file to completed
    mv $moreSubdomainsTodoDirectory/$subdomainFile $moreSubdomainsDoneDirectory

done;

cd $startDirectory && ./subdomainParse.sh
