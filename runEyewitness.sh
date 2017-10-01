# /bin/sh

pwd=$(pwd)
targetDir="./state/eyewitness/todo"
doneDir="./state/eyewitness/done"
reportDir="./state/eyewitness/report"

for i in $(ls $targetDir); do

  ./other-tools/EyeWitness/EyeWitness.py -f $targetDir/$i --web 

  #mv $targetDir/$i $doneDir

done
