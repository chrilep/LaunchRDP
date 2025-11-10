Unicode true

####
## Please note: Template replacements don't work in this file. They are provided with default defines like
## mentioned underneath.
## If the keyword is not defined, "wails_tools.nsh" will populate them with the values from ProjectInfo.
## If they are defined here, "wails_tools.nsh" will not touch them. This allows to use this project.nsi manually
## from outside of Wails for debugging and development of the installer.
##
## For development first make a wails nsis build to populate the "wails_tools.nsh":
## > wails build --target windows/amd64 --nsis
## Then you can call makensis on this file with specifying the path to your binary:
## For a AMD64 only installer:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app.exe
## For a ARM64 only installer:
## > makensis -DARG_WAILS_ARM64_BINARY=..\..\bin\app.exe
## For a installer with both architectures:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app-amd64.exe -DARG_WAILS_ARM64_BINARY=..\..\bin\app-arm64.exe
####
## The following information is taken from the ProjectInfo file, but they can be overwritten here.
####
## !define INFO_PROJECTNAME    "MyProject" # Default "{{.Name}}"
## !define INFO_COMPANYNAME    "MyCompany" # Default "{{.Info.CompanyName}}"
## !define INFO_PRODUCTNAME    "MyProduct" # Default "{{.Info.ProductName}}"
## !define INFO_PRODUCTVERSION "1.0.0"     # Default "{{.Info.ProductVersion}}"
## !define INFO_COPYRIGHT      "Copyright" # Default "{{.Info.Copyright}}"
###
## !define PRODUCT_EXECUTABLE  "Application.exe"      # Default "${INFO_PROJECTNAME}.exe"
## !define UNINST_KEY_NAME     "UninstKeyInRegistry"  # Default "${INFO_COMPANYNAME}${INFO_PRODUCTNAME}"
####
## !define REQUEST_EXECUTION_LEVEL "admin"            # Default "admin"  see also https://nsis.sourceforge.io/Docs/Chapter4.html
####
## Include the wails tools
####
!include "wails_tools.nsh"

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

# Enable HiDPI support. https://nsis.sourceforge.io/Reference/ManifestDPIAware
ManifestDPIAware true

!include "MUI.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
# !define MUI_WELCOMEFINISHPAGE_BITMAP "resources\leftimage.bmp" #Include this to add a bitmap on the left side of the Welcome Page. Must be a size of 164x314
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.

# Customize directory page text
!define MUI_DIRECTORYPAGE_TEXT_TOP "Setup will install ${INFO_PRODUCTNAME} in the following folder.$\r$\n$\r$\nTo continue, click Next."

Page custom DirectoryAndOptionsPage DirectoryAndOptionsPageLeave # Install location and desktop shortcut
!insertmacro MUI_PAGE_INSTFILES # Installing page (includes WebView2 install if needed)

!insertmacro MUI_UNPAGE_CONFIRM # Uninstall confirmation page
!insertmacro MUI_UNPAGE_COMPONENTS # Uninstall components page (for user data option)
!insertmacro MUI_UNPAGE_INSTFILES # Uninstalling page

!insertmacro MUI_LANGUAGE "English" # Set the Language of the installer

## The following two statements can be used to sign the installer and the uninstaller. The path to the binaries are provided in %1
#!uninstfinalize 'signtool --file "%1"'
#!finalize 'signtool --file "%1"'

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PRODUCTNAME} ${INFO_PRODUCTVERSION} Installer.exe" # Name of the installer's file.
InstallDir "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}" # Default installing folder ($PROGRAMFILES is Program Files folder).
ShowInstDetails show # This will always show the installation details.

# Silent install support for SCCM/Intune deployment
SilentInstall normal
SilentUnInstall normal

Var CreateDesktopShortcut
Var DesktopCheckbox
Var DirRequest
Var SpaceRequiredLabel
Var SpaceAvailableLabel

Function .onInit
   !insertmacro wails.checkArchitecture
   StrCpy $CreateDesktopShortcut "0"  # Default: no desktop shortcut
   
   # Check if running in silent mode
   ${If} ${Silent}
       # Skip custom page in silent mode - use defaults
   ${EndIf}
FunctionEnd

