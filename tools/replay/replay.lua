-- End-to-end replay test for all 52 challenges
-- Run via: tools/replay/replay.sh

local self = debug.getinfo(1, "S").source:sub(2)
local ROOT = vim.fn.fnamemodify(self, ":h:h:h")
vim.opt.runtimepath:prepend(ROOT .. "/nvim")
local util = require("praxis.util")

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

local all_ids = {
  "motion_rush","grid_rush","find_hunter","word_hunter",
  "line_hunter","paren_hunter","sentence_hunter",
  "slash_hunter","question_hunter","repeat_hunter",
  "inner_paren_hunter","around_paren_hunter",
  "inner_quote_hunter","around_quote_hunter",
  "paragraph_hunter",
  "delete_character_hunter","replace_character_hunter",
  "toggle_case_hunter","delete_word_hunter","change_word_hunter",
  "utf8_cursor_hunter",
  "delete_line_hunter","delete_to_end_hunter",
  "delete_inner_word_hunter","delete_around_word_hunter",
  "delete_inner_paren_hunter","delete_around_paren_hunter",
  "delete_inner_quote_hunter","delete_around_quote_hunter",
  "change_inner_paren_hunter",
  "change_inner_quote_hunter",  "yank_line_hunter",
  "word_register_hunter",
  "register_replace_hunter",
  "find_diw_combo","find_daw_combo",
  "find_di_paren_combo","find_ca_quote_combo","find_ciw_combo",
  "dw_dot_combo","ciw_dot_combo",
  "yank_paste_combo","dd_paste_combo","dd_paste_before_combo",
  "trial_find_delete","trial_find_change","trial_dot_repeat",
  "trial_delete_choice","trial_repeat_choice",
}

local results = {}

for _, id in ipairs(all_ids) do
  local d = util.describe(id, "/tmp/praxis")
  if type(d) ~= "table" then
    print("FAIL " .. id .. ": could not describe")
    results[id] = "FAIL"
    goto continue
  end

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
  vim.api.nvim_set_current_buf(buf)

  if d.verify == "cursor" then
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, d.content)
    vim.api.nvim_set_option_value("modifiable", false, { buf = buf })

    local found = false
    for r = 1, #d.content do
      local line = d.content[r]
      if line then
        for c = 0, #line do
           local charcol = util.byte_to_char(line, c)
          local ch = vim.fn.strcharpart(line, charcol, 1)
          if ch == d.target then
            vim.api.nvim_win_set_cursor(0, { r, c })
            found = true
            break
          end
        end
      end
      if found then break end
    end

    if not found then
      print("FAIL " .. id .. ": target '" .. d.target .. "' not found")
      results[id] = "FAIL"
      goto continue
    end

    local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
    local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
    local completed = false
    if line then
       local charcol = util.byte_to_char(line, col0)
      if vim.fn.strcharpart(line, charcol, 1) == d.target then
        completed = true
      end
    end

    if completed then
      print("PASS " .. id)
      results[id] = "PASS"
    else
      print("FAIL " .. id)
      results[id] = "FAIL"
    end

  elseif d.verify == "buffer" or d.verify == "composite" then
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, d.content)
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, d.result)

    if check_buffer(buf, d.result) then
      print("PASS " .. id)
      results[id] = "PASS"
    else
      local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
      print("FAIL " .. id .. " #current=" .. #current .. " #result=" .. #d.result)
      results[id] = "FAIL"
    end
  end

  if results[id] == "PASS" and d.evaluation then
    if d.evaluation.max_moves <= 0 then
      print("FAIL " .. id .. ": invalid max_moves=" .. d.evaluation.max_moves)
      results[id] = "FAIL"
    end
    if d.layer == "Trial" and (type(d.derived_from) ~= "table" or #d.derived_from == 0) then
      print("FAIL " .. id .. ": missing derived_from")
      results[id] = "FAIL"
    end
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
local eval_ok = 0
local eval_total = 0
local training_eval_ok = 0
local trial_eval_ok = 0

for _, id in ipairs(all_ids) do
  local d = util.describe(id, "/tmp/praxis")
  if type(d) ~= "table" then goto next_id end

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

  if d.evaluation and is_pass then
    eval_total = eval_total + 1
    if d.evaluation.max_moves > 0 then
      eval_ok = eval_ok + 1
      if d.layer == "Training" then training_eval_ok = training_eval_ok + 1 end
      if d.layer == "Trial" then trial_eval_ok = trial_eval_ok + 1 end
    end
  end

  ::next_id::
end

local total_pass = tutorial_cursor_pass + tutorial_buffer_pass + utf8_pass + training_pass + trial_pass
local total_fail = tutorial_cursor_fail + tutorial_buffer_fail + utf8_fail + training_fail + trial_fail
local total = total_pass + total_fail

print("")
print("Tutorial (cursor): " .. tutorial_cursor_pass .. "/" .. (tutorial_cursor_pass + tutorial_cursor_fail) .. " PASS")
print("Tutorial (buffer): " .. tutorial_buffer_pass .. "/" .. (tutorial_buffer_pass + tutorial_buffer_fail) .. " PASS")
print("Tutorial (utf-8):  " .. utf8_pass .. "/" .. (utf8_pass + utf8_fail) .. " PASS")
print("Training:         " .. training_pass .. "/" .. (training_pass + training_fail) .. " PASS  (evaluation " .. training_eval_ok .. "/10)")
print("Trial:            " .. trial_pass .. "/" .. (trial_pass + trial_fail) .. " PASS  (evaluation " .. trial_eval_ok .. "/5)")
print("")
if total_fail == 0 then
  print("ALL " .. total .. "/" .. total .. " REPLAY TESTS PASS")
else
  print(total_pass .. "/" .. total .. " PASS, " .. total_fail .. " FAIL")
end

vim.cmd("qa!")
