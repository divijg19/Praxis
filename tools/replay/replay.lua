-- Replay verification for every catalog challenge: executes a canonical Vim
-- keystroke trace and checks the resulting goal state and verified key count.
-- Par is derived from the trace and cross-checked against describe output.
-- Run via: go run ./tools/verify (PRAXIS_BIN may override the binary)

local self = debug.getinfo(1, "S").source:sub(2)
local ROOT = vim.fn.fnamemodify(self, ":h:h:h")
vim.opt.runtimepath:prepend(ROOT .. "/nvim")
local util = require("praxis.util")
local bin = util.bin()

-- IDs come from the CLI so the harness cannot drift from the catalog.
local all_ids = vim.fn.systemlist({ bin, "catalog", "--ids" })
if #all_ids == 0 then
  print("FAIL: praxis catalog --ids returned no challenges")
  vim.cmd("cquit 1")
end

local results = {}
local descs = {}

for _, id in ipairs(all_ids) do
  local d = util.describe(id)
  descs[id] = d
  if type(d) ~= "table" or type(d.content) ~= "table" or #d.content == 0 then
    print("FAIL " .. id .. ": could not describe")
    results[id] = "FAIL"
    goto continue
  end

  local tr = util.trace(id)
  if not tr then
    print("FAIL " .. id .. ": no valid canonical solution trace")
    results[id] = "FAIL"
    goto continue
  end

  local row = math.max(1, math.min(tr.start_row, #d.content))
  local line = d.content[row] or ""
  local col = math.max(0, math.min(tr.start_col, #line))
  if tr.start_row < 1 or tr.start_row > #d.content then
    print("FAIL " .. id .. ": trace start row out of range")
    results[id] = "FAIL"
    goto continue
  end

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, d.content)
  vim.api.nvim_set_option_value("modifiable", d.verify ~= "cursor", { buf = buf })
  vim.api.nvim_win_set_cursor(0, { row, col })
  vim.api.nvim_feedkeys(tr.keys, "x", true)

  local key_count = vim.fn.strchars(tr.keys)
  local solved = false

  if d.verify == "cursor" then
    local r, col0 = unpack(vim.api.nvim_win_get_cursor(0))
    local current = vim.api.nvim_buf_get_lines(buf, r - 1, r, false)[1]
    if current then
      local charcol = util.byte_to_char(current, col0)
      solved = vim.fn.strcharpart(current, charcol, 1) == d.target
    end
  else
    solved = util.check_buffer(buf, d.result)
  end

  if d.layer == "Trial" and (type(d.derived_from) ~= "table" or #d.derived_from == 0) then
    print("FAIL " .. id .. ": missing derived_from")
    results[id] = "FAIL"
  elseif d.par and d.par ~= key_count then
    print("FAIL " .. id .. ": par " .. d.par .. " does not match trace keys " .. key_count)
    results[id] = "FAIL"
  elseif solved then
    print("PASS " .. id .. " (" .. key_count .. " keys)")
    results[id] = "PASS"
  else
    print("FAIL " .. id .. ": canonical trace did not solve the challenge")
    results[id] = "FAIL"
  end

  ::continue::
end

local tutorial_cursor_pass = 0
local tutorial_cursor_fail = 0
local tutorial_buffer_pass = 0
local tutorial_buffer_fail = 0
local training_pass = 0
local training_fail = 0
local trial_pass = 0
local trial_fail = 0
local utf8_pass = 0
local utf8_fail = 0
local describe_fail = 0

for _, id in ipairs(all_ids) do
  local d = descs[id]
  if type(d) ~= "table" then
    describe_fail = describe_fail + 1
    goto next_id
  end

  local result = results[id] or "FAIL"
  local is_pass = result == "PASS"

  if d.layer == "Tutorial" and id == "utf8_cursor_hunter" then
    if is_pass then utf8_pass = utf8_pass + 1 else utf8_fail = utf8_fail + 1 end
  elseif d.layer == "Tutorial" and d.verify == "cursor" then
    if is_pass then tutorial_cursor_pass = tutorial_cursor_pass + 1 else tutorial_cursor_fail = tutorial_cursor_fail + 1 end
  elseif d.layer == "Tutorial" and (d.verify == "buffer" or d.verify == "composite") then
    if is_pass then tutorial_buffer_pass = tutorial_buffer_pass + 1 else tutorial_buffer_fail = tutorial_buffer_fail + 1 end
  elseif d.layer == "Training" then
    if is_pass then training_pass = training_pass + 1 else training_fail = training_fail + 1 end
  elseif d.layer == "Trial" then
    if is_pass then trial_pass = trial_pass + 1 else trial_fail = trial_fail + 1 end
  end

  ::next_id::
end

local total_pass = tutorial_cursor_pass + tutorial_buffer_pass + utf8_pass + training_pass + trial_pass
local total_fail = tutorial_cursor_fail + tutorial_buffer_fail + utf8_fail + training_fail + trial_fail + describe_fail
local total = total_pass + total_fail

print("")
print("Tutorial (cursor): " .. tutorial_cursor_pass .. "/" .. (tutorial_cursor_pass + tutorial_cursor_fail) .. " PASS")
print("Tutorial (buffer): " .. tutorial_buffer_pass .. "/" .. (tutorial_buffer_pass + tutorial_buffer_fail) .. " PASS")
print("Tutorial (utf-8):  " .. utf8_pass .. "/" .. (utf8_pass + utf8_fail) .. " PASS")
print("Training:         " .. training_pass .. "/" .. (training_pass + training_fail) .. " PASS")
print("Trial:            " .. trial_pass .. "/" .. (trial_pass + trial_fail) .. " PASS")
print("")
if describe_fail > 0 then
  print("Describe failures: " .. describe_fail)
end
if total_fail == 0 then
  print("ALL " .. total .. "/" .. total .. " REPLAY TESTS PASS")
  vim.cmd("qa!")
else
  print(total_pass .. "/" .. total .. " PASS, " .. total_fail .. " FAIL")
  vim.cmd("cquit 1")
end
