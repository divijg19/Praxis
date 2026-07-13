local M = {}

local buf_seq = 0
function M.create_buffer(name)
  for _, b in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_get_name(b) == name then
      pcall(vim.api.nvim_buf_delete, b, { force = true })
    end
  end
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  if not pcall(vim.api.nvim_buf_set_name, buf, name) then
    buf_seq = buf_seq + 1
    vim.api.nvim_buf_set_name(buf, name .. " #" .. buf_seq)
  end
  return buf
end

function M.set_lines(buf, lines)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
end

function M.set_modifiable(buf, val)
  vim.api.nvim_set_option_value("modifiable", val, { buf = buf })
end

function M.recovery(title, lines)
  local buf = M.create_buffer("Praxis")
  local display = { title, "" }
  for _, l in ipairs(lines) do
    table.insert(display, l)
  end
  M.set_lines(buf, display)
  M.set_modifiable(buf, false)
  vim.api.nvim_set_current_buf(buf)

  local function back()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local rb = vim.g.praxis_return_buf
    if rb and vim.api.nvim_buf_is_valid(rb) then
      vim.api.nvim_set_current_buf(rb)
    end
  end

  vim.keymap.set("n", "<CR>", back, { buffer = buf, nowait = true, silent = true })
  vim.keymap.set("n", "q", back, { buffer = buf, nowait = true, silent = true })
end

return M
