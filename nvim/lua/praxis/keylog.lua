local M = {}

-- Captures physical keys while a challenge is active. Keys reported as
-- non-typed (synthetic/autocmd input) taint the log so it cannot earn a mark.
function M.start()
  local ns = vim.api.nvim_create_namespace("praxis_keylog")
  local keys = {}
  local tainted = false
  vim.on_key(function(key, typed)
    if not typed and vim.env.PRAXIS_DEV ~= "1" then tainted = true end
    table.insert(keys, key)
  end, ns)
  return {
    stop = function()
      vim.on_key(nil, ns)
      return { keys = keys, tainted = tainted }
    end,
  }
end

return M
