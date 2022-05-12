#!/bin/sh

cp ./menu-stop.json ./menu.json
$(./bin/libgen-opds serve; cp ./menu-start.json ./menu.json) &
