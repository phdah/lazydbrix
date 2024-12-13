<h1 align="center">
  lazydbrix
</h1>
<p align="center">
A simple, minimalistic, easy plugin to work with Databricks & Pyspark locally in Neovim
</p>

## 📦 Install

Currently only supports `lazy` package manager:

```lua
{
    'phdah/lazydbrix',
    dependencies = {"voldikss/vim-floaterm"}
    -- NOTE: Uses both go and make to install.
    -- Make sure they are present on the system.
}
```

After fetching the plugin, please run the following command in Neovim to install
the binary:

```vim
:lua require("lazydbrix").lazydbrix:install()
```

> [!IMPORTANT]
> This has to be run after every update through `lazy`.

### 📋 Requirements
- `make`
- `go`
- A Databricks config present at: `~/.databrickscfg`. With profiles to all needed workspaces:
```bash
[DEFAULT]
host = <your_host>
token = <your_token>
cluster_id = <your_cluster_id>
org_id = <your_org_id>
jobs-api-version = 2.1

[other_profile]
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
| `<C-c>` | Select the currently hovering cluster. This also exits `lazydbrix` |
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
| --- | --- |
| `tbd` | 🟡 |
