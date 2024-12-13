local utils = require("utils")
local M = {}

---@return string _ path to package manager dir
function M.source_path()
return "~/repos/privat/lazydbrix" end
-- return vim.fn.stdpath("data") .. "/lazy/lazydbrix" end

---@return string _ path to install dir
function M.dir()
return vim.fn.stdpath("data") .. "/lazydbrix" end

---@return string _ path to binary
function M.bin()
return M.source_path() .. "/bin/lazydbrix" end
-- return M.dir() .. "/bin" end

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

    utils.log_info("Running command: " .. cmd .. " " .. table.concat(args, " "))

    -- Create a new job for running the CLI tool
    local run = vim.fn.jobstart({cmd, unpack(args)}, {
        stdout_buffered = true,
        stderr_buffered = true,
        on_exit = function(_, exit_code)
            if exit_code == 0 then
                utils.log_info("Command completed successfully")
            else
                utils.log_error("Command exited with code: " .. exit_code)
            end
        end,
        on_stdout = function(_, data)
            if data then
                utils.log_info("Output: " .. table.concat(data, "\n"))
            end
        end,
        on_stderr = function(_, data)
            if data then
                utils.log_error("Error: " .. table.concat(data, "\n"))
            end
        end,
    })

    if run == 0 then
        utils.log_error("Failed to start job: " .. cmd)
    end
end

function M.exec()
    local installBinary = M.dir()
    vim.fn.mkdir(installBinary, "p")
    vim.fn.mkdir(M.output(), "p")

    -- Define the job details
    local job = {
        cmd = "go",
        args = {"build", "-o", installBinary, M.source_path()},
    }

    utils.log_info("Installing lazydbrix with: " .. vim.inspect(job))

    run_job(job)
end

return M
