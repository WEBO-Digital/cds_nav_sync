@echo off
rem Get current working batch file directory 
rem Then remove "crons" from the path
set "batchFilePath=%~dp0"
set "sourceFolderPath=%batchFilePath:\crons=%"

rem print the paths
echo Original Path: %batchFilePath%
echo Modified Path: %sourceFolderPath%

rem Move console to the source directory 
rem And Then execute the command
set "exePath=%sourceFolderPath%"
cd "%exePath%"
nav_sync_test.exe -action invoice_sync
cd %SystemDrive%
