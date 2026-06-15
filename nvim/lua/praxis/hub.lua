local ui = require('praxis.ui')
local challenge = require('praxis.challenge')

local M = {}

function M.open()
  local next_id = vim.fn.systemlist({ "praxis", "next" })[1] or ""
  local stats_lines = vim.fn.systemlist({ "praxis", "stats" })

  local display = {
    "── Praxis ──────────────────────────────────────",
    "",
  }

  local completed, total
  local in_mastery = false
  local mastery_unseen, mastery_learning, mastery_practiced, mastery_experienced
  local review_challenge

  for _, line in ipairs(stats_lines) do
    local comp, tot = line:match("Challenges Completed: (%d+)/(%d+)")
    if comp then
      completed, total = comp, tot
    end

    if line:match("^Mastery:") then
      in_mastery = true
    elseif in_mastery then
      local n = line:match("^  Unseen: (%d+)")
      if n then mastery_unseen = n end
      n = line:match("^  Learning: (%d+)")
      if n then mastery_learning = n end
      n = line:match("^  Practiced: (%d+)")
      if n then mastery_practiced = n end
      n = line:match("^  Experienced: (%d+)")
      if n then mastery_experienced = n end
      if line:match("^Highest Tier:") then
        in_mastery = false
      end
    end
  end

  for i, line in ipairs(stats_lines) do
    if line:match("^Recommended Review:") then
      local rc = stats_lines[i + 1]
      if rc then
        review_challenge = rc:match("^  (.+)")
      end
      break
    end
  end

  local stage = ""
  if next_id ~= "" then
    stage = vim.fn.systemlist({ "praxis", "stage", next_id })[1] or ""
  end

  if stage ~= "" then
    table.insert(display, "  Location: " .. stage)
  else
    table.insert(display, "  Location: Complete")
  end
  table.insert(display, "  Progress: " .. (completed or "0") .. "/" .. (total or "41"))
  table.insert(display, "")

  table.insert(display, "  Direction:")
  if next_id ~= "" then
    table.insert(display, "    Next: " .. next_id .. " — " .. stage)
  else
    table.insert(display, "    Complete")
  end
  if review_challenge and review_challenge ~= "" then
    local review_stage = vim.fn.systemlist({ "praxis", "stage", review_challenge })[1] or ""
    table.insert(display, "    Review: " .. review_challenge .. " — " .. review_stage)
  end
  table.insert(display, "")

  table.insert(display, "  Mastery:")
  local parts = {}
  if mastery_unseen then table.insert(parts, "Unseen: " .. mastery_unseen) end
  if mastery_learning then table.insert(parts, "Learning: " .. mastery_learning) end
  if mastery_practiced then table.insert(parts, "Practiced: " .. mastery_practiced) end
  if mastery_experienced then table.insert(parts, "Experienced: " .. mastery_experienced) end
  if #parts > 0 then
    table.insert(display, "    " .. table.concat(parts, "   "))
  end
  table.insert(display, "")

  table.insert(display, "  Press Enter to continue.")

  local buf = ui.create_buffer("Praxis")
  ui.set_lines(buf, display)
  ui.set_modifiable(buf, false)
  vim.api.nvim_set_current_buf(buf)

  vim.keymap.set("n", "<CR>", function()
    if next_id and next_id ~= "" then
      pcall(vim.api.nvim_buf_delete, buf, { force = true })
      challenge.open(next_id)
    end
  end, { buffer = buf, nowait = true, silent = true })
end

return M
