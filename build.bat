@echo off
set year=%date:~,4%
set month=%date:~5,2%
set day=%date:~8,2%
set buildTime=%year%-%month%-%day% %time%
set /p version=<version.txt
go build -trimpath -ldflags "-w -s -X 'main.version=%version%' -X 'main.buildTime=%buildTime%'"