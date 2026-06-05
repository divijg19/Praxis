local M = {}

function M.show(opts)
  local lines
  if opts and opts.args and opts.args ~= "" then
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
end

vim.api.nvim_create_user_command('Praxis', M.show, { nargs = '?' })

return M
