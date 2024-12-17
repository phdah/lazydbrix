local install = require("install")

return {
    sourceOnStart = false, -- Boolean | source the output file on startup
    dependencies = {"go"}, -- table | dependencies for running lazydbrix:install()
    branch = "main", -- string | which branch to install lazydrix binary from

    -- Only change the delow if you know what you're doing
    file = install.file(), -- string | output file for cluster selection, defaults to ~/.cache/nvim/lazydbrix/cluster_selection.nvim
    bin = install.bin(), -- string | path to installed lazydbrix binary, defaults to ~/go/bin/lazydbrix (see install.bin() for more info)

    -- Below are things you should never change
    repo = "github.com/phdah/lazydbrix/cmd/lazydbrix" -- string | github repo to install lazydbrix binary from
}
