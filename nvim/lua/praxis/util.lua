local M = {}

-- Canonical invocation name. The CLI binary and the Neovim entry point share
-- this lowercase name. Vim requires user commands to start with an uppercase
-- letter, so the Neovim command is derived from it at that single boundary.
M.name = "praxis"

function M.command()
  local name = M.name:gsub("^%l", string.upper)
  return name
end

-- Run the Neovim entry point. All internal invocations go through here so the
-- command name cannot drift from the canonical lowercase name.
function M.invoke(args)
  local suffix = args and args ~= "" and (" " .. args) or ""
  vim.cmd(M.command() .. suffix)
end

-- The certification gate injects an absolute binary path; interactive use
-- falls back to `praxis` on $PATH. Every CLI call must go through here.
function M.bin()
  return vim.env.PRAXIS_BIN or "praxis"
end

local function is_scratch(buf)
  if not vim.api.nvim_buf_is_valid(buf) then return false end
  local name = vim.api.nvim_buf_get_name(buf)
  local base = vim.fn.fnamemodify(name, ":t")
  return base:match("^Praxis") ~= nil
    and vim.api.nvim_get_option_value("buftype", { buf = buf }) == "nofile"
end

function M.is_scratch(buf)
  return is_scratch(buf)
end

function M.close_current()
  local buf = vim.api.nvim_get_current_buf()
  if is_scratch(buf) then
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
  end
end

-- Describe a challenge by id, returning the decoded table or nil on failure.
function M.describe(id)
  local raw = vim.fn.systemlist({ M.bin(), "describe", id })
  if vim.v.shell_error ~= 0 then return nil end
  local ok, desc = pcall(vim.fn.json_decode, table.concat(raw, ""))
  if not ok or type(desc) ~= "table" then return nil end
  return desc
end

-- Preview a hint (cost only) or reveal it. Reveal persists the rung in progress.
function M.hint(id, level, preview)
  local args = { M.bin(), "hint", id, tostring(level) }
  if preview then table.insert(args, "--preview") end
  local out = vim.fn.systemlist(args)
  if vim.v.shell_error ~= 0 then return nil end
  return out
end

-- Byte-exact buffer comparison shared by the product validator and replay.
function M.check_buffer(buf, result_lines)
  local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
  if #current ~= #result_lines then return false end
  for i = 1, #current do
    if current[i] ~= result_lines[i] then return false end
  end
  return true
end

-- Fetch an internal solution trace. Requires PRAXIS_DEV=1 in the environment.
function M.trace(id)
  local raw = vim.fn.systemlist({ M.bin(), "trace", id })
  if vim.v.shell_error ~= 0 then return nil end
  local ok, tr = pcall(vim.fn.json_decode, table.concat(raw, ""))
  if not ok or type(tr) ~= "table" then return nil end
  if type(tr.keys) ~= "string" or tr.keys == "" then return nil end
  if type(tr.start_row) ~= "number" or type(tr.start_col) ~= "number" then return nil end
  return tr
end

-- Convert a byte column index to a character column index within a line.
function M.byte_to_char(line, bytecol)
  return vim.fn.strchars(string.sub(line, 1, bytecol))
end

-- Run a praxis subcommand capturing stdout. Returns nil on a non-zero exit so
-- callers can show recovery instead of interpreting an error as data.
function M.systemlist(args)
  local out = vim.fn.systemlist(vim.list_extend({ M.bin() }, args))
  if vim.v.shell_error ~= 0 then return nil end
  return out
end

-- Send a JSON payload to a praxis subcommand and return stdout, or nil on
-- non-zero exit. Used by the verification-based completion path.
function M.send(args, payload)
  local out = vim.fn.system(vim.list_extend({ M.bin() }, args), vim.fn.json_encode(payload))
  if vim.v.shell_error ~= 0 then return nil end
  return out
end

-- Run a praxis subcommand. Returns true when it exited successfully, so callers
-- can surface failures instead of silently assuming progress was saved.
function M.praxis(args)
  vim.fn.system(vim.list_extend({ M.bin() }, args))
  return vim.v.shell_error == 0
end

-- Close the current Praxis buffer and open `id` (or the hub when id is empty)
-- through the praxis command, so navigation state resets cleanly. This is the
-- single open-next/back path shared by the hub and the challenge screen.
function M.continue(id)
  M.close_current()
  M.invoke(id)
end

-- Return to the buffer that was active before Praxis opened (the "return
-- buffer"), or to the hub if none. Single owner of the back-navigation path
-- used by the recovery screen and the hub.
function M.return_to_previous()
  M.close_current()
  local rb = vim.g.praxis_return_buf
  if rb and vim.api.nvim_buf_is_valid(rb) then
    vim.api.nvim_set_current_buf(rb)
  else
    M.invoke("")
  end
end

return M
