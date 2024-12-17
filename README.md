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
:lua require("lazydbrix").install()
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

These are the functions to be called from `require("lazydbrix")`
| Functions | Description |
| :--- | --- |
| `install()` | This installs `lazydbrix` |
| `open()` | This opens `lazydbrix` in a floating window |
| `setup()` | Update configurations to be used for the `Lazydbrix` class |

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

For lazy loading, set the `keys` object to suitable keymapping as well as `filetype` to python. For sourcing any previous set clusters, set `sourceOnStart = true`:

```lua
{
    'phdah/lazydbrix',
    ft = {"python"},
    opts = {sourceOnStart = true},
    keys = {
        {
            "<leader>do", ':lua require("lazydbrix").open()<CR>',
            'n'
        }
    },
    dependencies = {"voldikss/vim-floaterm"}
}
```

### 🛠️ Setup

These are all the available and default configurations (found in `defaults.lua`) that can be passed to the `setup(opts)` function:
```lua
{
    sourceOnStart = false, -- Boolean | source the output file on startup
    dependencies = {"go"}, -- table | dependencies for running lazydbrix:install()
    branch = "main", -- string | which branch to install lazydrix binary from, usefull for debuggin

    -- Only change the delow if you know what you're doing
    file = install.file(), -- string | output file for cluster selection, defaults to ~/.cache/nvim/lazydbrix/cluster_selection.nvim
    bin = install.bin() -- string | path to installed lazydbrix binary, defaults to ~/go/bin/lazydbrix (see install.bin() for more info)
}
```

## 📶 Roadmap

| Feature | Status |
| ------- | ------ |
| `tbd`   | 🟡     |
