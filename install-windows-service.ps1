# Delete and stop the service if it already exists.
if (Get-Service gostIndex -ErrorAction SilentlyContinue) {
  $service = Get-WmiObject -Class Win32_Service -Filter "name='ez-search'"
  $service.StopService()
  Start-Sleep -s 1
  $service.delete()
}

$workdir = Split-Path $MyInvocation.MyCommand.Path

Write-Host $workdir $MyInvocation
write-Host "$workdir\ez-search.exe -c `"$workdir\config.json`" "
# Create the new service.
New-Service -name gostIndex `
  -displayName gostIndex `
  -binaryPathName "$workdir\ez-search.exe -c `"$workdir\config.json`" -wd `"$workdir`" "
# Attempt to set the service to delayed start using sc config.
Try {
  Start-Process -FilePath sc.exe -ArgumentList 'config ez-search start= delayed-auto'
}
Catch { Write-Host -f red "An error occured setting the service to delayed start." }
