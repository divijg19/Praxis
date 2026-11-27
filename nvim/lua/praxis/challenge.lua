local M = {}

local util = require("praxis.util")
local keylog = require("praxis.keylog")

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
  local buf = ui.show("Praxis: " .. desc.name, desc.content, editable)
  local start_row = desc.start_row and desc.start_row > 0 and desc.start_row or 3
  local start_col = desc.start_col or 0
  local first_line = desc.content[start_row] or ""
  start_col = math.max(0, math.min(start_col, #first_line))
  vim.api.nvim_win_set_cursor(0, { math.min(start_row, #desc.content), start_col })

  if not util.praxis({ "start", id }) then
    ui.recovery("Praxis couldn't record this start.", {
      "",
      "[Enter] or [q] Back.",
    })
    return
  end

  local state = {
    done            = false,
    edits           = 0,
    start_ns        = vim.uv.hrtime(),
    challenge_lines = desc.content,
    target          = desc.target,
    verify          = desc.verify,
    result_lines    = desc.result,
    challenge_id    = id,
    layer           = desc.layer,
    start_row       = start_row,
    start_col       = start_col,
    par             = desc.par,
    editable        = editable,
    capture         = keylog.start(),
  }

  local reset_challenge

  local function show_failure(buf, lines)
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
    vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  end

  local function render_result(verdict)
    if not verdict.verified then
      show_failure(buf, {
        "Not saved - replay failed.",
        "",
        "Reason: " .. (verdict.reason or "unknown"),
        "No completion was recorded. You can retry.",
        "",
        "[<leader>r] Retry   [q] Back.",
      })
      return
    end

    local display = { "Run verified.", "" }
    if state.layer ~= "Tutorial" then
      table.insert(display, "Keys: " .. (verdict.keys or 0) .. "   (verified by clean replay)")
      table.insert(display, "Par: " .. (verdict.par or 0) .. " canonical keys")
      table.insert(display, "Marks: " .. string.rep("*", verdict.marks or 1) .. string.rep("-", 3 - (verdict.marks or 1)))
      table.insert(display, "")
    end
    table.insert(display, "Completed " .. (verdict.completions or 1) .. " times.")
    table.insert(display, "")
    table.insert(display, "[r] Retry.")
    table.insert(display, "[Enter] Continue.")
    table.insert(display, "[q] Back.")
    show_failure(buf, display)
    vim.keymap.set("n", "r", function()
      reset_challenge()
    end, { buffer = buf, nowait = true, silent = true })
  end

  local function complete_challenge()
    state.done = true
    vim.api.nvim_echo({ { "Verifying run..." } }, false, {})
    local captured = state.capture.stop()
    local payload = {
      keys = captured.keys,
      tainted = captured.tainted,
    }
    local raw = util.send({ "complete", id }, payload)
    if not raw then
      show_failure(buf, {
        "Could not verify this run.",
        "",
        "Praxis did not respond. No completion was recorded.",
        "",
        "[<leader>r] Retry   [q] Back.",
      })
      return
    end
    local ok, verdict = pcall(vim.fn.json_decode, raw)
    if not ok or type(verdict) ~= "table" then
      show_failure(buf, { "Could not read the verification result.", "", "[<leader>r] Retry   [q] Back." })
      return
    end
    render_result(verdict)
  end

  local function check_buffer()
    if not util.check_buffer(buf, state.result_lines) then return end
    vim.api.nvim_echo({ { "Challenge complete." } }, false, {})
    complete_challenge()
  end

  reset_challenge = function()
    pcall(vim.keymap.del, "n", "r", { buffer = buf })
    state.done = false
    state.edits = 0
    state.start_ns = vim.uv.hrtime()
    state.capture = keylog.start()
    vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, state.challenge_lines)
    if not state.editable then
      vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
    end
    local retry_line = state.challenge_lines[state.start_row] or ""
    local retry_col = math.max(0, math.min(state.start_col, #retry_line))
    vim.api.nvim_win_set_cursor(0, { math.min(state.start_row, #state.challenge_lines), retry_col })
    if not util.praxis({ "start", state.challenge_id }) then
      vim.api.nvim_echo({ { "Couldn't record this retry. Press [<leader>r] to try again." } }, false, {})
    end
    vim.api.nvim_echo({ { "[<leader>r] Retry   [Enter] Continue   [q] Back.", "MsgArea" } }, false, {})
  end

  vim.api.nvim_create_autocmd({ "BufWipeout", "BufDelete" }, {
    buffer = buf,
    callback = function()
      state.capture.stop()
    end,
  })

  vim.api.nvim_create_autocmd("CursorMoved", {
    buffer = buf,
    callback = function()
      if state.done then return end
      if state.verify ~= "cursor" then return end
      local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
      if #state.challenge_lines > 1 and row == 1 then return end
      state.edits = state.edits + 1
      local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
      if line then
        local charcol = util.byte_to_char(line, col0)
        if vim.fn.strcharpart(line, charcol, 1) == state.target then
          complete_challenge()
        end
      end
    end,
  })

  vim.api.nvim_create_autocmd("TextChanged", {
    buffer = buf,
    callback = function()
      if state.done then return end
      state.edits = state.edits + 1
      if state.verify == "buffer" or state.verify == "composite" then
        check_buffer()
      end
    end,
  })

  vim.api.nvim_create_autocmd("TextChangedI", {
    buffer = buf,
    callback = function()
      if state.done then return end
      state.edits = state.edits + 1
      if state.verify == "buffer" or state.verify == "composite" then
        check_buffer()
      end
    end,
  })

  vim.keymap.set("n", "<leader>r", function()
    reset_challenge()
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "<CR>", function()
    if state.done then
      local next_lines = util.systemlist({ "next" })
      if not next_lines then
        vim.api.nvim_echo({ { "Couldn't read the next challenge. Returning to the Hub." } }, false, {})
      end
      util.continue(next_lines and next_lines[1])
    end
  end, { buffer = buf, nowait = true, silent = true })

  vim.keymap.set("n", "q", function()
    util.continue("")
  end, { buffer = buf, nowait = true, silent = true })

  vim.api.nvim_echo({ { "[<leader>r] Retry   [Enter] Continue   [q] Back.", "MsgArea" } }, false, {})
end

return M