# Custom page with directory selection and desktop shortcut checkbox
Function DirectoryAndOptionsPage
    # Skip custom page in silent mode
    ${If} ${Silent}
        Abort
    ${EndIf}
    
    nsDialogs::Create 1018
    Pop $0
    
    ${If} $0 == error
        Abort
    ${EndIf}
    
    ${NSD_CreateLabel} 0 0 100% 32u "Setup will install ${INFO_PRODUCTNAME} in the following folder.$\r$\n$\r$\nTo install in a different folder, click Browse and select another folder. Click Install to continue."
    Pop $0
    
    ${NSD_CreateLabel} 0 36u 100% 12u "Destination Folder"
    Pop $0
    
    ${NSD_CreateDirRequest} 0 50u 85% 12u "$INSTDIR"
    Pop $DirRequest
    
    ${NSD_CreateBrowseButton} 86% 50u 14% 12u "Browse..."
    Pop $1
    ${NSD_OnClick} $1 BrowseDirectory
    
    ${NSD_CreateLabel} 0 68u 50% 12u "Space required: 11 MB"
    Pop $SpaceRequiredLabel
    
    ${NSD_CreateLabel} 50% 68u 50% 12u "Space available: calculating..."
    Pop $SpaceAvailableLabel
    
    ${NSD_CreateCheckbox} 0 86u 100% 12u "Create &desktop shortcut"
    Pop $DesktopCheckbox
    ${NSD_SetState} $DesktopCheckbox ${BST_UNCHECKED}
    
    Call UpdateSpaceAvailable
    
    nsDialogs::Show
FunctionEnd

Function UpdateSpaceAvailable
    # Get the drive from INSTDIR
    StrCpy $0 $INSTDIR 3
    
    # Get available space in KB
    System::Call 'kernel32::GetDiskFreeSpaceEx(t "$0", *l .r1, *l, *l) i .r2'
    
    ${If} $2 != 0
        # Convert KB to MB
        System::Int64Op $1 / 1048576
        Pop $1
        ${NSD_SetText} $SpaceAvailableLabel "Space available: $1 MB"
    ${Else}
        ${NSD_SetText} $SpaceAvailableLabel "Space available: unknown"
    ${EndIf}
FunctionEnd

Function BrowseDirectory
    nsDialogs::SelectFolderDialog "Select Installation Folder" $INSTDIR
    Pop $0
    ${If} $0 != error
        StrCpy $INSTDIR $0
        ${NSD_SetText} $DirRequest $INSTDIR
        Call UpdateSpaceAvailable
    ${EndIf}
FunctionEnd

Function NormalizeInstallDir
    # Remove trailing backslash if present
    StrCpy $1 $INSTDIR 1 -1
    ${If} $1 == "\"
        StrCpy $INSTDIR $INSTDIR -1
    ${EndIf}
    
    # Get the last folder name
    ${GetFileName} $INSTDIR $2
    
    # If last folder is not the product name, append it
    ${If} $2 != "${INFO_PRODUCTNAME}"
        StrCpy $INSTDIR "$INSTDIR\${INFO_PRODUCTNAME}"
    ${EndIf}
FunctionEnd

Function DirectoryAndOptionsPageLeave
    ${NSD_GetText} $DirRequest $INSTDIR
    
    # Remove trailing backslash if present
    StrCpy $1 $INSTDIR 1 -1
    ${If} $1 == "\"
        StrCpy $INSTDIR $INSTDIR -1
    ${EndIf}
    
    # Check for root directory (e.g., C:, D:, X:)
    StrLen $0 $INSTDIR
    ${If} $0 <= 2
        MessageBox MB_OK|MB_ICONEXCLAMATION "You cannot install to this location.$\r$\nUsing default installation directory."
        StrCpy $INSTDIR "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
        ${NSD_SetText} $DirRequest $INSTDIR
        Abort
    ${EndIf}
    
    # Get the part after the drive letter and check against forbidden directories
    StrCpy $5 $INSTDIR "" 2
    ${If} $5 == "\Windows"
    ${OrIf} $5 == "\Users"
    ${OrIf} $5 == "\Config.Msi"
    ${OrIf} $5 == "\ProgramData"
    ${OrIf} $5 == "\$Recycle.Bin"
    ${OrIf} $5 == "\OneDriveTemp"
        MessageBox MB_OK|MB_ICONEXCLAMATION "You cannot install to this location.$\r$\nUsing default installation directory."
        StrCpy $INSTDIR "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
        ${NSD_SetText} $DirRequest $INSTDIR
        Abort
    ${EndIf}
    
    # Normalize the installation directory (append product name if needed)
    Call NormalizeInstallDir
    
    # Update the text field to show the normalized path
    ${NSD_SetText} $DirRequest $INSTDIR
    
    ${NSD_GetState} $DesktopCheckbox $CreateDesktopShortcut
