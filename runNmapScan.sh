# bin/sh

# Run Nmap
sudo nmap -sS -Pn -A -v -iL ./state/hosts_nmap.txt -oX output_nmap.xml

# Insert Records into Sqlite
./that-shouldnt-be-there -config=./state/config.json -parseNmap

# Alert logic
./that-shouldnt-be-there -config=./state/config.json -alertPort

# Archive output file
now=$(date +"%m-%d-%Y")
fileName="output_nmap_"$now".xml"
mv output_nmap.xml $fileName
mv $fileName ./archive/

