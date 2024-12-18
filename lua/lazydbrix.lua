local install = require("install")
local config = require("defaults")
local utils = require("utils")

--- Global variable to store the module
local M = {}

--- Lazydbrix class definition
---@class Lazydbrix
---@field sourceOnStart boolean
---@field dependencies string[]
---@field file string
---@field bin string
---@field repo string
---@field branch string
local Lazydbrix = {}
Lazydbrix.__index = Lazydbrix

--- Constructor for Lazydbrix
---@param opts table options to use
---@return Lazydbrix instance of Lazydbrix
function Lazydbrix.newLazydbrix(opts)
    local self = setmetatable({}, Lazydbrix)
    for key, value in pairs(opts) do self[key] = value end

    return self
end

--- Update instance options dynamically
---@param new_opts table New options to merge into the instance
function Lazydbrix:updateOptions(new_opts)
    -- Dynamically update attributes from the new options
    for key, value in pairs(new_opts) do self[key] = value end
end

--- Get full github source
---@return string source concatination of repo and branch
function Lazydbrix:source() return self.repo .. "@" .. self.branch end

--- Source output file
---@return nil
function Lazydbrix:sourceFile()
    if self.file == "" then
        utils.log_error("Got empty filename, can't source it")
        return
    end
    if vim.loop.fs_stat(self.file) ~= nil then
        vim.cmd(":source " .. self.file)
    else
        utils.log_debug("Output file does not exist yet, checked file: " ..
                            self.file)
    end
    self:notifyClusterSelection()

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

--- Function to notify the cluster selection
function Lazydbrix:notifyClusterSelection()
    utils.log_info("Cluster selected:\n" ..
                       vim.inspect(self:getClusterSelections()))
end

--- Function to install lazydbrix
function Lazydbrix:install()
    local allInstalled = install.runVarifyBin(self.dependencies)
    if not allInstalled then
        utils.log_error("Stopping lazydbrix install")
        return
    end

    utils.log_debug("All dependencies exists")
    local success = install.exec(self:source())
    if not success then
        utils.log_error("Installation of lazydbrix was not successfully!")
        return
    end

    local lazydbrixBinaryExists = install.runVarifyBin({self.bin})
    if not lazydbrixBinaryExists then
        utils.log_error("There exists no lazydbrix binary")
        return
    end
    utils.log_info("Installation of lazydbrix was successfully, dbrix away!")
end

--- Function to open Floaterm with the command
---@return nil
function Lazydbrix:open()
    if not self.bin then
        utils.log_error("No command specified for Lazydbrix")
        return
    end
    utils.log_debug("Binary: " .. self.bin)
    local term_cmd = string.format(
                         ":FloatermNew --width=0.9 --height=0.9 %s -nvim %s",
                         self.bin, self.file)
    vim.cmd(term_cmd)

    -- Autocommand to source the output file on closing
    vim.api.nvim_create_autocmd("TermClose", {
        desc = [[Source the Databricks environmental variables
from the output file, at terminal close event]],
        once = true,
        callback = function() self:sourceFile() end
    })
end

-- Create an instance of Lazydbrix
local lazydbrix = Lazydbrix.newLazydbrix(config)

--- Setup function to overwrite default configuration
---@param user_config table User-provided configuration
function M.setup(user_config)
    if lazydbrix == nil then
        return
    else
        lazydbrix:updateOptions(user_config or {})
        if lazydbrix.sourceOnStart then lazydbrix:sourceFile() end
    end
end

--- Wrapper function of Lazydbrix:install()
function M.install()
    if lazydbrix == nil then
        return
    else
        lazydbrix:install()
    end
end

--- Wrapper function of Lazydbrix:open()
function M.open()
    if lazydbrix == nil then
        return
    else
        lazydbrix:open()
    end
end

--- Wrapper function of Lazydbrix:sourceFile()
function M.show()
    if lazydbrix == nil then
        return
    else
        lazydbrix:sourceFile()
    end
end

return M
