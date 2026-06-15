local ui = require('praxis.ui')

local M = {}

function M.open(id)
  local lines = vim.fn.systemlist({ "praxis", "challenge", id })
  local verify = vim.fn.systemlist({ "praxis", "verify", id })[1] or ""
  local target = vim.fn.systemlist({ "praxis", "target", id })[1]
  local result = vim.fn.systemlist({ "praxis", "result", id }) or {}

  local buf = ui.create_buffer("Praxis")
  ui.set_lines(buf, lines)
  ui.set_modifiable(buf, verify == "buffer")
  vim.api.nvim_set_current_buf(buf)

  local session = require('praxis.session')
  session.start()
  vim.fn.system({ "praxis", "attempt", id })

  local function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
  end

  local state = {
    done            = false,
    moves           = 0,
    start_ns        = vim.uv.hrtime(),
    challenge_lines = lines,
    target          = target,
    verify          = verify,
    result_lines    = result,
    challenge_id    = id,
  }

  local function check_buffer()
    local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    if #current ~= #result then return end
    for i = 1, #current do
      if current[i] ~= result[i] then return end
    end
    state.done = true
    vim.api.nvim_echo({ { "Success" } }, false, {})
    render_result()
  end

  local function render_result()
    local elapsed_ms = math.floor((vim.uv.hrtime() - state.start_ns) / 1e6)
    session.record(id, state.moves, elapsed_ms)
    local stats_out = vim.fn.systemlist({ "praxis", "stats", id })
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    local display = {
      "Complete.", "",
      "Moves: " .. state.moves,
      "Time: "  .. elapsed_ms .. "ms", "",
    }
    for _, line in ipairs(stats_out) do
      table.insert(display, line)
    end
    table.insert(display, "")
    table.insert(display, "[r] Replay")
    table.insert(display, "[Enter] Continue Journey")
    table.insert(display, "[q] Quit")
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, display)
    vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  end

  local function reset_challenge()
    state.done     = false
    state.moves    = 0
    state.start_ns = vim.uv.hrtime()
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, state.challenge_lines)
    if state.verify ~= "buffer" then
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
      if state.verify == "buffer" then
        check_buffer()
      end
    end,
  })

  vim.api.nvim_create_autocmd("TextChangedI", {
    buffer = buf,
    callback = function()
      if state.done then return end
      if state.verify == "buffer" then
        check_buffer()
      end
    end,
  })

  vim.keymap.set("n", "r", function()
    if state.done then reset_challenge() end
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
    if state.done then
      pcall(vim.api.nvim_buf_delete, buf, { force = true })
      vim.cmd("Praxis")
    end
  end, { buffer = buf, nowait = true, silent = true })
end

return M
