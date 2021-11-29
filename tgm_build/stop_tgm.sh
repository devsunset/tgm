#!/bin/bash
echo "tgm service stop ..."
sudo kill -9 $(ps aux | grep 'tgm' | awk '{print $2}')


