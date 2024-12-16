local install = require("install")
local utils = require("utils")

--- Global variable to store the module
local M = {}

--- Lazydbrix class definition
local Lazydbrix = {}
Lazydbrix.__index = Lazydbrix

--- Constructor for Lazydbrix
---@param opts table options to use
---@return metatable Lazydbrix instance of a Lazydbrix
function Lazydbrix.newLazydbrix(opts)
    opts = opts or {}
    local self = setmetatable({}, Lazydbrix)

    self.dependencies = opts.dependencies
    self.file = opts.file
    self.bin = opts.bin
    self.source = opts.source

    self.clustyerSelectionTbl = {}
    return self
end

--- Get the environmental variables for cluster selection
---@return table _ with profile, clusterName and clusterID
function Lazydbrix:getClusterSelections()
    return {
        profile = vim.fn.getenv("PROFILE"),
        clusterName = vim.fn.getenv("CLUSTER_NAME"),
        clusterID = vim.fn.getenv("CLUSTER_ID")
    }
end

--- Set the environmental variables for cluster selection as an attribute
--- in the Lazydbrix instance
function Lazydbrix:setClusterSelections()
    self.clustyerSelectionTbl = self:getClusterSelections()
    self:notifyClusterSelection()
end

--- Function to notify the cluster selection
function Lazydbrix:notifyClusterSelection()
    utils.log_info("Cluster selected:\n" ..
                       vim.inspect(self.clustyerSelectionTbl))
end

--- Function to install lazydbrix
function Lazydbrix:install()
    local allInstalled = install.runVarifyBin(self.dependencies)
    if not allInstalled then
        utils.log_error("Stopping lazydbrix install")
        return
    end

    utils.log_debug("All dependencies exists")
    install.exec(self.source)

    local lazydbrixInstalled  = install.runVarifyBin({self.source})
    if not lazydbrixInstalled  then
        utils.log_error("Installation of lazydbrix was not successfully!")
    else
        utils.log_info("Installation of lazydbrix was successfully, dbrix away!")
    end
end

--- Function to open Floaterm with the command
---@return nil
function Lazydbrix:open()
    if not self.bin then
        utils.log_error("No command specified for Lazydbrix")
        return
    end
    local term_cmd = string.format(
                         ":FloatermNew --width=0.9 --height=0.9 %s -nvim %s",
                         self.bin, self.file)
    vim.cmd(term_cmd)

    -- Autocommand to source the output file on closing
    vim.api.nvim_create_autocmd("TermClose", {
        desc = [[Source the Databricks environmental variables
from the output file, at terminal close event]],
        once = true,
        callback = function()
            vim.cmd(":source " .. self.file)
            self:setClusterSelections()
        end
    })
end

-- Create an instance of Lazydbrix
local lazydbrix = Lazydbrix.newLazydbrix({
    dependencies = {"go"},
    file = install.file(),
    bin = install.bin(),
    source = "github.com/phdah/lazydbrix/cmd/lazydbrix@latest"
})

M.lazydbrix = lazydbrix

return M
