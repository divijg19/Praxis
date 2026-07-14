local M = {}

function M.create_buffer(name)
  for _, b in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_get_name(b) == name then
      pcall(vim.api.nvim_buf_delete, b, { force = true })
    end
  end
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_buf_set_name(buf, name)
  return buf
end

function M.show(name, lines, modifiable)
	local buf = M.create_buffer(name)
	vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
	vim.api.nvim_set_option_value("modifiable", modifiable, { buf = buf })
	vim.api.nvim_set_current_buf(buf)
	return buf
end

function M.recovery(title, lines)
	local display = { title, "" }
	for _, l in ipairs(lines) do
		table.insert(display, l)
	end
	local buf = M.show("Praxis", display, false)

  local function back()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local rb = vim.g.praxis_return_buf
    if rb and vim.api.nvim_buf_is_valid(rb) then
      vim.api.nvim_set_current_buf(rb)
    else
      vim.cmd("Praxis")
    end
  end

  vim.keymap.set("n", "<CR>", back, { buffer = buf, nowait = true, silent = true })
  vim.keymap.set("n", "q", back, { buffer = buf, nowait = true, silent = true })
end

return M
