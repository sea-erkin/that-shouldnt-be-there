# /bin/sh

currentDir=$(pwd)
toolDir="/other-tools/phantomjs/"
targetDir="/state/eyewitness/todo"
doneDir="/state/eyewitness/done"
reportDir="/state/eyewitness/report"

for i in $(ls $currentDir/$targetDir); do

  for j in $(cat $currentDir/$targetDir/$i); do

    cd $currentDir/$toolDir && phantomjs screenshot.js "http://"$j
    
    cd $currentDir/$toolDir && phantomjs screenshot.js "https://"$j

  done

done

mv $currentDir/$toolDir/success* $currentDir/$doneDir

now=$(date +"%m-%d-%Y")
fileName="screenshot_"$now".tar.gz"

cd $currentDir/$doneDir && tar -czvf $fileName success*

rm $currentDir/$doneDir/success*