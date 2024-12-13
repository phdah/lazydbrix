local install = require("install")
local utils = require("utils")

-- Global variable to store the last output
local M = {}

-- Lazydbrix class definition
local Lazydbrix = {}
Lazydbrix.__index = Lazydbrix

-- Constructor for Lazydbrix
function Lazydbrix.newLazydbrix(opts)
    opts = opts or {}
    local self = setmetatable({}, Lazydbrix)
    self.cmd = opts.cmd
    self.file = opts.file
    self.bin = opts.bin
    self.clustyerSelectionTbl = self:getClusterSelections()
    self.dependencies = opts.dependencies
    return self
end

function Lazydbrix:getClusterSelections()
    return {
        profile = vim.fn.getenv("PROFILE"),
        clusterName = vim.fn.getenv("CLUSTER_NAME"),
        clusterID = vim.fn.getenv("CLUSTER_ID")
    }
end

-- Function to print the cluster selection
function Lazydbrix:printClusterSelection()
    print(vim.inspect(self.clustyerSelectionTbl))
end

-- Function to install lazydbrix
function Lazydbrix:install()
    install.exec(self.dependencies)
end

-- Function to open Floaterm with the command
function Lazydbrix:open()
    if not self.bin then
        utils.log_error("No command specified for Lazydbrix")
        return
    end
    local term_cmd = string.format(":FloatermNew --width=0.9 --height=0.9 %s -nvim %s", self.bin, self.file)
    vim.cmd(term_cmd)

    -- Autocommand to source the output file on closing
    vim.api.nvim_create_autocmd("TermClose", {
        desc =
[[Source the Databricks environmental variables
from the output file, at terminal close event]],
        once = true,
        callback = function()
            vim.cmd(":source " .. self.file)
            self:getClusterSelections()
        end
    })
end

-- Create an instance of Lazydbrix
local lazydbrix = Lazydbrix.newLazydbrix({
    cmd = install.bin(),
    file = install.file(),
    bin = install.bin(),
    dependencies = { "go", "make"}
})

M.lazydbrix = lazydbrix

return M