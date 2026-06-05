local M = {}

function M.show(opts)
	local lines
	local is_challenge = opts and opts.args and opts.args ~= ""

	if is_challenge then
		lines = vim.fn.systemlist({ "praxis", "challenge", opts.args })
	else
		lines = vim.fn.systemlist({ "praxis" })
		table.insert(lines, "")
		table.insert(lines, "CLI Connected")
	end

	local buf = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
	vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
	vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
	vim.api.nvim_buf_set_name(buf, "Praxis")
	vim.api.nvim_set_current_buf(buf)

	if is_challenge then
		local done = false
		local moves = 0
		local start_ns = vim.uv.hrtime()

		vim.api.nvim_create_autocmd("CursorMoved", {
			buffer = buf,
			callback = function()
				if done then
					return
				end
				moves = moves + 1
				local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
				local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
				if line and vim.fn.strcharpart(line, col0, 1) == "★" then
					done = true
					local elapsed_ms = math.floor((vim.uv.hrtime() - start_ns) / 1e6)
					local result = {
						"Success",
						"",
						"Moves: " .. moves,
						"Time: " .. elapsed_ms .. "ms",
					}
					vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
					vim.api.nvim_buf_set_lines(buf, 0, -1, false, result)
					vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
					vim.api.nvim_echo({ { "Success" } }, false, {})
				end
			end,
		})
	end
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })

return M
