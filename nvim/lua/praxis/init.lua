local M = {}

function M.show(opts)
  local lines
  local is_challenge = opts and opts.args and opts.args ~= ""

  if is_challenge then
    lines = vim.fn.systemlist({ 'praxis', 'challenge', opts.args })
  else
    lines = vim.fn.systemlist({ 'praxis' })
    table.insert(lines, '')
    table.insert(lines, 'CLI Connected')
  end

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.api.nvim_set_option_value('modifiable', false, { buf = buf })
  vim.api.nvim_set_option_value('buftype', 'nofile', { buf = buf })
  vim.api.nvim_buf_set_name(buf, 'Praxis')
  vim.api.nvim_set_current_buf(buf)

  if is_challenge then
    local done = false
    vim.api.nvim_create_autocmd('CursorMoved', {
      buffer = buf,
      callback = function()
        if done then return end
        local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
        local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
        if line and line:sub(col0 + 1):match('^★') then
          done = true
          vim.api.nvim_echo({{"Success"}}, false, {})
        end
      end,
    })
  end
end

vim.api.nvim_create_user_command('Praxis', M.show, { nargs = '?' })

return M
