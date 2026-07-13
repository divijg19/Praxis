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
    require("praxis.challenge").open(opts.args)
  elseif first_time() then
    require("praxis.onboarding").open()
  else
    require("praxis.hub").open()
  end
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })

return M
