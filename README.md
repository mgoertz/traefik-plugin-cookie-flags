# traefik-plugin-cookie-flags

## Configuration

### Static

```yaml
experimental:
  plugins:
    cookieFlags:
      modulename: "github.com/Lambda-IT/traefik-plugin-cookie-flags"
      version: "v0.1.0" #replace with newest version
```

### Dynamic

To configure the plugin you should create a middleware in your dynamic configuration as explained here. The following example creates and uses the cookie path prefix middleware plugin to add the prefix "/foo" to the cookie paths:

```yaml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares:
        - "cookieFlags"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    cookieFlags:
      plugin:
        cookieFlags:
          sameSite: "None"
```

Inspired by https://github.com/SchmitzDan/traefik-plugin-cookie-path-prefix
