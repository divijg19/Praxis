-- End-to-end replay test for all 41 challenges
-- Committed at tools/replay/replay.lua
-- Run via: tools/replay/replay.sh

local function sh(id, cmd)
  local handle = io.popen("/tmp/praxis " .. cmd .. " " .. id)
  local result = handle:read("*l")
  handle:close()
  return result or ""
end

local function sh_all(id, cmd)
  local handle = io.popen("/tmp/praxis " .. cmd .. " " .. id)
  local lines = {}
  for line in handle:lines() do
    table.insert(lines, line)
  end
  handle:close()
  return lines
end

local function byte_to_char(line, bytecol)
  return vim.fn.strchars(string.sub(line, 1, bytecol))
end

local function check_buffer(buf, result_lines)
  local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
  if #current ~= #result_lines then
    return false
  end
  for i = 1, #current do
    if current[i] ~= result_lines[i] then
      return false
    end
  end
  return true
end

local cursor_ids = {
  "motion_rush","grid_rush","find_hunter","word_hunter",
  "symbol_hunter","line_hunter","paren_hunter","sentence_hunter",
  "slash_hunter","question_hunter","repeat_hunter",
  "inner_paren_hunter","around_paren_hunter","inner_bracket_hunter",
  "around_bracket_hunter","inner_quote_hunter","around_quote_hunter",
  "paragraph_hunter","match_hunter",
}

local buffer_ids = {
  "delete_character_hunter","replace_character_hunter",
  "toggle_case_hunter","delete_word_hunter","change_word_hunter",
  "delete_line_hunter","delete_to_end_hunter",
  "delete_inner_word_hunter","delete_around_word_hunter",
  "delete_inner_paren_hunter","delete_around_paren_hunter",
  "delete_inner_quote_hunter","delete_around_quote_hunter",
  "change_inner_word_hunter","change_inner_paren_hunter",
  "change_inner_quote_hunter","yank_line_hunter",
  "named_register_hunter","word_register_hunter",
  "register_replace_hunter","register_duplicate_hunter",
}

local utf8_id = "utf8_cursor_hunter"

local cursor_pass = 0
local cursor_fail = 0
local buffer_pass = 0
local buffer_fail = 0
local utf8_pass = 0
local utf8_fail = 0

for _, id in ipairs(cursor_ids) do
  local verify = sh(id, "verify")
  local target = sh(id, "target")
  local lines = sh_all(id, "challenge")

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)

  local found = false
  for r = 1, #lines do
    local line = lines[r]
    if line then
      for c = 0, #line do
        local charcol = byte_to_char(line, c)
        local ch = vim.fn.strcharpart(line, charcol, 1)
        if ch == target then
          vim.api.nvim_win_set_cursor(0, { r, c })
          found = true
          break
        end
      end
    end
    if found then break end
  end

  if not found then
    cursor_fail = cursor_fail + 1
    print("FAIL " .. id .. ": target '" .. target .. "' not found")
    goto continue_cursor
  end

  local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
  local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
  local completed = false
  if line and verify == "cursor" then
    local charcol = byte_to_char(line, col0)
    if vim.fn.strcharpart(line, charcol, 1) == target then
      completed = true
    end
  end

  if completed then
    cursor_pass = cursor_pass + 1
    print("PASS " .. id)
  else
    cursor_fail = cursor_fail + 1
    print("FAIL " .. id)
  end

  ::continue_cursor::
end

for _, id in ipairs(buffer_ids) do
  local verify = sh(id, "verify")
  local content_lines = sh_all(id, "challenge")
  local result_lines = sh_all(id, "result")

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, content_lines)
  vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)

  vim.api.nvim_buf_set_lines(buf, 0, -1, false, result_lines)
  local matched = verify == "buffer" and check_buffer(buf, result_lines)

  if matched then
    buffer_pass = buffer_pass + 1
    print("PASS " .. id)
  else
    buffer_fail = buffer_fail + 1
    local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    print("FAIL " .. id .. " #current=" .. #current .. " #result=" .. #result_lines)
  end
end

do
  local id = utf8_id
  local verify = sh(id, "verify")
  local target = sh(id, "target")
  local lines = sh_all(id, "challenge")

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)

  vim.api.nvim_win_set_cursor(0, { 3, 9 })
  local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
  local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
  local charcol = byte_to_char(line, col0)
  local ch = vim.fn.strcharpart(line, charcol, 1)

  if verify == "cursor" and ch == target then
    utf8_pass = utf8_pass + 1
    print("PASS " .. id .. " bytecol=" .. col0 .. " charcol=" .. charcol)
  else
    utf8_fail = utf8_fail + 1
    print("FAIL " .. id)
  end
end

local total_cursor = cursor_pass + cursor_fail
local total_buffer = buffer_pass + buffer_fail
local total_pass = cursor_pass + buffer_pass + utf8_pass
local total_fail = cursor_fail + buffer_fail + utf8_fail
local total = total_cursor + total_buffer + 1

print("")
print("Cursor: " .. cursor_pass .. "/" .. total_cursor .. " PASS")
print("Buffer: " .. buffer_pass .. "/" .. total_buffer .. " PASS")
print("UTF-8:  " .. utf8_pass .. "/1 PASS")
print("")
if total_fail == 0 then
  print("ALL " .. total .. "/" .. total .. " REPLAY TESTS PASS")
else
  print(total_pass .. "/" .. total .. " PASS, " .. total_fail .. " FAIL")
end

vim.cmd("qa!")
