local M = {}

function M.create_buffer(name)
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_buf_set_name(buf, name)
  return buf
end

function M.set_lines(buf, lines)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
end

function M.set_modifiable(buf, val)
  vim.api.nvim_set_option_value("modifiable", val, { buf = buf })
end

return M
