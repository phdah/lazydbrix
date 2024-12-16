local utils = require("utils")
local M = {}

---@return string _ path to go installed binary
function M.dir()
    -- Check if go bin dir is defined
    local goBin = vim.fn.getenv("GOBIN")
    if goBin ~= vim.NIL then return goBin .. "/go/bin" end

    -- Default to home dir
    return vim.fn.getenv("HOME") .. "/go/bin"
end

---@return string _ path to binary
function M.bin() return M.dir() .. "/lazydbrix" end

---@return string _ path to output file
function M.output() return vim.fn.stdpath("cache") .. "/lazydbrix" end

---@return string _ path to output file
function M.file() return M.output() .. "/cluster_selection.nvim" end

---@param job table build command for go
---@return nil
local function run_job(job)
    if not job then
        utils.log_error("No job passed to run_job")
        return
    end

    if not job.cmd or type(job.cmd) ~= "string" then
        utils.log_error("Invalid or missing command in job")
        return
    end

    local cmd = job.cmd
    local args = job.args or {}

    utils.log_debug("Running command: " .. cmd .. " " .. table.concat(args, " "))

    -- Create pipes for stdout and stderr
    local uv = vim.loop
    local stdout = uv.new_pipe(false)
    local stderr = uv.new_pipe(false)

    -- Start the process
    handle, pid = uv.spawn(cmd, {args = args, stdio = {nil, stdout, stderr}},
                           function(exit_code, signal)
        -- Cleanup after the process exits
        handle:close()
        stdout:close()
        stderr:close()

        if exit_code == 0 then
            utils.log_debug("Command completed successfully")
        else
            utils.log_error("Command exited with code: " .. exit_code ..
                                " and signal: " .. signal)
        end
    end)

    if not handle then
        utils.log_error("Failed to start job: " .. cmd)
        return
    end

    utils.log_debug("Process started with PID: " .. pid)

    -- Capture stdout
    stdout:read_start(function(err, data)
        if err then
            utils.log_error("Error reading stdout: " .. err)
        elseif data then
            utils.log_debug("Job run output: " .. data)
        end
    end)

    -- Capture stderr
    stderr:read_start(function(err, data)
        if err then
            utils.log_error("Error reading stderr: " .. err)
        elseif data then
            utils.log_error("Error output: " .. data)
        end
    end)
end

--- Verify if a binary is installed
---@param binary string The name of the binary to check
---@return string|nil The path to the binary if found, or nil if not found
local function varifyBin(binary)
    local cmd = "which " .. binary .. " 2>/dev/null" -- Suppress stderr to keep output clean
    local handle = io.popen(cmd)
    if not handle then return nil end

    local result = handle:read("*all")
    handle:close()

    -- Trim trailing newline and return the result, or nil if not found
    result = result and result:gsub("\n$", "") or nil
    if result == "" then return nil end
    return result
end

--- Run binary verification for all dependencies
---@param dependencies table The name of the binary to check
---@return boolean successfully if all dependencies exists
function M.runVarifyBin(dependencies)
    local successfully = true
    for _, binary in ipairs(dependencies) do
        local path = varifyBin(binary)
        if path then
            utils.log_debug("Dependency" .. binary .. " is installed at: " .. path)
        else
            utils.log_error("Dependency" .. binary .. " is not installed or not in PATH.")
            successfully = false
        end
    end
    return successfully
end

--- Execute installation of lazydbrix go backend
---@param source string The github repo path to install
---@return nil
function M.exec(source)
    vim.fn.mkdir(M.output(), "p")

    local job = {cmd = "go", args = {"install", source}}

    run_job(job)
end

return M
