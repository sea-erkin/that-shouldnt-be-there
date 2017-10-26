#!/bin/bash
# Import config
. ./config.cfg

currentDir=$(pwd)
toolDir="/other-tools/phantomjs/"
targetDir="/state/screenshot/todo"
doneDir="/state/screenshot/done"
archiveDir="/state/screenshot/archive"

for i in $(ls $currentDir/$targetDir); do

  # remove duplicates
  cat $currentDir/$targetDir/$i | sort | uniq > $currentDir/$targetDir/$i.tmp
  cat $currentDir/$targetDir/$i.tmp > $currentDir/$targetDir/$i
  rm $currentDir/$targetDir/$i.tmp

  if [[ $CONFIG_SCREENSHOT_ALWAYS = false ]]; then
    echo "Filtering out hosts that have been already screenshotted"

	# get the current file name and keep domain and source
	currentFileDomainSource=$(echo $i | cut -d "_" -f 1,2)

	echo "Current domain source: ""$currentFileDomainSource"

	# check to see if there are previous files
	pastFilesExist=$(ls -l $currentDir/$archiveDir/ | grep $currentFileDomainSource | wc -l)

	if [[ $pastFilesExist != 0 ]]; then
	  echo "Past files exist: ""$pastFilesExist"


   	        # get the latest file that has the same domain and source in the archive directory for this domain & source
        	latestFileInArchive=$(ls -t $currentDir/$archiveDir/$currentFileDomainSource* | head -1)
    
    		echo "Latest file in archive: ""$latestFileInArchive"
    
    		# difference the two files
    		difference=$(diff $currentDir/$targetDir/$i $latestFileInArchive)
    
    		echo "Difference" "$difference"
    
    		# if files not different, remove the file from the todo directory
    		if [[ $difference = "" ]]; then
    
    			echo "Removing file in screenshot todo because not different from last run: "$currentDir/$targetDir/$i
    			rm $currentDir/$targetDir/$i
			continue    

    		fi

	fi

  fi

  for j in $(cat $currentDir/$targetDir/$i); do

    cd $currentDir/$toolDir && phantomjs screenshot.js "http://"$j
    
    cd $currentDir/$toolDir && phantomjs screenshot.js "https://"$j

  done

done

totalCount=$(ls $currentDir/$toolDir | grep success | wc -l)

echo "Total screenshot count: ""$totalCount"

if [[ $totalCount > 0 ]]; then

  echo "Sending "$totalCount" Screenshots"

  mv $currentDir/$toolDir/success* $currentDir/$doneDir

  now=$(date +"%m-%d-%Y")
  fileName="screenshot_"$now".tar.gz"
         
  cd $currentDir/$doneDir && tar -czvf $fileName success*
         
  rm $currentDir/$doneDir/success*
         
  mv $currentDir/$targetDir/* $currentDir/$archiveDir/
         
  cd $currentDir && ./that-shouldnt-be-there -sendScreenshot -c=./state/config.json -d

else

  echo "No new screenshots to send"

  exit 0

fi
