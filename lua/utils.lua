local M = {}

M.log_error = function(mes) vim.notify("[lazydbrix]: " .. mes, vim.log.levels.ERROR) end
M.log_info = function(mes) vim.notify("[lazydbrix]: " .. mes, vim.log.levels.INFO) end

return M
