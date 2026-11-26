-- Learner journey harness for Praxis (v0.3.10)
-- Validates the EXPERIENCE (navigation, recovery, solving, completion),
-- not just content correctness (that is tools/replay/replay.lua's job).
-- Run via: go run ./tools/verify (PRAXIS_BIN / PRAXIS_VERIFY_DATA may be set by it)

local self = debug.getinfo(1, "S").source:sub(2)
local ROOT = vim.fn.fnamemodify(self, ":h:h:h")
vim.opt.runtimepath:prepend(ROOT .. "/nvim")
local util = require("praxis.util")
local bin = vim.env.PRAXIS_BIN or "praxis"
vim.cmd("runtime plugin/praxis.lua")
local data_home = vim.env.PRAXIS_VERIFY_DATA or "/tmp/praxis_journey"
vim.env.XDG_DATA_HOME = data_home
vim.env.XDG_CONFIG_HOME = data_home .. "/cfg"
if vim.env.PRAXIS_BIN then
  vim.env.PATH = vim.fn.fnamemodify(vim.env.PRAXIS_BIN, ":h") .. ":" .. vim.env.PATH
end

local results = {}
local function ok(name, cond)
  print((cond and "PASS " or "FAIL ") .. name)
  results[name] = cond
end

local function snap()
  return table.concat(vim.api.nvim_buf_get_lines(vim.api.nvim_get_current_buf(), 0, -1, false), "\n")
end
local function has(t, p) return t:match(p) ~= nil end
local function count_praxis()
  local n = 0
  for _, b in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_valid(b) and vim.api.nvim_buf_get_name(b):match("Praxis") then
      n = n + 1
    end
  end
  return n
end
local function press(k)
  vim.cmd("normal " .. k)
end

-- Drive the current challenge buffer through the product's validation
-- autocmds for representative cursor, buffer, composite and Trial modes.
local function solve(id)
  local d = util.describe(id)
  if type(d) ~= "table" then return false end
  local buf = vim.api.nvim_get_current_buf()
  if d.verify == "cursor" then
    local start_row = 1
    if #d.content > 1 then start_row = 2 end
    local trow = start_row
    for r = start_row, #d.content do
      if d.content[r]:find(d.target, 1, true) then trow = r; break end
    end
    if trow > 1 then press(string.rep("j", trow - 1)) end
    local c = d.content[trow]:find(d.target, 1, true)
    press(string.rep("l", c - 1))
    vim.cmd("doautocmd CursorMoved")
    return has(snap(), "Complete%.")
  else
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, d.result)
    vim.cmd("doautocmd TextChanged")
    return true
  end
end

-- 1. reset -> onboarding (fresh learner)
vim.fn.system({ bin, "reset", "--yes" })
vim.cmd("Praxis")
ok("onboarding", has(snap(), "Welcome to Praxis%."))
ok("onboarding_single_action", has(snap(), "%[Enter%] Start Onboarding%.") and not has(snap(), "%[e%] Explore%."))
press("e")
ok("catalog_open", has(snap(), "Praxis Catalog"))
press("q")
ok("catalog_return", has(snap(), "Welcome to Praxis%."))

-- 2. cursor Tutorial: solve via real keystrokes, then continue
vim.cmd("Praxis grid_rush")
ok("cursor_tutorial_open", has(snap(), "Use h, j, k and l to move"))
ok("cursor_tutorial_solved", solve("grid_rush"))
ok("cursor_tutorial_result", has(snap(), "Complete%."))
ok("tutorial_result_no_metrics", not has(snap(), "Moves:"))
press("<CR>")
ok("cursor_tutorial_continue", not has(snap(), "Unknown challenge") and not has(snap(), "executable not found"))

