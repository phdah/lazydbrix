local Terminal = require('toggleterm.terminal').Terminal

-- Global variable to store the last output
local M = {}

-- Lazydbrix class definition
local Lazydbrix = setmetatable({}, {__index = Terminal})
Lazydbrix.__index = Lazydbrix

-- Constructor for Lazydbrix
function Lazydbrix.newLazydbrix(opts, terminalOpts)
    opts = opts or {}
    local self = setmetatable(Terminal:new(terminalOpts), Lazydbrix)
    self.cmd = opts.cmd or ""
    self.lazydbrixClusterSelectionStr = ""
    self.lazydbrixClusterSelectionTbl = {}
    return self
end

-- Method to handle stdout
function Lazydbrix:onLazydbrixStdout(_, _, data, _)
    for _, line in ipairs(data) do
        if line ~= "" then self.lazydbrixClusterSelectionStr = line end
    end
end

function Lazydbrix:decodeJSON()
    local success, result = pcall(vim.fn.json_decode,
                                  self.lazydbrixClusterSelectionStr)
    if success then
        self.lazydbrixClusterSelectionTbl = result
    end
end

-- Function to print the cluster selection
function Lazydbrix:printClusterSelection()
    print(vim.inspect(Lazydbrix.lazydbrixClusterSelectionTbl))
end

-- Create an instance of Lazydbrix
local lazydbrix = Lazydbrix.newLazydbrix({
    cmd = "~/repos/privat/lazydbrix/bin/lazydbrix"
}, {
    direction = "float",
    on_stdout = function(term, job, data, name)
        Lazydbrix:onLazydbrixStdout(term, job, data, name)
        Lazydbrix:decodeJSON()
    end
})

M.Lazydbrix = lazydbrix

return M
