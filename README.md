
# TGM

## Install

git clone  https://github.com/devsunset/tgm.git

cd tgm/tgm_build

./tgm_build.sh

start_tgm.sh

http://localhost:8282

## Reference Open Source
filebrowser - https://github.com/filebrowser/filebrowser


## DB (boltdb - golang storm) - Tool

### cli version - boltbrowser
https://github.com/br0xen/boltbrowser

    - install
    go get github.com/br0xen/boltbrowser

    - execute 
    $GOPATH/bin/boltbrowser tgm.db

### web version - boltdbweb
https://github.com/evnix/boltdbweb

    - install
    go get github.com/evnix/boltdbweb

    - execute 
    $GOPATH/bin/boltbrowser --db-name=tgm.db
    Goto: http://localhost:8089

    Usage
    boltdbweb --db-name=<DBfilename>[required] --port=<port>[optional] --static-path=<static-path>[optional]
        --db-name: The file name of the DB.
            NOTE: If 'file.db' does not exist. it will be created as a BoltDB file.
        --port: Port for listening on... (Default: 8080)
        --static-path: If you moved the binary to different folder you can determin the path of the web folder. (Default: Same folder where the binary is located.)
