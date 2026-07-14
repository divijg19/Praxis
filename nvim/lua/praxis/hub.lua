local M = {}

local function describe_for_id(id)
  local raw = vim.fn.systemlist({ "praxis", "describe", id })
  local desc = vim.fn.json_decode(table.concat(raw, ""))
  if type(desc) == "table" then
    return desc.name or "", desc.layer or "", desc.stage or ""
  end
  return "", "", ""
end

function M.open()
  local ui = require('praxis.ui')
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
      if line:match("^Most mastered:") then
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

  local name, layer, stage = "", "", ""
  if next_id ~= "" then
    name, layer, stage = describe_for_id(next_id)
  end

  if layer ~= "" then
    table.insert(display, "  Current: " .. layer .. " — " .. stage)
  else
    table.insert(display, "  Current: Complete")
  end
  table.insert(display, "  Progress: " .. (completed or "0") .. "/" .. (total or tostring(#vim.fn.systemlist({ "praxis", "catalog" }))))
  table.insert(display, "")

  table.insert(display, "  Direction:")
  if next_id ~= "" then
    table.insert(display, "    Next: " .. name .. " — " .. stage)
  else
    table.insert(display, "    Complete")
  end
  if review_challenge and review_challenge ~= "" then
    local review_name, _, review_stage = describe_for_id(review_challenge)
    table.insert(display, "    Review: " .. review_name .. " — " .. review_stage)
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

  if next_id ~= "" then
    if review_challenge and review_challenge ~= "" then
      table.insert(display, "  [Enter] Continue, or [r] Review.")
    else
      table.insert(display, "  [Enter] Continue.")
    end
    table.insert(display, "  [q] Back.")
  else
    table.insert(display, "  Curriculum complete.")
    if review_challenge and review_challenge ~= "" then
      table.insert(display, "  [r] Review.")
    end
    table.insert(display, "  [q] Back.")
  end

  local buf = ui.show("Praxis", display, false)

  local function open_target(id)
    if not id or id == "" then return end
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    require('praxis.challenge').open(id)
  end

  vim.keymap.set("n", "<CR>", function()
    if next_id and next_id ~= "" then
      open_target(next_id)
    else
      open_target(review_challenge)
    end
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "r", function()
    open_target(review_challenge)
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "q", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local rb = vim.g.praxis_return_buf
    if rb and vim.api.nvim_buf_is_valid(rb) then
      vim.api.nvim_set_current_buf(rb)
    end
  end, { buffer = buf, nowait = true, silent = true })
end

return M
