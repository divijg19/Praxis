local ui = require('praxis.ui')
local challenge = require('praxis.challenge')

local M = {}

function M.open()
  local buf = ui.create_buffer("Praxis")
  ui.set_lines(buf, {
    "Welcome to Praxis.",
    "",
    "Praxis teaches Vim through deliberate practice.",
    "",
    "You'll begin with a simple movement exercise.",
    "",
    "Press Enter to begin.",
  })
  ui.set_modifiable(buf, false)
  vim.api.nvim_set_current_buf(buf)

  vim.keymap.set("n", "<CR>", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    challenge.open("motion_rush")
  end, { buffer = buf, nowait = true, silent = true })
end

return M
