local M = {}

function M.open()
  local ui = require('praxis.ui')
  local next_id = vim.fn.systemlist({ "praxis", "next" })[1] or ""

  local buf = ui.show("Praxis", {
    "── Praxis ──",
    "",
    "Welcome to Praxis.",
    "",
    "Learn Vim by solving short editing challenges.",
    "",
    "Progress through Tutorials, Training, and Trials.",
    "",
    "If you're new, press [s] to begin.",
    "",
    "[s] Start.",
    "    Follow the guided curriculum from the beginning.",
    "",
    "[e] Explore.",
    "    Browse all challenges.",
    "",
    "[h] About.",
    "    Learn how the curriculum works.",
    "",
    "[p] View progress.",
    "    See completed challenges and mastery.",
  }, false)

  vim.keymap.set("n", "s", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    if next_id ~= "" then
      require('praxis.challenge').open(next_id)
    else
      require('praxis.hub').open()
    end
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "e", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local catalog = vim.fn.systemlist({ "praxis", "catalog" })
    local cat_lines = { "── Praxis Catalog ──", "", "[q] Back.", "" }
    for _, c in ipairs(catalog) do
      table.insert(cat_lines, "  " .. c)
    end
    local cat_buf = ui.show("Praxis Catalog", cat_lines, false)
    vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = cat_buf })
    local function back()
      pcall(vim.api.nvim_buf_delete, cat_buf, { force = true })
      M.open()
    end
    vim.keymap.set("n", "q", back, { buffer = cat_buf, nowait = true, silent = true })
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
      "not paused, done.",
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
      "You are done when every challenge shows as complete.",
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
