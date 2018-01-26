
# build image

```
docker build -t overlaps .
```

# run container
```
docker run -p 8080:8080 -v ./:/host_dir/ -d overlaps
```