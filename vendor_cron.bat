@echo off
set "exePath=%SystemDrive%\Users\Suman Neupane\Live Projects\CDS\nav_sync"
cd "%exePath%"
nav_sync_test.exe -action vendor_fetch 
cd %SystemDrive%
