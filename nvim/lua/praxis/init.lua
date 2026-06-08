local M = {}

function M.show(opts)
	local lines
	local is_challenge = opts and opts.args and opts.args ~= ""

	local verify = ""
	if is_challenge then
		lines = vim.fn.systemlist({ "praxis", "challenge", opts.args })
		verify = vim.fn.systemlist({ "praxis", "verify", opts.args })[1] or ""
	else
		lines = vim.fn.systemlist({ "praxis" })
		table.insert(lines, "")
		table.insert(lines, "CLI Connected")
	end

	local buf = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
	if verify == "buffer" then
		vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
	else
		vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
	end
	vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
	vim.api.nvim_buf_set_name(buf, "Praxis")
	vim.api.nvim_set_current_buf(buf)

	if is_challenge then
		local target = vim.fn.systemlist({ "praxis", "target", opts.args })[1]
		local result = vim.fn.systemlist({ "praxis", "result", opts.args }) or {}

		local function byte_to_char(line, bytecol)
			return vim.fn.strchars(string.sub(line, 1, bytecol))
		end

		local function check_buffer()
			local current = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
			if #current ~= #result then
				return
			end
			for i = 1, #current do
				if current[i] ~= result[i] then
					return
				end
			end
			state.done = true
			vim.api.nvim_echo({ { "Success" } }, false, {})
			render_result()
		end

		local state = {
			done            = false,
			moves           = 0,
			start_ns        = vim.uv.hrtime(),
			challenge_lines = lines,
			target          = target,
			verify          = verify,
			result_lines    = result,
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
			if state.verify ~= "buffer" then
				vim.api.nvim_set_option_value("modifiable", false, { buf = buf })
			end
			vim.api.nvim_win_set_cursor(0, { 1, 0 })
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
	end
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })

return M
