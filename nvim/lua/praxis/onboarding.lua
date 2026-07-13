local M = {}

function M.open()
  local ui = require('praxis.ui')
  local buf = ui.create_buffer("Praxis")
  local next_id = vim.fn.systemlist({ "praxis", "next" })[1] or ""

  ui.set_lines(buf, {
    "Welcome to Praxis.",
    "",
    "Learn Vim by solving short editing challenges.",
    "",
    "Progress from single motions → combined skills → practical challenges.",
    "",
    "If you're new, press s to begin.",
    "",
    "[s] Start Learning",
    "    Follow the guided curriculum from the beginning.",
    "",
    "[e] Explore Curriculum",
    "    Browse all challenges.",
    "",
    "[h] About Praxis",
    "    Learn how the curriculum works.",
    "",
    "[p] View Progress",
    "    See completed challenges and mastery.",
  })
  ui.set_modifiable(buf, false)
  vim.api.nvim_set_current_buf(buf)

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
    local cat_buf = ui.create_buffer("Praxis Catalog")
    local catalog = vim.fn.systemlist({ "praxis", "catalog" })
    local cat_lines = { "── Praxis Catalog ──", "", "[q] Back.", "" }
    for _, c in ipairs(catalog) do
      table.insert(cat_lines, "  " .. c)
    end
    ui.set_lines(cat_buf, cat_lines)
    ui.set_modifiable(cat_buf, false)
    vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = cat_buf })
    vim.api.nvim_set_current_buf(cat_buf)
    local function back()
      pcall(vim.api.nvim_buf_delete, cat_buf, { force = true })
      M.open()
    end
    vim.keymap.set("n", "q", back, { buffer = cat_buf, nowait = true, silent = true })
    vim.keymap.set("n", "<CR>", back, { buffer = cat_buf, nowait = true, silent = true })
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "h", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    local about_buf = ui.create_buffer("Praxis About")
    local about_lines = {
      "Praxis — Mastery through practice.",
      "",
      "You learn Vim by solving short editing challenges:",
      "  • Tutorials  — learn one technique, with a hint.",
      "  • Training   — combine techniques, within a move budget.",
      "  • Trials     — apply it yourself, with no hints.",
      "",
      "Each challenge you finish builds mastery. You can always",
      "return to where you left off with :Praxis.",
      "",
      "Each challenge is solved a few times to build lasting mastery.",
      "",
      "You are done when every challenge shows as complete.",
      "",
      "[q] Back.",
    }
    ui.set_lines(about_buf, about_lines)
    ui.set_modifiable(about_buf, false)
    vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = about_buf })
    vim.api.nvim_set_current_buf(about_buf)
    vim.keymap.set("n", "q", function()
      pcall(vim.api.nvim_buf_delete, buf, { force = true })
      M.open()
    end, { buffer = about_buf, nowait = true, silent = true })
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "p", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    require('praxis.hub').open()
  end, { buffer = buf, nowait = true, silent = true })
end

return M
