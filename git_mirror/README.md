
Usage:

```
docker build -t git-http-backend .
docker run -d -p 4080:80 -v ./repos:/git git-http-backend
```

Unauthenticated push will not work unless you enable it in repositories:

```
cd /path/to/host/gitdir
git init --bare test.git
cd test.git
git config http.receivepack true
```

Repos will then be accessible at `http://localhost:4080/<repo>.git`.
