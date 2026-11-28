local M = {}

-- Captures physical keys while a challenge is active. Since Neovim 0.11 the
-- second on_key argument is the raw typed string (not a boolean); an empty
-- string means synthetic input. Synthetic keys are allowed only in the dev
-- harness. UI mapping events are filtered so hint/retry keys never reach the
-- verifier.
function M.start()
  local ns = vim.api.nvim_create_namespace("praxis_keylog")
  local reveal = vim.api.nvim_replace_termcodes("<Leader>h", true, true, true)
  local retry = vim.api.nvim_replace_termcodes("<Leader>r", true, true, true)
  local keys = {}
  local tainted = false
  vim.on_key(function(key, typed)
    if typed == reveal or typed == retry or key == reveal or key == retry then return end
    if typed == "" or typed == nil then
      if vim.env.PRAXIS_DEV == "1" then
        table.insert(keys, key)
      else
        tainted = true
      end
      return
    end
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
