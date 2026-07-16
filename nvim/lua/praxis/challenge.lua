local M = {}

local util = require("praxis.util")

-- Completions needed to reach the Practiced tier (mirrors stats.MasteryTier).
local PRACTICED_THRESHOLD = 3

function M.open(id)
  local ui = require('praxis.ui')
  local desc = util.describe(id)
  if not desc then
    ui.recovery("That challenge doesn't exist.", {
      "",
      "[Enter] or [q] Back.",
    })
    return
  end

  local editable = desc.verify == "buffer" or desc.verify == "composite"
  local buf = ui.show("Praxis — " .. desc.name, desc.content, editable)

  if not util.praxis({ "attempt", id }) then
    ui.recovery("Praxis couldn't record this attempt.", {
      "Your progress may not have been saved.",
      "",
      "[Enter] or [q] Back.",
    })
    return
  end

  local state = {
    done            = false,
    moves           = 0,
    start_ns        = vim.uv.hrtime(),
    challenge_lines = desc.content,
    target          = desc.target,
    verify          = desc.verify,
    result_lines    = desc.result,
    challenge_id    = id,
    maxmoves        = desc.evaluation and desc.evaluation.max_moves,
    editable        = editable,
  }

  local function render_result()
    local elapsed_ms = math.floor((vim.uv.hrtime() - state.start_ns) / 1e6)
    if not util.praxis({ "record", id, tostring(state.moves), tostring(elapsed_ms) }) then
      vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
      vim.api.nvim_buf_set_lines(buf, 0, -1, false, {
        "Complete, but your progress wasn't saved.",
        "",
        "Praxis reported an error while recording this challenge.",
        "Your edits are not lost, but stats were not updated.",
        "",
        "[Enter] or [q] Back.",
      })
      vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
      return
    end
    local stats_out = vim.fn.systemlist({ "praxis", "stats", id })
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    local display = {
      "Complete.", "",
    }
    if state.maxmoves then
      table.insert(display, "Moves: " .. state.moves .. " / " .. state.maxmoves)
    else
      table.insert(display, "Moves: " .. state.moves)
    end
    table.insert(display, "Time: "  .. elapsed_ms .. " ms")
    table.insert(display, "")
    local completions = 0
    for _, line in ipairs(stats_out) do
      local n = line:match("^Completions: (%d+)")
      if n then completions = tonumber(n) end
    end
    if completions > 0 then
      table.insert(display, "Completed " .. completions .. " times.")
    end
    if completions < PRACTICED_THRESHOLD then
      table.insert(display, "Practice this " .. (PRACTICED_THRESHOLD - completions) .. " more times to build mastery.")
    end
    table.insert(display, "")
    table.insert(display, "[r] Retry.")
    table.insert(display, "[Enter] Continue.")
    table.insert(display, "[q] Back.")
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, display)
    vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  end

  local function check_buffer()
    if state.maxmoves then
      vim.api.nvim_echo({ { "Moves: " .. state.moves .. " / " .. state.maxmoves } }, false, {})
    end
    local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    if #current ~= #state.result_lines then return end
    for i = 1, #current do
      if current[i] ~= state.result_lines[i] then return end
    end
    if state.verify == "composite" and state.maxmoves and state.moves > state.maxmoves then
      vim.api.nvim_echo({ { "Over the move limit — press [r] to retry." } }, false, {})
      return
    end
    state.done = true
    vim.api.nvim_echo({ { "Challenge complete." } }, false, {})
    render_result()
  end

  local function reset_challenge()
    state.done     = false
    state.moves    = 0
    state.start_ns = vim.uv.hrtime()
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, state.challenge_lines)
    if not state.editable then
      vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
    end
    vim.api.nvim_win_set_cursor(0, { 1, 0 })
    util.praxis({ "attempt", state.challenge_id })
  end

  vim.api.nvim_create_autocmd("CursorMoved", {
    buffer = buf,
    callback = function()
      if state.done then return end
      if state.verify ~= "cursor" then return end
      local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
      if #state.challenge_lines > 1 and row == 1 then return end
      state.moves = state.moves + 1
      local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
      if line then
        local charcol = util.byte_to_char(line, col0)
        if vim.fn.strcharpart(line, charcol, 1) == state.target then
          state.done = true
          vim.api.nvim_echo({ { "Challenge complete." } }, false, {})
          render_result()
        end
      end
    end,
  })

  vim.api.nvim_create_autocmd("TextChanged", {
    buffer = buf,
    callback = function()
      if state.done then return end
      state.moves = state.moves + 1
      if state.verify == "buffer" or state.verify == "composite" then
        check_buffer()
      end
    end,
  })

  vim.api.nvim_create_autocmd("TextChangedI", {
    buffer = buf,
    callback = function()
      if state.done then return end
      if state.verify == "buffer" or state.verify == "composite" then
        check_buffer()
      end
    end,
  })

  vim.keymap.set("n", "r", function()
    reset_challenge()
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "<CR>", function()
    if state.done then
      local next_id = vim.fn.systemlist({ "praxis", "next" })[1]
      util.continue(next_id)
    end
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "q", function()
    util.continue("")
  end, { buffer = buf, nowait = true, silent = true })

  vim.api.nvim_echo({ { "[r] Retry   [Enter] Continue   [q] Back.", "MsgArea" } }, false, {})
end

return M
