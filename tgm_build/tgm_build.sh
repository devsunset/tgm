!/bin/bash
echo '-------------------------------------'
echo '-------------------------------------'
echo '_____  _____ __  ___  '
echo '|_   _||  __ \|  \/  |'
echo '| |  | |  \/| .  . |  '
echo '  | |  | | __ | |\/| |'
echo '  | |  | |_\ \| |  | |' _  _  _
echo '  \_/   \____/\_|  |_/ (_)(_)(_)'
echo '-------------------------------------'
echo '-------------------------------------'
echo 'Welocme to TGM Build System'
echo '-------------------------------------'
echo '*** Pre-installation environment ***'
echo 'go & nodejs & npm & git install'
echo '-------------------------------------'
sleep 1s

pwd 
pwd=$(pwd)


# 0.TGM_BUILD FOLDER CHECK
if [ ! -d ./build_tgm ] ; then
    echo 'tgm_build directory create ...'
    mkdir ./build_tgm
else
    echo 'tgm_deploy directory clear ...'
    rm -rf ./build_tgm/*
fi 


# 1.GET GIT TGM PROJECT
if [ ! -d ./tgm ] ; then
    echo 'git clone ...'
    git clone https://github.com/devsunset/tgm.git
else
    echo 'git pull ...' 
    git pull
fi

cd "$pwd/tgm/frontend"

if [ ! -d ./dist ] ; then
    mkdir dist
fi 

# 2.NPM INSTLL
# 3.NPM RUN BUILD
# 4.GO BUILD
# 5.WEBSSH2 COPY
# 6.WEBSSH2 NPM INSTALL
npm install && npm run build && cd "$pwd/tgm" && go build && cp  -R "$pwd/tgm/webssh2" "$pwd/build_tgm/" && cd "$pwd/build_tgm/webssh2" && npm install 

# 7.BUILD FILE PATCH
cd "$pwd/build_tgm" && cp "$pwd/tgm/tgm_build/start_tgm.sh" start_tgm.sh && cp "$pwd/tgm/tgm_build/stop_tgm.sh" stop_tgm.sh && mv "$pwd/tgm/tgm" tgm
date
date +"%FORMAT"
var=$(date)
echo "BUILD_DATE : $var" > "version.txt"

ls -altr
