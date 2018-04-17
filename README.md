# Kubernetes Management API

## How to run?

1. Build + install API
```shell
go install
```

### macOS daemon ([launchd](http://www.launchd.info/))

1. Create a new Launch Agent (e.g. `~/Library/LaunchAgents/com.platzii.KubernetesManagementAPI.plist`)
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.platzii.KubernetesManagementAPI.plist</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/Username/go/bin/kubernetes-management-api</string>
        <string>--kubectl</string>
        <string>/usr/local/bin/kubectl</string>
        <string>--host</string>
        <string>localhost</string>
    </array>
    <key>RunAtLoad</key>
    <true />
</dict>
</plist>
```

If needed, change the path to the `kubernetes-management-api` and `kubectl` binaries.

2. Load Job
```shell
launchctl load ~/Library/LaunchAgents/com.platzii.KubernetesManagementAPI.plist
```

This will run the **kubernetes-management-api** on boot.

## API ref

### /proxy

#### /proxy/list

#### /proxy/info

#### /proxy/start

#### /proxy/stop
