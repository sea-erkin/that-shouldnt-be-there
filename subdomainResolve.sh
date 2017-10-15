# bin/bash

# Resolves subdomains found from recon to ip addresses
# Puts IPs in nmap todo and hosts/ips in screenshot todo

pwd=$(pwd)
targetDir="./state/subdomains-resolve/todo"
doneDir="./state/subdomains-resolve/done"
archiveDir="./state/subdomains-resolve/archive"
nmapDir="./state/nmap/todo"
screenshotDir="./state/screenshot/todo"

# Resolve hostnames to ips
for i in $(ls $targetDir); do

  for j in $(cat $targetDir/$i); do
    
    nslookup $j | grep Address | grep -v 192.168.1.1 | grep -v '#53' | cut -d " " -f 2 | awk '{print $1"_""'"$j"'"}' >> $doneDir/$i.tmp
 
  done 

  cat $doneDir/$i.tmp | sort | uniq > $doneDir/$i

done

rm $doneDir/*.tmp
rm $targetDir/*

# Prep resolved ips for nmap and screenshot
for i in $(ls $doneDir); do 

  for j in $(cat $doneDir/$i | sort | uniq); do

    # prep ips for nmap output
    echo $j | cut -d "_" -f 1 >> $nmapDir/$i.tmp
  
    # prep ips, hosts, for screenshot
    echo $j | tr _ '\n' >> $screenshotDir/$i.tmp

  done

  cat $nmapDir/$i.tmp | sort | uniq > $nmapDir/$i
  cat $screenshotDir/$i.tmp | sort | uniq > $screenshotDir/$i

done

mv $doneDir/* $archiveDir
rm $nmapDir/*.tmp
rm $screenshotDir/*.tmp

cd $pwd && ./runNmapScan.sh
