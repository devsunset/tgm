#!/bin/bash
sudo kill -9 $(ps aux | grep 'tgm' | awk '{print $2}')
echo "tgm service stop ..."


