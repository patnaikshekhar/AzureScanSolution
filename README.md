# Azure AntiVirus Solution (using Windows Defender)

Often when allowing customers to upload files via your website there is a need to scan the file for viruses before the file can be accessed by your internal team who would be processing the file. 

This is a proof of concept that creates such an architecture. The proof of concept consists of a frontend (nodejs / express) that is used by users to upload a file. When the file is uploaded it is sent to a Quarantine container where it triggers a logic app. The logic app in turn then invokes an API (via HTTP). The API is a web service (App Service) written in Go which sits in a VM behind an application gateway. The Go web service then invokes the Windows Defender executable (MpCmdRun.exe) to run a custom scan passing in the file. If the file has a virus it is moved to the Virus container else it is moved to the Clean container.

This is the overall architecture:

![Architecture](https://raw.githubusercontent.com/patnaikshekhar/AzureScanSolution/windows_defender/architecture.png)

This is how the logic app needs to be constructed so that it is triggered when a blob is added to the container

![Logic App](https://raw.githubusercontent.com/patnaikshekhar/AzureScanSolution/windows_defender/logicApp.png)