-- 3. buffer Tutorial: solve via real validation path, then escape
vim.cmd("Praxis delete_character_hunter")
ok("buffer_tutorial_open", has(snap(), "Use x to delete the extra letter"))
ok("buffer_tutorial_solved", solve("delete_character_hunter"))
ok("buffer_tutorial_result", has(snap(), "Complete%."))
-- The solve path records the completion through the real praxis binary; verify
-- progress was actually persisted and not silently dropped.
local persisted = vim.fn.systemlist({ bin, "stats", "delete_character_hunter" })
local saved = false
for _, l in ipairs(persisted) do
  if l:match("^Completions: (%d+)") and tonumber(l:match("(%d+)")) >= 1 then saved = true end
end
ok("buffer_tutorial_persisted", saved)
press("q")
ok("buffer_tutorial_escape", has(snap(), "So far:"))
local core_total = "0"
for _, line in ipairs(util.systemlist({ "stats" }) or {}) do
  local n = line:match("^Tutorial: %d+/(%d+) core")
  if n then core_total = n end
end
ok("hub_tutorial_status", has(snap(), "Tutorial: %d+/" .. core_total .. " core%. Finish Core"))

-- 4. Training (composite): solve, then escape
vim.cmd("Praxis find_diw_combo")
ok("training_open", has(snap(), "Use f{char} then diw to delete the word%."))
ok("training_solved", solve("find_diw_combo"))
ok("training_result", has(snap(), "Complete%."))
press("q")
ok("training_escape", has(snap(), "So far:"))

-- 5. Trial: solve, then escape
vim.cmd("Praxis trial_find_delete")
ok("trial_open", has(snap(), "Remove the third word%."))
ok("trial_solved", solve("trial_find_delete"))
ok("trial_result", has(snap(), "Complete%."))
press("q")
ok("trial_escape", has(snap(), "So far:"))

-- 5b. move-counter determinism: cursor navigation must not inflate the move
--     count for composite challenges. Real keystrokes fire CursorMoved, so
--     navigation alone must leave moves at 0; only the solving edit increments
--     it. Under the old logic, moving the cursor on a composite challenge
--     inflated the count and could trip the move limit.
vim.cmd("Praxis find_diw_combo")
press("j")
press("l")
press("l")
ok("moves_nav_no_inflate", solve("find_diw_combo") and has(snap(), "Moves: 1 /"))
press("q")

-- 5c. insert-mode edits must count toward the move budget too.
vim.cmd("Praxis find_diw_combo")
local combo = util.describe("find_diw_combo")
local combo_buf = vim.api.nvim_get_current_buf()
vim.api.nvim_set_option_value("modifiable", true, { buf = combo_buf })
vim.api.nvim_buf_set_lines(combo_buf, 0, -1, false, combo.result)
vim.cmd("doautocmd TextChangedI")
ok("insert_mode_moves_counted", has(snap(), "Moves: 1 /"))
press("q")

-- 6. mid-challenge escape from a fresh challenge (recovery)
vim.cmd("Praxis grid_rush")
press("q")
ok("mid_challenge_escape", has(snap(), "So far:"))

-- 6b. interrupted session: re-entering Praxis must not orphan buffers
vim.cmd("Praxis grid_rush")
vim.cmd("Praxis")
ok("no_orphan_buffers", count_praxis() == 1)

-- 6c. a real file whose name contains "Praxis" must never be force-deleted.
local user_buf = vim.api.nvim_create_buf(false, false)
vim.api.nvim_buf_set_name(user_buf, "/tmp/PraxisNotes.md")
vim.api.nvim_buf_set_lines(user_buf, 0, -1, false, { "unsaved user work" })
vim.api.nvim_set_option_value("modified", true, { buf = user_buf })
vim.api.nvim_set_current_buf(user_buf)
vim.cmd("Praxis")
ok("user_buffer_preserved", vim.api.nvim_buf_is_valid(user_buf)
  and vim.api.nvim_get_option_value("modified", { buf = user_buf }))

