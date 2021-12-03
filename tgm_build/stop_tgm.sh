#!/bin/bash
sudo kill -9 $(ps aux | grep 'webssh2' | awk '{print $2}')
echo "tgm service stop ..."

sudo kill -9 $(ps aux | grep 'tgm' | awk '{print $2}')
echo "tgm service stop ..."



