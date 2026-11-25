local M = {}

function M.create_buffer(name)
  local target = name
  local suffix = 2
  while true do
    local conflict = false
    for _, b in ipairs(vim.api.nvim_list_bufs()) do
      local existing = vim.fn.fnamemodify(vim.api.nvim_buf_get_name(b), ":t")
      if vim.api.nvim_buf_is_valid(b) and existing == target then
        local buftype = vim.api.nvim_get_option_value("buftype", { buf = b })
        local modified = vim.api.nvim_get_option_value("modified", { buf = b })
        if buftype == "nofile" and not modified then
          pcall(vim.api.nvim_buf_delete, b, { force = true })
        else
          conflict = true
        end
      end
    end
    if not conflict then break end
    target = name .. " " .. suffix
    suffix = suffix + 1
  end
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_buf_set_name(buf, target)
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
    require("praxis.util").return_to_previous()
  end

  vim.keymap.set("n", "<CR>", back, { buffer = buf, nowait = true, silent = true })
  vim.keymap.set("n", "q", back, { buffer = buf, nowait = true, silent = true })
end

return M
