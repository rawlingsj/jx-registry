### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-registry/releases/download/v{{.Version}}/jx-registry-linux-amd64.tar.gz | tar xzv 
sudo mv jx-registry /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-registry/releases/download/v{{.Version}}/jx-registry-darwin-amd64.tar.gz | tar xzv
sudo mv jx-registry /usr/local/bin
```

