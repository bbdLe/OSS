#!/bin/bash
ps -ef | grep dataServer | grep -v grep | awk '{print $2}' | xargs kill -s 9
