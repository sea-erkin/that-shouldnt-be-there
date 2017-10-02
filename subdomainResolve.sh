# /bin/sh

# Resolves subdomains found from recon to ip addresses
# Puts IPs in nmap todo and hosts/ips in eyewitness todo

pwd=$(pwd)
targetDir="./state/subdomains-resolve/todo"
doneDir="./state/subdomains-resolve/done"
nmapDir="./state/nmap/todo"
eyeWitnessDir="./state/eyewitness/todo"

# Resolve hostnames to ips
for i in $(ls $targetDir); do

  for j in $(cat $targetDir/$i); do
    
    nslookup $j | grep Address | grep -v 192.168.1.1 | cut -d " " -f 2 | awk '{print $1"_""'"$j"'"}' >> $doneDir/$i.tmp
 
  done 

  cat $doneDir/$i.tmp | sort | uniq > $doneDir/$i

done

rm $doneDir/*.tmp
rm $targetDir/*

# Prep resolved ips for nmap and eyewitness
for i in $(ls $doneDir); do 

  for j in $(cat $doneDir/$i | sort | uniq); do

    # prep ips for nmap output
    echo $j | cut -d "_" -f 1 >> $nmapDir/$i.tmp
  
    # prep ips, hosts, for eyewitness
    echo $j | tr _ '\n' >> $eyeWitnessDir/$i.tmp

  done

  cat $nmapDir/$i.tmp | sort | uniq > $nmapDir/$i
  cat $eyeWitnessDir/$i.tmp | sort | uniq > $eyeWitnessDir/$i

done

rm $nmapDir/*.tmp
rm $eyeWitnessDir/*.tmp