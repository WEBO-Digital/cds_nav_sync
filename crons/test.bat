@echo off
set "originalPathToCron=%~dp0"
rem Remove "crons" from the path
set "sourceFolderPath=%originalPathToCron:\crons=%"

echo Original Path: %originalPathToCron%
echo Modified Path: %sourceFolderPath%