local M = {}

local session = nil

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

		if not session then
			session = {
				start_ns     = vim.uv.hrtime(),
				challenges   = 0,
				completions  = 0,
				total_moves  = 0,
				total_time_ms = 0,
			}
		end
		session.challenges = session.challenges + 1
	vim.fn.system({ "praxis", "attempt", opts.args })

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
			challenge_id    = opts.args,
		}

		local function render_result()
			local elapsed_ms = math.floor((vim.uv.hrtime() - state.start_ns) / 1e6)
			session.completions = session.completions + 1
			session.total_moves = session.total_moves + state.moves
			session.total_time_ms = session.total_time_ms + elapsed_ms
			vim.fn.system({ "praxis", "record", opts.args, tostring(state.moves), tostring(elapsed_ms) })
			local stats_out = vim.fn.systemlist({ "praxis", "stats", opts.args })
			vim.api.nvim_set_option_value("modifiable", true, { buf = buf })
			local display = {
				"Success", "",
				"Moves: " .. state.moves,
				"Time: "  .. elapsed_ms .. "ms", "",
			}
			for _, line in ipairs(stats_out) do
				table.insert(display, line)
			end
			table.insert(display, "")
			table.insert(display, "Press r to replay")
			vim.api.nvim_buf_set_lines(buf, 0, -1, false, display)
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
			vim.fn.system({ "praxis", "attempt", state.challenge_id })
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

local function format_time(total_s)
	local m = math.floor(total_s / 60)
	local s = total_s % 60
	if m > 0 then return string.format("%dm%ds", m, s) end
	return string.format("%ds", s)
end

function M.show_session()
	local completions = session and session.completions or 0
	local challenges  = session and session.challenges or 0
	local moves       = session and session.total_moves or 0
	local time_ms     = session and session.total_time_ms or 0
	local total_s     = math.floor(time_ms / 1000)
	local avg_moves   = completions > 0 and math.floor(moves / completions) or 0
	local avg_time_s  = completions > 0 and math.floor(time_ms / completions / 1000) or 0
	local elapsed_s   = session and math.floor((vim.uv.hrtime() - session.start_ns) / 1e9) or 0

	local buf = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_lines(buf, 0, -1, false, {
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
	vim.api.nvim_set_option_value("buftype", "nofile", { buf = buf })
	vim.api.nvim_buf_set_name(buf, "Praxis Session")
	vim.api.nvim_set_current_buf(buf)
end

vim.api.nvim_create_user_command("Praxis", M.show, { nargs = "?" })
vim.api.nvim_create_user_command("PraxisSession", M.show_session, {})

return M
