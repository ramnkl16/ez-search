# Delete and stop the service if it already exists.
if (Get-Service ez-search -ErrorAction SilentlyContinue) {
  $service = Get-WmiObject -Class Win32_Service -Filter "name='ez-search'"
  $service.StopService()
  Start-Sleep -s 1
  $service.delete()
}