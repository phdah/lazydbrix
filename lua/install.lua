local utils = require("utils")
local M = {}

-- TODO: Figure out how to not be dependent on lazy
---@return string _ path to package manager dir
function M.source_path() return vim.fn.stdpath("data") .. "/lazy/lazydbrix" end

---@return string _ path to install dir
function M.dir() return vim.fn.stdpath("data") .. "/lazydbrix" end

---@return string _ path to binary
function M.bin() return M.dir() .. "/lazydbrix" end

---@return string _ path to output file
function M.output() return vim.fn.stdpath("cache") .. "/lazydbrix" end

---@return string _ path to output file
function M.file() return M.output() .. "/cluster_selection.nvim" end

---@param job table build command for go
local run_job = function(job)
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
            utils.log_info("Command completed successfully")
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
            utils.log_info("Job run output: " .. data)
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

---@brief Verify if a binary is installed
---@param binary string The name of the binary to check
---@return string|nil The path to the binary if found, or nil if not found
local varifyBin = function(binary)
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

---@brief Run binary verification for all dependencies
M.runVarifyBin = function(dependencies)
    local successfully = true
    for _, binary in ipairs(dependencies) do
        local path = varifyBin(binary)
        if path then
            utils.log_info(binary .. " is installed at: " .. path)
        else
            utils.log_error(binary .. " is not installed or not in PATH.")
            successfully = false
        end
    end
    return successfully
end

function M.exec(dependencies)
    local allInstalled = M.runVarifyBin(dependencies)
    if not allInstalled then
        utils.log_error("Stopping lazydbrix install")
        return
    end
    vim.fn.mkdir(M.dir(), "p")
    vim.fn.mkdir(M.output(), "p")

    local job = {
        cmd = "make",
        args = {"-C", M.source_path(), "build", "BIN=" .. M.bin()}
    }

    utils.log_debug("Installing lazydbrix with: " .. vim.inspect(job))

    run_job(job)
end

return M
