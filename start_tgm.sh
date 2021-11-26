#!/bin/bash
 
IP=$(hostname -I)
#echo $IP

OLD_IFS="$IFS"
IFS=" "
IP_ARRAY=( $IP )
IFS="$OLD_IFS"

HOST_IP="" 

for x in "${IP_ARRAY[@]}"
do
    echo "HOST IP : $x"
    HOST_IP=$x
    break
done

#echo $HOST_IP
sudo ./tgm -a $HOST_IP -l tgm.log  --disable-cmd-limit=true &
#sudo ./tgm -a $HOST_IP -l tgm.log  --disable-cmd-limit=false &
echo "tgm service start success ..."