FunctionEnd

Section "!${INFO_PRODUCTNAME}" MainSection
    SectionIn RO  # This section is required and cannot be deselected
    !insertmacro wails.setShellContext

    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    
    # Create desktop shortcut if checkbox was checked
    ${If} $CreateDesktopShortcut == ${BST_CHECKED}
        CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    ${EndIf}

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
SectionEnd

Section "un.${INFO_PRODUCTNAME}" un.MainSection
    SectionIn RO  # This section is required and cannot be deselected
    !insertmacro wails.setShellContext

    # Remove WebView2 Data (cache) of all users - ALWAYS
    DetailPrint "# Removing WebView2 cache for all users..."
    FindFirst $0 $1 "C:\Users\*"
    ${DoWhile} $1 != ""
        ${If} $1 != "."
        ${AndIf} $1 != ".."
        ${AndIf} $1 != "Public"
        ${AndIf} $1 != "Default"
        ${AndIf} $1 != "Default User"
        ${AndIf} $1 != "All Users"
            ${If} ${FileExists} "C:\Users\$1\AppData\Roaming\${PRODUCT_EXECUTABLE}\*.*"
                #DetailPrint "  Removing: C:\Users\$1\AppData\Roaming\${PRODUCT_EXECUTABLE}"
                RMDir /r /REBOOTOK "C:\Users\$1\AppData\Roaming\${PRODUCT_EXECUTABLE}"
            ${EndIf}
        ${EndIf}
        FindNext $0 $1
    ${Loop}
    FindClose $0

    # Remove local appdata (log, temp files) of all users - ALWAYS
    DetailPrint "# Removing app temp/log files for all users..."
    FindFirst $0 $1 "C:\Users\*"
    ${DoWhile} $1 != ""
        ${If} $1 != "."
        ${AndIf} $1 != ".."
        ${AndIf} $1 != "Public"
        ${AndIf} $1 != "Default"
        ${AndIf} $1 != "Default User"
        ${AndIf} $1 != "All Users"
            ${If} ${FileExists} "C:\Users\$1\AppData\Local\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}\*.*"
                #DetailPrint "  Removing: C:\Users\$1\AppData\Local\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
                RMDir /r /REBOOTOK "C:\Users\$1\AppData\Local\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
            ${EndIf}
        ${EndIf}
        FindNext $0 $1
    ${Loop}
    FindClose $0
    
    # Remove installation directory
    DetailPrint "# Removing installation..."
    RMDir /r /REBOOTOK $INSTDIR

    # Remove shortcuts
    DetailPrint "# Removing shortcut(s)..."
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro wails.deleteUninstaller
SectionEnd

Section /o "un.Remove user data and configurations" un.UserData
    # Get current user's profile
    ReadEnvStr $0 "USERPROFILE"
    
    # Remove application data folders (configurations, hosts, credentials) ONLY for current user
    DetailPrint "Removing configuration files for current user..."
    DetailPrint "  Target: $0\AppData\Local\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
    RMDir /r /REBOOTOK "$0\AppData\Local\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
SectionEnd

LangString DESC_un_MainSection ${LANG_ENGLISH} "Uninstall ${INFO_PRODUCTNAME}"
LangString DESC_un_UserData ${LANG_ENGLISH} "Remove all user data, configurations, host settings and saved credentials. Uncheck this to keep your data for future installations."

!insertmacro MUI_UNFUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${un.MainSection} $(DESC_un_MainSection)
  !insertmacro MUI_DESCRIPTION_TEXT ${un.UserData} $(DESC_un_UserData)
!insertmacro MUI_UNFUNCTION_DESCRIPTION_END
