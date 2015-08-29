
cd ~/go/src/git.eycia.me/eycia/msghub

git pull

if [ x"$*" != x ]
then

	git status

	git add *

	git commit -m "$*"

	git push

fi