-- A clean file named exactly "Praxis" must be kept; our buffer gets a unique name.
local clean_buf = vim.api.nvim_create_buf(false, false)
vim.api.nvim_buf_set_name(clean_buf, "/tmp/praxis-clean/Praxis")
vim.api.nvim_set_current_buf(clean_buf)
vim.cmd("Praxis")
local unique_scratch = false
for _, b in ipairs(vim.api.nvim_list_bufs()) do
  if vim.fn.fnamemodify(vim.api.nvim_buf_get_name(b), ":t") == "Praxis 2" then
    unique_scratch = true
  end
end
ok("clean_user_buffer_preserved", vim.api.nvim_buf_is_valid(clean_buf) and unique_scratch)
pcall(vim.api.nvim_buf_delete, clean_buf, { force = true })

-- 7. invalid challenge id -> recovery screen, not Lua error
vim.cmd("Praxis does_not_exist")
ok("invalid_id_recovery", has(snap(), "That challenge doesn't exist%."))

-- 7b. completing Core produces the Training transition message.
local core_ids = {
  "grid_rush", "find_hunter", "slash_hunter", "word_hunter", "line_hunter",
  "delete_character_hunter", "delete_word_hunter", "replace_character_hunter",
  "change_word_hunter", "delete_inner_word_hunter", "yank_line_hunter",
}
for _, id in ipairs(core_ids) do
  vim.fn.system({ bin, "attempt", id })
  for _ = 1, 3 do vim.fn.system({ bin, "record", id, "1", "1" }) end
end
vim.cmd("Praxis")
ok("training_transition", has(snap(), "Onboarding complete%. Training uses a move limit"))

-- 8. force completion of all 48, then completion screen
local guard = 0
while guard < 300 do
  local nid = vim.fn.systemlist({ bin, "next" })[1] or ""
  if nid == "" then break end
  vim.fn.system({ bin, "attempt", nid })
  vim.fn.system({ bin, "record", nid, "1", "1" })
  vim.fn.system({ bin, "record", nid, "1", "1" })
  vim.fn.system({ bin, "record", nid, "1", "1" })
  guard = guard + 1
end
vim.cmd("Praxis")
ok("completion_shown", has(snap(), "Curriculum complete%."))
local completed_total = "0"
for _, line in ipairs(util.systemlist({ "stats" }) or {}) do
  local n = line:match("^Challenges Completed: (%d+)/(%d+)")
  if n then completed_total = n end
end
ok("completion_progress", has(snap(), "So far: " .. completed_total .. " challenges completed%."))

-- 9. completion review opens a challenge, then escape
press("r")
ok("completion_review", not has(snap(), "Curriculum complete%.") and not has(snap(), "Unknown challenge"))
press("q")
ok("completion_escape", has(snap(), "Curriculum complete%.") or has(snap(), "So far:"))

-- 9b. corrupt progress -> recovery, never a false "complete" Hub.
vim.fn.writefile({ "{broken" }, data_home .. "/praxis/stats.json")
vim.cmd("Praxis")
ok("corrupt_progress_recovery", has(snap(), "couldn't read your progress%."))
vim.fn.system({ bin, "reset", "--yes" })

-- 9c. unsupported version -> recovery, never a false "complete" Hub.
vim.fn.writefile({ '{"version":99,"entries":{}}' }, data_home .. "/praxis/stats.json")
vim.cmd("Praxis")
ok("unknown_version_recovery", has(snap(), "couldn't read your progress%."))
vim.fn.system({ bin, "reset", "--yes" })

-- 10. missing binary -> recovery screen (graceful, no false completion)
vim.fn.system({ bin, "reset", "--yes" })
vim.env.PATH = "/usr/local/bin:/usr/bin:/bin"
vim.cmd("Praxis")
ok("missing_binary_recovery", has(snap(), "Praxis isn't installed%."))

local failed = false
for _, passed in pairs(results) do
  if not passed then
    failed = true
    break
  end
end
if failed then
  vim.cmd("cquit 1")
else
  vim.cmd("qa!")
end
