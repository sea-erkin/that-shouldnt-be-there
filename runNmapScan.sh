# bin/sh

currentDir=$(pwd)
todoDir="./state/nmap/todo/"
doneDir="./state/nmap/done/"

for i in $(ls $todoDir); do

    # only web
    cd $currentDir/$todoDir && sudo nmap -vvv -sS -Pn -T 2 -p 80,443,8080,8000, 8001,8443 -iL $i -oX $i.xml

    # Insert Records into Sqlite
    cd $currentDir && ./that-shouldnt-be-there -c=./state/config.json -parseNmap=$currentDir/$todoDir/$i.xml -d

    # Alert logic
    cd $currentDir && ./that-shouldnt-be-there -c=./state/config.json -alertPort=$i -d

    sleep 0.5

    # Archive output file
    mv $currentDir/$todoDir/$i $currentDir/$doneDir/$i
    mv $currentDir/$todoDir/$i.xml $currentDir/$doneDir/$i.xml

done;

cd $currentDir && ./runScreenshot.sh
