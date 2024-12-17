local utils = require("utils")
local M = {}

---@return string _ path to go installed binary
function M.dir()
    -- Check if go bin dir is defined
    local goBin = vim.fn.getenv("GOBIN")
    if goBin ~= vim.NIL then return goBin end

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
---@return boolean success, string|nil output or error
local function run_job(job)
    if not job then return false, "No job passed to run_job" end

    if not job.cmd or type(job.cmd) ~= "string" then
        return false, "Invalid or missing command in job"
    end

    -- Result of the run_job
    local success = true

    local uv = vim.loop
    local stdout = uv.new_pipe(false)
    local stderr = uv.new_pipe(false)
    local handle
    local result, err_output = {}, {}
    local done = false

    local cmd = job.cmd
    local args = job.args or {}

    -- Start the process
    handle, _ = uv.spawn(cmd, {args = args, stdio = {nil, stdout, stderr}},
                         function(exit_code, _)
        -- Cleanup after process exit
        stdout:close()
        stderr:close()
        handle:close()
        done = true

        if exit_code ~= 0 then
            success = false
            table.insert(err_output, "Command failed with exit code: " .. exit_code)
        end
    end)

    if not handle then
        success = false
        table.insert(err_output, "Failed to start job: " .. cmd)
    end

    -- Read stdout
    stdout:read_start(function(err, data)
        if err then
            success = false
            table.insert(err_output, "Error reading stdout: " .. err)
        elseif data then
            table.insert(result, data)
        end
    end)

    -- Read stderr
    stderr:read_start(function(err, data)
        if err then
            success = false
            table.insert(err_output, "Error reading stderr: " .. err)
        elseif data then
            table.insert(err_output, data)
        end
    end)

    -- Wait for the job to complete using a coroutine
    vim.wait(100000, function() return done end, 50) -- Timeout: 100 seconds, Check interval: 50ms

    if not done then return false, "Job timed out" end
    if not success then
        return false, table.concat(err_output, "\n") -- Success with output
    end

    -- If all succededs, return true and result
    return true, table.concat(result, "\n") -- Success with output

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
            utils.log_debug("Dependency: " .. binary .. " is installed at: " .. path)
        else
            utils.log_error("Dependency: " .. binary .. " is not installed or not in PATH.")
            successfully = false
        end
    end
    return successfully
end

--- Execute installation of lazydbrix go backend
---@param source string The github repo path to install
---@return boolean success
function M.exec(source)
    vim.fn.mkdir(M.output(), "p")

    local job = {cmd = "go", args = {"install", source}}

    local success, output_or_error = run_job(job)
    if success then
        utils.log_debug("Job completed successfully: " .. output_or_error)
    else
        utils.log_error("Job failed: " .. output_or_error)
    end
    return success
end

return M
