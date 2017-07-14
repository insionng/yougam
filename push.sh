echo "comment:";
DATE=$(date +%Y%m%d)

find ./vendor -name .git | xargs rm -fr
find ./libraries -name .git | xargs rm -fr
find ./libraries -name .lnk | xargs rm -fr
find ./modules -name .git | xargs rm -fr
find ./themes -name .git | xargs rm -fr

rm -rf ./yougam
rm -rf ./bin/*.*
rm -rf ./logs/*.*
rm -rf ./public

if [ -e "data" ]; then
	rm -rf ../you-$DATE-data
	mkdir ../you-$DATE-data
	mv ./data ../you-$DATE-data/
	if [ -e "file" ]; then
		mv ./file ../you-$DATE-data/
	fi
fi

git add . -A

echo "\"$@\""|xargs git commit -am
git push

