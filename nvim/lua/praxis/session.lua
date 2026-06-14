local ui = require('praxis.ui')

local M = {}

local session = nil

function M.start()
  if not session then
    session = {
      start_ns      = vim.uv.hrtime(),
      challenges    = 0,
      completions   = 0,
      total_moves   = 0,
      total_time_ms = 0,
    }
  end
  session.challenges = session.challenges + 1
end

function M.record(id, moves, elapsed_ms)
  session.completions = session.completions + 1
  session.total_moves = session.total_moves + moves
  session.total_time_ms = session.total_time_ms + elapsed_ms
  vim.fn.system({ "praxis", "record", id, tostring(moves), tostring(elapsed_ms) })
end

local function format_time(total_s)
  local m = math.floor(total_s / 60)
  local s = total_s % 60
  if m > 0 then return string.format("%dm%ds", m, s) end
  return string.format("%ds", s)
end

function M.show()
  local completions = session and session.completions or 0
  local challenges  = session and session.challenges or 0
  local moves       = session and session.total_moves or 0
  local time_ms     = session and session.total_time_ms or 0
  local total_s     = math.floor(time_ms / 1000)
  local avg_moves   = completions > 0 and math.floor(moves / completions) or 0
  local avg_time_s  = completions > 0 and math.floor(time_ms / completions / 1000) or 0
  local elapsed_s   = session and math.floor((vim.uv.hrtime() - session.start_ns) / 1e9) or 0

  local buf = ui.create_buffer("Praxis Session")
  ui.set_lines(buf, {
    "Session", "",
    "Challenges: " .. challenges,
    "Completions: " .. completions, "",
    "Session Length: " .. format_time(elapsed_s),
    "Practice Time: " .. format_time(total_s), "",
    "Moves: " .. moves, "",
    "Avg Moves: " .. avg_moves,
    "Avg Time: " .. avg_time_s .. "s",
  })
  vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
  vim.api.nvim_set_option_value("bufhidden", "wipe", { buf = buf })
  vim.api.nvim_set_option_value("swapfile", false, { buf = buf })
  vim.api.nvim_set_current_buf(buf)
end

return M
