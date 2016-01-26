@echo on

MKDIR c:\FFMPEG\
MKDIR c:\FFMPEG\bin

XCOPY windows_binaries\windows_ffmpeg\bin c:\FFMPEG\bin /y /s
SETX PATH "%PATH%;c:\FFMPEG\bin" 