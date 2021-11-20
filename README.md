# ENV - a tiny .env wrapper

> this is probably not much to use to anybody else but feel free to do whatever you want with it

```go
import "go.atrox.dev/env"
```

### Notes:

- `.env` file is automatically loaded, if there is no `.env` file this wrapper does not care - no error is thrown
- `APP_ENV` must be set to the apps current environment (`development`/`staging`/`production`, ...)
- get environment variables with `Get(key)` and `GetDefault(key, default)`
- check if production with `IsProduction()` or development with `IsDevelopment()`
