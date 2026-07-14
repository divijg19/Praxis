local M = {}

function M.available()
  return vim.fn.executable("praxis") == 1
end

local function first_time()
  local xdg = os.getenv("XDG_DATA_HOME")
  if not xdg or xdg == "" then
    xdg = os.getenv("HOME") .. "/.local/share"
  end
  return vim.fn.filereadable(xdg .. "/praxis/stats.json") == 0
end

local function close_orphans()
  local cur = vim.api.nvim_get_current_buf()
  if vim.api.nvim_buf_get_name(cur):match("Praxis") then
    local rb = vim.g.praxis_return_buf
    if rb and vim.api.nvim_buf_is_valid(rb) then
      pcall(vim.api.nvim_set_current_buf, rb)
    end
  end
  for _, b in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_valid(b) and vim.api.nvim_buf_get_name(b):match("Praxis") then
      pcall(vim.api.nvim_buf_delete, b, { force = true })
    end
  end
end

function M.show(opts)
  local cur = vim.api.nvim_get_current_buf()
  local cur_name = vim.api.nvim_buf_get_name(cur)
  if cur_name == "" or not string.match(cur_name, "Praxis") then
    vim.g.praxis_return_buf = cur
  end

  close_orphans()

  if not M.available() then
    require("praxis.ui").recovery("Praxis isn't installed.", {
      "Install or build Praxis, then restart Neovim.",
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
