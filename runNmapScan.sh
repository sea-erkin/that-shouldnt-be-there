# bin/sh

currentDir=$(pwd)
todoDir="/state/nmap/todo"
doneDir="/state/nmap/done"

for i in $(ls $todoDir); do

    # only web
    cd $currentDir/$todoDir && sudo nmap -sS -T 0 -p 80,443,8080,8000,8443 -iL $i -oX output_nmap.xml

    # Insert Records into Sqlite
    cd $currentDir && ./that-shouldnt-be-there -config=./state/config.json -parseNmap

    # Alert logic
    cd $currentDir && ./that-shouldnt-be-there -config=./state/config.json -alertPort

    # Archive output file
    now=$(date +"%m-%d-%Y")
    fileName="output_nmap_"$now".xml"
    mv $currentDir/$todoDir/output_nmap.xml $currentDir/$todoDir/$fileName
    mv $currentDir/$todoDir/$fileName $currentDir/$doneDir/$fileName

done;
