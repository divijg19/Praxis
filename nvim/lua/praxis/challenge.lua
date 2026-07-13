local M = {}

function M.open(id)
  local ui = require('praxis.ui')
  local desc_raw = vim.fn.systemlist({ "praxis", "describe", id })
  local ok, desc = pcall(vim.fn.json_decode, table.concat(desc_raw, ""))
  if not ok or type(desc) ~= "table" then
    require("praxis.ui").recovery("Unknown challenge: " .. id, {
      "",
      "Press [Enter] or [q] to return.",
    })
    return
  end

  local buf = ui.create_buffer("Praxis")
  ui.set_lines(buf, desc.content)
  ui.set_modifiable(buf, desc.verify == "buffer" or desc.verify == "composite")
  vim.api.nvim_set_current_buf(buf)

  vim.fn.system({ "praxis", "attempt", id })

  local function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
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
  }

  local function render_result()
    local elapsed_ms = math.floor((vim.uv.hrtime() - state.start_ns) / 1e6)
    vim.fn.system({ "praxis", "record", id, tostring(state.moves), tostring(elapsed_ms) })
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
    table.insert(display, "Time: "  .. elapsed_ms .. "ms")
    table.insert(display, "")
    for _, line in ipairs(stats_out) do
      table.insert(display, line)
    end
    table.insert(display, "")
    table.insert(display, "[r] Retry")
    table.insert(display, "[Enter] Continue")
    table.insert(display, "[q] Quit")
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, display)
    vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  end

  local function check_buffer()
    local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    if #current ~= #state.result_lines then return end
    for i = 1, #current do
      if current[i] ~= state.result_lines[i] then return end
    end
    if state.verify == "composite" and state.maxmoves and state.moves > state.maxmoves then
        vim.api.nvim_echo({ { "Too many moves! Press [r] to retry.", "WarningMsg" } }, false, {})
      return
    end
    state.done = true
    vim.api.nvim_echo({ { "Success" } }, false, {})
    render_result()
  end

  local function reset_challenge()
    state.done     = false
    state.moves    = 0
    state.start_ns = vim.uv.hrtime()
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, state.challenge_lines)
    if state.verify ~= "buffer" and state.verify ~= "composite" then
      vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
    end
    vim.api.nvim_win_set_cursor(0, { 1, 0 })
    vim.fn.system({ "praxis", "attempt", state.challenge_id })
  end

  vim.api.nvim_create_autocmd("CursorMoved", {
    buffer = buf,
    callback = function()
      if state.done then return end
      state.moves = state.moves + 1
      if state.verify == "cursor" then
        local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
        local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
        if line then
          local charcol = byte_to_char(line, col0)
          if vim.fn.strcharpart(line, charcol, 1) == state.target then
            state.done = true
            vim.api.nvim_echo({ { "Success" } }, false, {})
            render_result()
          end
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
      pcall(vim.api.nvim_buf_delete, buf, { force = true })
      if next_id and next_id ~= "" then
        vim.cmd("Praxis " .. next_id)
      else
        vim.cmd("Praxis")
      end
    end
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "q", function()
    pcall(vim.api.nvim_buf_delete, buf, { force = true })
    vim.cmd("Praxis")
  end, { buffer = buf, nowait = true, silent = true })
end

return M
