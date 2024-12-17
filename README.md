<h1 align="center">
  lazydbrix
</h1>
<p align="center">
A simple, minimalistic, easy plugin to work with Databricks & Pyspark locally in Neovim
</p>

## 📦 Install

Using `lazy` package manager:
```lua
{
    'phdah/lazydbrix',
    dependencies = {"voldikss/vim-floaterm"}
    -- NOTE: Uses go to install.
    -- Make sure it's present on the system.
}
```

After fetching the plugin, please run the following command in Neovim to install
the binary:

```vim
:lua require("lazydbrix").lazydbrix:install()
```

This command will call `go install github.com/phdah/lazydbrix/cmd/lazydbrix@main` and put the binary in either of:
`$GOBIN/lazydbrix` or `$HOME/go/bin/lazydbrix`, dependent on if you have setup `GOBIN` or not (see [docs](https://pkg.go.dev/cmd/go#hdr-Environment_variables)).

> [!IMPORTANT]
> This has to be run after every update through `lazy`, for now.

### 📋 Requirements

- `go`
- A Databricks config present at: `~/.databrickscfg`. With profiles to all
  needed workspaces:

```bash
[DEFAULT]
host = <your_host>
token = <your_token>
cluster_id = <your_cluster_id>
org_id = <your_org_id>
jobs-api-version = 2.1

[TEST]
...
[PROD]
...
```

> [!CAUTION]
> All profiles that are in this file, needs to have valid setups, otherwise the plugin won't work

## 🚀 How to

To open the window, run:

```vim
:lua require("lazydbrix").lazydbrix:open()
```

Once inside of the floating window, this is how you navigate inside of `lazydbrix`:
| Keymaps | Description |
| :--- | --- |
| `<C-c>` | This exits `lazydbrix` |
| `<enter>` | Select the currently hovering cluster. This inly takes effect once exiting `lazydbrix` |
| `j` | Move down in  the current window |
| `k` | Move up in  the current window |
| `l` | Move down a window |
| `h` | Move up a window |
| `<Tab>` | Move to the right window |
| `<S-Tab>` | Move to the left window |

#### 💤 Lazy loading

For lazy loading, set the `keys` object to suitable key mapping:

```lua
{
    'phdah/lazydbrix',
    keys = {
        {
            "<leader>do", ':lua require("lazydbrix").lazydbrix:open()<CR>',
            'n'
        }
    },
    dependencies = {"voldikss/vim-floaterm"}
}
```

## 📶 Roadmap

| Feature | Status |
| ------- | ------ |
| `tbd`   | 🟡     |
