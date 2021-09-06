# awair_exporter

Update `AWAIR_URL` with the IP address to your Awair. Make sure to set a static IP address for your Awair.

```
docker run -p 8181:8181 -e AWAIR_URL="http://1.1.1.1" -e POLL_DURATION="30s" dewski/awair_exporter:latest
```
