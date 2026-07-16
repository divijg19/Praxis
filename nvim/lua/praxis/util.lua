local M = {}

-- Describe a challenge by id, returning the decoded table or nil on failure.
-- `bin` overrides the praxis binary; defaults to "praxis" on $PATH. The replay
-- harness passes the freshly built "/tmp/praxis" so it exercises the built
-- binary rather than whatever is installed.
function M.describe(id, bin)
  bin = bin or "praxis"
  local raw = vim.fn.systemlist({ bin, "describe", id })
  local ok, desc = pcall(vim.fn.json_decode, table.concat(raw, ""))
  if not ok or type(desc) ~= "table" then return nil end
  return desc
end

-- Convert a byte column index to a character column index within a line.
function M.byte_to_char(line, bytecol)
  return vim.fn.strchars(string.sub(line, 1, bytecol))
end

-- Close the current Praxis buffer and open `id` (or the hub when id is empty)
-- through the Praxis command, so navigation state resets cleanly. This is the
-- single open-next/back path shared by the hub and the challenge screen.
function M.continue(id)
  pcall(vim.api.nvim_buf_delete, vim.api.nvim_get_current_buf(), { force = true })
  if id and id ~= "" then
    vim.cmd("Praxis " .. id)
  else
    vim.cmd("Praxis")
  end
end

-- Return to the buffer that was active before Praxis opened (the "return
-- buffer"), or to the hub if none. Single owner of the back-navigation path
-- used by the recovery screen and the hub.
function M.return_to_previous()
  pcall(vim.api.nvim_buf_delete, vim.api.nvim_get_current_buf(), { force = true })
  local rb = vim.g.praxis_return_buf
  if rb and vim.api.nvim_buf_is_valid(rb) then
    vim.api.nvim_set_current_buf(rb)
  else
    vim.cmd("Praxis")
  end
end

return M
