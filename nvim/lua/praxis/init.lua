local challenge = require('praxis.challenge')
local session = require('praxis.session')
local ui = require('praxis.ui')

local M = {}

function M.show(opts)
  local is_challenge = opts and opts.args and opts.args ~= ""

  if is_challenge then
    challenge.open(opts.args)
  else
    local lines = vim.fn.systemlist({ "praxis" })
    table.insert(lines, "")
    table.insert(lines, "CLI Connected")
    local buf = ui.create_buffer("Praxis")
    ui.set_lines(buf, lines)
    ui.set_modifiable(buf, false)
    vim.api.nvim_set_current_buf(buf)
  end
end

function M.show_session()
  session.show()
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })
vim.api.nvim_create_user_command("PraxisSession", M.show_session, {})

return M
