local M = {}

local util = require("praxis.util")

function M.available()
  return vim.fn.executable(util.bin()) == 1
end

local function xdg_data()
  local xdg = os.getenv("XDG_DATA_HOME")
  if not xdg or xdg == "" then
    xdg = os.getenv("HOME")
    if not xdg or xdg == "" then return nil end
    xdg = xdg .. "/.local/share"
  end
  return xdg
end

local function first_time()
  local xdg = xdg_data()
  if not xdg then return true end
  return vim.fn.filereadable(xdg .. "/praxis/progress.json") == 0
end

-- Close every open Praxis buffer before re-entering, so a fresh screen always
-- starts from a clean slate. Running :Praxis is an explicit return to the entry
-- point; any in-progress challenge is intentionally abandoned at that point.
local function close_orphans()
  for _, b in ipairs(vim.api.nvim_list_bufs()) do
    if util.is_scratch(b) then
      pcall(vim.api.nvim_buf_delete, b, { force = true })
    end
  end
end

function M.show(opts)
  local cur = vim.api.nvim_get_current_buf()
  if not util.is_scratch(cur) then
    vim.g.praxis_return_buf = cur
  end

  close_orphans()

  if not M.available() then
    require("praxis.ui").recovery("Praxis isn't installed.", {
      "Install or build Praxis, then run :Praxis again.",
      "",
      "[Enter] or [q] Back.",
    })
    return
  end

  local is_challenge = opts and opts.args and opts.args ~= ""

  if is_challenge then
    require("praxis.challenge").open(opts.args)
  elseif first_time() then
    require("praxis.onboarding").open()
  else
    require("praxis.hub").open()
  end
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })

return M
