# emoji server and pasting tool

a simple image file server and tool for pasting emoji on discord and wherever else

## emoji server (main.go)

### api

**Authorization:** Bearer `$MANAGE_KEY`

- `OPTIONS /`: list all emoji
- `POST /`: upload new emoji
- `DELETE /:name`: delete emoji by name

**Public**

- `GET /:name`: get emoji by name
- `GET /:name/:size`: get emoji by name and resize

## emojitool (tool/main.go)

_for linux + rofi + xdotool setups only_

### configuration

`$EMOJITOOL_CONFIG_FILE` or `$HOME/.config/emojitool`

```sh
BASE_URL="https://emoji.example.com"
MANAGE_KEY="SECRET_MANAGE_KEY"
# EMOJITOOL_CACHE_FILE="$HOME/.cache/emojitool" # default
# EMOJITOOL_LOG_FILE="$HOME/.cache/emojitool.log" # default
```

### cache file

`$EMOJITOOL_CACHE_FILE` or `$HOME/.cache/emojitool`
