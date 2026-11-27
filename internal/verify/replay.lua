-- Embedded, config-free verifier used by `praxis complete` and `praxis verify`.
-- It replays a validated key log in a clean Neovim, ignores instruction-line
-- cursor targets, rejects goals already satisfied at the start, and reports
-- the shortest prefix that reaches the goal.

local desc = vim.fn.json_decode(vim.env.PRAXIS_VERIFY_DESC)
local keys_file = assert(io.open(vim.env.PRAXIS_VERIFY_KEYS, "r"))
local keys = vim.fn.json_decode(keys_file:read("*a"))
keys_file:close()
local out = vim.env.PRAXIS_VERIFY_OUT

local function replay(prefix)
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, desc.content)
  vim.api.nvim_set_option_value("modifiable", desc.verify ~= "cursor", { buf = buf })

  local row = math.max(1, math.min(desc.start_row or 3, #desc.content))
  local line = desc.content[row] or ""
  local col = math.max(0, math.min(desc.start_col or 0, #line))
  vim.api.nvim_win_set_cursor(0, { row, col })

  local before = table.concat(vim.api.nvim_buf_get_lines(buf, 0, -1, false), "\n")
  -- One typeahead batch keeps Vim counts (`3l`, `2dd`) adjacent.
  vim.api.nvim_feedkeys(vim.api.nvim_replace_termcodes(table.concat(prefix, ""), true, true, true), "x", true)
  local after = table.concat(vim.api.nvim_buf_get_lines(buf, 0, -1, false), "\n")
  local edits = before == after and 0 or 1

  local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
  local verified = false
  if desc.verify == "cursor" then
    local r, byte_col = unpack(vim.api.nvim_win_get_cursor(0))
    -- The product frontend never evaluates the instruction line (row 1).
    if not (#desc.content > 1 and r == 1) then
      local current_line = current[r] or ""
      local char_col = vim.fn.strchars(string.sub(current_line, 1, byte_col))
      verified = vim.fn.strcharpart(current_line, char_col, 1) == desc.target
    end
  else
    local expected = desc.result or {}
    verified = #current == #expected
    if verified then
      for i = 1, #current do
        if current[i] ~= expected[i] then
          verified = false
          break
        end
      end
    end
  end
  return verified, current, edits
end

local verified = false
local matched = 0
local edit_count = 0
local final_lines = {}
local already, start_lines = replay({})
final_lines = start_lines
if already then
  final_lines = start_lines
else
  for n = 1, #keys do
    local prefix = {}
    for i = 1, n do prefix[i] = keys[i] end
    local ok, lines, edits = replay(prefix)
    final_lines = lines
    if ok then
      verified = true
      matched = n
      edit_count = edits
      break
    end
  end
end

local verdict = {
  verified = verified,
  reason = verified and "goal reached" or (already and "goal already satisfied at start" or "goal not reached"),
  keys = matched,
  edits = edit_count,
  par = desc.par or 0,
  final_lines = final_lines,
}

local f = assert(io.open(out, "w"))
f:write(vim.fn.json_encode(verdict))
f:close()
vim.cmd("qa!")
