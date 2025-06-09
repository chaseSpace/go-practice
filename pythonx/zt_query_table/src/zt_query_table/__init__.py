def main() -> None:
    from .impl import mcp
    mcp.run(transport='sse')
