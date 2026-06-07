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
		local target = vim.fn.systemlist({ "praxis", "target", opts.args })[1]
		local state = {
			done            = false,
			moves           = 0,
			start_ns        = vim.uv.hrtime(),
			challenge_lines = lines,
			target          = target,
		}

		local function render_result()
			local elapsed_ms = math.floor((vim.uv.hrtime() - state.start_ns) / 1e6)
			vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
			vim.api.nvim_buf_set_lines(buf, 0, -1, false, {
				"Success", "",
				"Moves: " .. state.moves,
				"Time: "  .. elapsed_ms .. "ms", "",
				"Press r to replay",
			})
			vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
		end

		local function reset_challenge()
			state.done     = false
			state.moves    = 0
			state.start_ns = vim.uv.hrtime()
			vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
			vim.api.nvim_buf_set_lines(buf, 0, -1, false, state.challenge_lines)
			vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
			vim.api.nvim_win_set_cursor(0, { 1, 0 })
		end

		vim.api.nvim_create_autocmd("CursorMoved", {
			buffer = buf,
			callback = function()
				if state.done then return end
				state.moves = state.moves + 1
				local row, col0 = unpack(vim.api.nvim_win_get_cursor(0))
				local line = vim.api.nvim_buf_get_lines(buf, row - 1, row, false)[1]
				if line and vim.fn.strcharpart(line, col0, 1) == state.target then
					state.done = true
					vim.api.nvim_echo({ { "Success" } }, false, {})
					render_result()
				end
			end,
		})

		vim.keymap.set("n", "r", function()
			if state.done then reset_challenge() end
		end, { buffer = buf, nowait = true, silent = true })
	end
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })

return M
