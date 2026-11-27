local M = {}

local util = require("praxis.util")

local function start(next_id, buf)
  pcall(vim.api.nvim_buf_delete, buf, { force = true })
  if next_id and next_id ~= "" then
    require('praxis.challenge').open(next_id)
  else
    require('praxis.hub').open()
  end
end

function M.open()
  local ui = require('praxis.ui')
  local next_lines = util.systemlist({ "next" })
  local next_id = next_lines and next_lines[1] or ""

  local advanced = false
  local stats_lines = util.systemlist({ "stats" })
  if stats_lines then
    for _, line in ipairs(stats_lines) do
      local n = line:match("^Challenges Completed: (%d+)/")
      if n and tonumber(n) > 0 then advanced = true end
    end
  end

  local lines = {
    "── Praxis ──",
    "",
    "Welcome to Praxis.",
    "",
    "Learn Vim by solving short editing challenges.",
    "",
    "Progress through Tutorials, Training, and Trials.",
    "",
  }
  if advanced then
    table.insert(lines, "[Enter] Start.")
    table.insert(lines, "    Continue the guided curriculum.")
    table.insert(lines, "")
    table.insert(lines, "[e] Explore.   Browse all challenges.")
    table.insert(lines, "[h] About.     Learn how the curriculum works.")
    table.insert(lines, "[p] Progress.  See completed challenges.")
  else
    table.insert(lines, "[Enter] Start Onboarding.")
  end

  local buf = ui.show("Praxis", lines, false)

  vim.keymap.set("n", "<CR>", function()
    start(next_id, buf)
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "s", function()
    start(next_id, buf)
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "e", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local catalog = util.systemlist({ "catalog" }) or {}
    local cat_lines = { "── Praxis Catalog ──", "", "[q] Back.", "" }
    for _, c in ipairs(catalog) do
      table.insert(cat_lines, "  " .. c)
    end
    local cat_buf = ui.show("Praxis Catalog", cat_lines, false)
    vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = cat_buf })
    vim.keymap.set("n", "q", function()
      pcall(vim.api.nvim_buf_delete, cat_buf, { force = true })
      M.open()
    end, { buffer = cat_buf, nowait = true, silent = true })
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "h", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local about_lines = {
      "── Praxis About ──",
      "",
      "Praxis is deliberate practice, not a course.",
      "",
      "Tutorial is onboarding: it teaches just enough that you can",
      "improve on your own. Finish the core, and Tutorial is done,",
      "not paused.",
      "",
      "  • Tutorial  : teaches. Core is required; the rest are Additional",
      "               Lessons you can browse anytime, not a new mode.",
      "  • Training  : never teaches mechanics; it refines fluency.",
      "  • Trials    : integrates. Solve real problems your own way.",
      "",
      "Tutorial Complete means Core is done and you are independent.",
      "Training and Trials never truly finish. Fluency is a process.",
      "Curriculum Complete (every exercise seen) is an achievement,",
      "not the goal.",
      "",
      "After Tutorial, Training and Trials are the product. You choose",
      "what to practice, and repetition builds lasting mastery.",
      "",
      "You can always return to where you left off with :Praxis.",
      "",
      "[q] Back.",
    }
    local about_buf = ui.show("Praxis About", about_lines, false)
    vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = about_buf })
    vim.keymap.set("n", "q", function()
      pcall(vim.api.nvim_buf_delete, about_buf, { force = true })
      M.open()
    end, { buffer = about_buf, nowait = true, silent = true })
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "p", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    require('praxis.hub').open()
  end, { buffer = buf, nowait = true, silent = true })
end

return M
