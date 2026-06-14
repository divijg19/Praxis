local challenge = require('praxis.challenge')
local session = require('praxis.session')
local ui = require('praxis.ui')
local onboarding = require('praxis.onboarding')

local M = {}

local function first_time()
  local xdg = os.getenv("XDG_DATA_HOME")
  if not xdg or xdg == "" then
    xdg = os.getenv("HOME") .. "/.local/share"
  end
  return vim.fn.filereadable(xdg .. "/praxis/stats.json") == 0
end

function M.show(opts)
  local is_challenge = opts and opts.args and opts.args ~= ""

  if is_challenge then
    challenge.open(opts.args)
  elseif first_time() then
    onboarding.open()
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
