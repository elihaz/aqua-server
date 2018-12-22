# server that listens to file info requests

To build image:

```bash
docker build {image_name} -t .
```
To Run it:
```bash
docker run -p{port}:{port} {image_name} -port {port} -protocol {HTTP/HTTPS} -format{JSON/XML}
```
Or for help:

```bash
docker run {image_name} --help
```

