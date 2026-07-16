-- Learner journey harness for Praxis (v0.2.7)
-- Validates the EXPERIENCE (navigation, recovery, solving, completion),
-- not just content correctness (that is tools/replay/replay.lua's job).
-- Run via: tools/journey/journey.sh

local self = debug.getinfo(1, "S").source:sub(2)
local ROOT = vim.fn.fnamemodify(self, ":h:h:h")
vim.opt.runtimepath:prepend(ROOT .. "/nvim")
local util = require("praxis.util")
vim.cmd("runtime plugin/praxis.lua")
vim.env.XDG_DATA_HOME = "/tmp/praxis_journey"
vim.env.XDG_CONFIG_HOME = "/tmp/praxis_journey/cfg"

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

-- Solve the CURRENT challenge buffer using real product logic.
local function solve(id)
  local d = util.describe(id)
  if type(d) ~= "table" then return false end
  local buf = vim.api.nvim_get_current_buf()
  if d.verify == "cursor" then
    local trow = 1
    for r = 1, #d.content do
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
vim.fn.system({ "rm", "-rf", "/tmp/praxis_journey" })
vim.fn.system({ "praxis", "reset", "--yes" })
vim.cmd("Praxis")
ok("onboarding", has(snap(), "Welcome to Praxis%."))
press("e")
ok("catalog_open", has(snap(), "Praxis Catalog"))
press("q")
ok("catalog_return", has(snap(), "Welcome to Praxis%."))

-- 2. cursor Tutorial: solve via real keystrokes, then continue
vim.cmd("Praxis motion_rush")
ok("cursor_tutorial_open", has(snap(), "Use h, j, k and l to move"))
ok("cursor_tutorial_solved", solve("motion_rush"))
ok("cursor_tutorial_result", has(snap(), "Complete%."))
press("<CR>")
ok("cursor_tutorial_continue", not has(snap(), "Unknown challenge") and not has(snap(), "executable not found"))

-- 3. buffer Tutorial: solve via real validation path, then escape
vim.cmd("Praxis delete_character_hunter")
ok("buffer_tutorial_open", has(snap(), "Use x to delete the extra letter"))
ok("buffer_tutorial_solved", solve("delete_character_hunter"))
ok("buffer_tutorial_result", has(snap(), "Complete%."))
press("q")
ok("buffer_tutorial_escape", has(snap(), "Progress:"))

-- 4. Training (composite): solve, then escape
vim.cmd("Praxis find_diw_combo")
ok("training_open", has(snap(), "Use f then a letter, followed by diw%."))
ok("training_solved", solve("find_diw_combo"))
ok("training_result", has(snap(), "Complete%."))
press("q")
ok("training_escape", has(snap(), "Progress:"))

-- 5. Trial: solve, then escape
vim.cmd("Praxis trial_find_delete")
ok("trial_open", has(snap(), "Remove the third word%."))
ok("trial_solved", solve("trial_find_delete"))
ok("trial_result", has(snap(), "Complete%."))
press("q")
ok("trial_escape", has(snap(), "Progress:"))

-- 5b. move-counter determinism: cursor navigation must not inflate the move
--     count for composite challenges. In production every cursor move fires
--     CursorMoved; we simulate that with doautocmd so the test exercises the
--     real event path. Navigation alone must leave moves at 0; only the
--     solving edit increments it. Under the old logic, moving the cursor on a
--     composite challenge inflated the count and could trip the move limit.
vim.cmd("Praxis find_diw_combo")
press("j")
vim.cmd("doautocmd CursorMoved")   -- simulate real navigation (inflated moves under old logic)
press("l")
vim.cmd("doautocmd CursorMoved")
press("l")
vim.cmd("doautocmd CursorMoved")
ok("moves_nav_no_inflate", solve("find_diw_combo") and has(snap(), "Moves: 1 /"))
press("q")

-- 6. mid-challenge escape from a fresh challenge (recovery)
vim.cmd("Praxis motion_rush")
press("q")
ok("mid_challenge_escape", has(snap(), "Progress:"))

-- 6b. interrupted session: re-entering Praxis must not orphan buffers
vim.cmd("Praxis motion_rush")
vim.cmd("Praxis")
ok("no_orphan_buffers", count_praxis() == 1)

-- 7. invalid challenge id -> recovery screen, not Lua error
vim.cmd("Praxis does_not_exist")
ok("invalid_id_recovery", has(snap(), "That challenge doesn't exist%."))

-- 8. force completion of all 52, then completion screen
local guard = 0
while guard < 300 do
  local nid = vim.fn.systemlist({ "praxis", "next" })[1] or ""
  if nid == "" then break end
  vim.fn.system({ "praxis", "attempt", nid })
  vim.fn.system({ "praxis", "record", nid, "1", "1" })
  vim.fn.system({ "praxis", "record", nid, "1", "1" })
  vim.fn.system({ "praxis", "record", nid, "1", "1" })
  guard = guard + 1
end
vim.cmd("Praxis")
ok("completion_shown", has(snap(), "Curriculum complete%."))
ok("completion_progress", has(snap(), "Progress: 52/52"))

-- 9. completion review opens a challenge, then escape
press("r")
ok("completion_review", not has(snap(), "Curriculum complete%.") and not has(snap(), "Unknown challenge"))
press("q")
ok("completion_escape", has(snap(), "Curriculum complete%.") or has(snap(), "Progress:"))

-- 10. missing binary -> recovery screen (graceful, no false completion)
vim.fn.system({ "rm", "-rf", "/tmp/praxis_journey" })
vim.fn.system({ "praxis", "reset", "--yes" })
vim.env.PATH = "/usr/local/bin:/usr/bin:/bin"
vim.cmd("Praxis")
ok("missing_binary_recovery", has(snap(), "Praxis isn't installed%."))

vim.cmd("qa!")
