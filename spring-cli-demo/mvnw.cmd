@echo off
where mvn >nul 2>&1 || (echo Maven not found. Please install Maven. & exit /b 1)
mvn %*
