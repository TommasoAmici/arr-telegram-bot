venv:
	uv sync

lint: venv
	ruff check

format: venv
	ruff format
