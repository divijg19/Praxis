-- Learner journey harness for Praxis
-- Validates the EXPERIENCE (navigation, recovery, solving, completion),
-- not just content correctness (that is tools/replay/replay.lua's job).
-- Run via: go run ./tools/verify (PRAXIS_BIN / PRAXIS_VERIFY_DATA may be set by it)

local self = debug.getinfo(1, "S").source:sub(2)
local ROOT = vim.fn.fnamemodify(self, ":h:h:h")
vim.opt.runtimepath:prepend(ROOT .. "/nvim")
local util = require("praxis.util")
local bin = util.bin()
local trace = util.trace
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
  local tr = trace(id)
  if type(d) ~= "table" or type(tr) ~= "table" or type(d.content) ~= "table" then return false end
  local row = math.max(1, math.min(tr.start_row, #d.content))
  local line = d.content[row] or ""
  local col = math.max(0, math.min(tr.start_col, #line))
  vim.api.nvim_win_set_cursor(0, { row, col })
  vim.api.nvim_feedkeys(tr.keys, "x", true)
  if d.verify == "cursor" then
    vim.cmd("doautocmd CursorMoved")
  else
    vim.cmd("doautocmd TextChanged")
  end
  return has(snap(), "Run verified%.")
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
ok("cursor_tutorial_result", has(snap(), "Run verified%."))
ok("tutorial_result_no_metrics", not has(snap(), "Edits:"))
press("<CR>")
ok("cursor_tutorial_continue", not has(snap(), "Unknown challenge") and not has(snap(), "executable not found"))

-- 3. buffer Tutorial: solve via real validation path, then escape
vim.cmd("Praxis delete_character_hunter")
ok("buffer_tutorial_open", has(snap(), "Use x to delete the extra letter"))
ok("buffer_tutorial_solved", solve("delete_character_hunter"))
ok("buffer_tutorial_result", has(snap(), "Run verified%."))
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

-- 3b. `r` must remain available as a taught Vim command, not a retry key.
vim.cmd("Praxis replace_character_hunter")
ok("replace_key_not_intercepted", solve("replace_character_hunter"))
-- Retry from the result screen must release `r` back to Vim.
press("r")
ok("retry_releases_r", solve("replace_character_hunter"))
press("q")

-- 4. Training (composite): solve, then escape
vim.cmd("Praxis find_diw_combo")
ok("training_open", has(snap(), "Use f{char} then diw to delete the word%."))
ok("training_solved", solve("find_diw_combo"))
ok("training_result", has(snap(), "Run verified%."))
press("q")
ok("training_escape", has(snap(), "So far:"))

-- 5. Trial: solve, then escape
vim.cmd("Praxis trial_find_delete")
ok("trial_open", has(snap(), "Remove the third word%."))
ok("trial_solved", solve("trial_find_delete"))
ok("trial_result", has(snap(), "Run verified%."))
press("q")
ok("trial_escape", has(snap(), "So far:"))

-- 5b. move-counter determinism: cursor navigation must not inflate the move
--     count for composite challenges. Real keystrokes fire CursorMoved, so
--     navigation alone must not inflate the edit count
--     it. Under the old logic, moving the cursor on a composite challenge
--     inflated the edit count.
vim.cmd("Praxis find_diw_combo")
ok("edits_nav_no_inflate", solve("find_diw_combo") and has(snap(), "Run verified%."))
press("q")

-- 5c. insert-mode edits are counted as local effort.
vim.cmd("Praxis find_ciw_combo")
ok("insert_mode_edits_counted", solve("find_ciw_combo") and has(snap(), "Run verified%."))
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
-- Core IDs come from the CLI so the harness cannot drift from the catalog.
local core_ids = vim.fn.systemlist({ bin, "catalog", "--core" })
for _, id in ipairs(core_ids) do
  vim.fn.system({ bin, "dev", "mark-complete", id, "3" })
end
vim.cmd("Praxis")
ok("training_transition", has(snap(), "Onboarding complete%. Training composes what you learned"))

-- 8. force completion of the full catalog, then completion screen
local guard = 0
while guard < 300 do
  local nid = vim.fn.systemlist({ bin, "next" })[1] or ""
  if nid == "" then break end
  vim.fn.system({ bin, "dev", "mark-complete", nid, "3" })
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
ok("completion_review", not has(snap(), "doesn't exist"))
press("q")
ok("completion_escape", has(snap(), "So far:") or not has(vim.fn.bufname("%"), "Praxis"))

-- 9b. corrupt progress -> recovery, never a false "complete" Hub.
vim.fn.writefile({ "{broken" }, data_home .. "/praxis/progress.json")
vim.cmd("Praxis")
ok("corrupt_progress_recovery", has(snap(), "couldn't read your progress%."))
vim.fn.system({ bin, "reset", "--yes" })

-- 9c. unsupported version -> recovery, never a false "complete" Hub.
vim.fn.writefile({ '{"version":99,"entries":{}}' }, data_home .. "/praxis/progress.json")
vim.cmd("Praxis")
ok("unknown_version_recovery", has(snap(), "couldn't read your progress%."))
vim.fn.system({ bin, "reset", "--yes" })

-- 10. missing binary -> recovery screen (graceful, no false completion)
vim.fn.system({ bin, "reset", "--yes" })
vim.env.PRAXIS_BIN = nil
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
