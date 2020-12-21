# Blogger

Copyright 2020, The Basebandit

## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## About The Project

Blogger is a simple golang http server that serves html files from a directory.

## How to use it
- Download the project from github  

```go
go get github.com/blogger
```  

- Run it  

```bash
└──╼ $ go run blogger.go -http ":8888" -content="/home/capricorn/code/go/blogger/content/" -index
/index
/read
Server listening on address :8888
```

Replace the content path with the path to your content.

