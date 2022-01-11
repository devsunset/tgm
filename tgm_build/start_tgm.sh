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
sudo ./tgm -a $HOST_IP -p 8282 -w 8383 -l tgm.log &
echo "tgm service start success ..."

cd ./webssh2 && npm start > ../webssh2.log &
echo "webssh2 service start success ..."
