import html2text

html_content = '''<div class="rounded-lg border bg-card text-card-foreground shadow-xs"><div class="flex flex-col space-y-1.5 p-6"><div class="text-2xl font-semibold leading-none tracking-tight">Content</div></div><div class="p-6 pt-0"><div data-color-mode="light"><div class="wmde-markdown wmde-markdown-color markdown text-xs"><h1 id="context7-mcp---up-to-date-code-docs-for-any-prompt"><a class="anchor" aria-hidden="true" tabindex="-1" href="#context7-mcp---up-to-date-code-docs-for-any-prompt" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>Context7 MCP - Up-to-date Code Docs For Any Prompt</h1>
<p><a href="https://context7.com" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Website" src="https://img.shields.io/badge/Website-context7.com-blue"></a> <a href="https://smithery.ai/server/@upstash/context7-mcp" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="smithery badge" src="https://smithery.ai/badge/@upstash/context7-mcp"></a> <a href="https://insiders.vscode.dev/redirect?url=vscode%3Amcp%2Finstall%3F%7B%22name%22%3A%22context7%22%2C%22command%22%3A%22npx%22%2C%22args%22%3A%5B%22-y%22%2C%22%40upstash%2Fcontext7-mcp%40latest%22%5D%7D" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Install in VS Code (npx)" src="https://img.shields.io/badge/VS_Code-VS_Code?style=flat-square&amp;label=Install%20Context7%20MCP&amp;color=0098FF"></a></p>
<p><a href="./docs/README.zh-CN.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="中文文档" src="https://img.shields.io/badge/docs-%E4%B8%AD%E6%96%87%E7%89%88-yellow"></a> <a href="./docs/README.ko.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="한국어 문서" src="https://img.shields.io/badge/docs-%ED%95%9C%EA%B5%AD%EC%96%B4-green"></a> <a href="./docs/README.es.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Documentación en Español" src="https://img.shields.io/badge/docs-Espa%C3%B1ol-orange"></a> <a href="./docs/README.fr.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Documentation en Français" src="https://img.shields.io/badge/docs-Fran%C3%A7ais-blue"></a> <a href="./docs/README.pt-BR.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Documentação em Português (Brasil)" src="https://img.shields.io/badge/docs-Portugu%C3%AAs%20(Brasil)-purple"></a> <a href="./docs/README.it.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Documentazione in italiano" src="https://img.shields.io/badge/docs-Italian-red"></a> <a href="./docs/README.id-ID.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Dokumentasi Bahasa Indonesia" src="https://img.shields.io/badge/docs-Bahasa%20Indonesia-pink"></a> <a href="./docs/README.de.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Dokumentation auf Deutsch" src="https://img.shields.io/badge/docs-Deutsch-darkgreen"></a> <a href="./docs/README.ru.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Документация на русском языке" src="https://img.shields.io/badge/docs-%D0%A0%D1%83%D1%81%D1%81%D0%BA%D0%B8%D0%B9-darkblue"></a> <a href="./docs/README.tr.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Türkçe Doküman" src="https://img.shields.io/badge/docs-T%C3%BCrk%C3%A7e-blue"></a> <a href="./docs/README.ar.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Arabic Documentation" src="https://img.shields.io/badge/docs-Arabic-white"></a></p>
<h2 id="-without-context7"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-without-context7" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>❌ Without Context7</h2>
<p>LLMs rely on outdated or generic information about the libraries you use. You get:</p>
<ul>
<li>❌ Code examples are outdated and based on year-old training data</li>
<li>❌ Hallucinated APIs don't even exist</li>
<li>❌ Generic answers for old package versions</li>
</ul>
<h2 id="-with-context7"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-with-context7" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>✅ With Context7</h2>
<p>Context7 MCP pulls up-to-date, version-specific documentation and code examples straight from the source — and places them directly into your prompt.</p>
<p>Add <code>use context7</code> to your prompt in Cursor:</p>
<pre class="language-txt"><code class="language-txt code-highlight hljs language-plaintext" data-highlighted="yes">Create a basic Next.js project with app router. use context7
</code><div class="copied" data-code="Create a basic Next.js project with app router. use context7
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<pre class="language-txt"><code class="language-txt code-highlight hljs language-plaintext" data-highlighted="yes">Create a script to delete the rows where the city is "" given PostgreSQL credentials. use context7
</code><div class="copied" data-code="Create a script to delete the rows where the city is &quot;&quot; given PostgreSQL credentials. use context7
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>Context7 fetches up-to-date code examples and documentation right into your LLM's context.</p>
<ul>
<li>1️⃣ Write your prompt naturally</li>
<li>2️⃣ Tell the LLM to <code>use context7</code></li>
<li>3️⃣ Get working code answers</li>
</ul>
<p>No tab-switching, no hallucinated APIs that don't exist, no outdated code generations.</p>
<h2 id="-adding-projects"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-adding-projects" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>📚 Adding Projects</h2>
<p>Check out our <a href="./docs/adding-projects.md" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">project addition guide</a> to learn how to add (or update) your favorite libraries to Context7.</p>
<h2 id="️-installation"><a class="anchor" aria-hidden="true" tabindex="-1" href="#️-installation" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>🛠️ Installation</h2>
<h3 id="requirements"><a class="anchor" aria-hidden="true" tabindex="-1" href="#requirements" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>Requirements</h3>
<ul>
<li>Node.js &gt;= v18.0.0</li>
<li>Cursor, Windsurf, Claude Desktop or another MCP Client</li>
</ul>
<details>
<summary><b>Installing via Smithery</b></summary>
<p>To install Context7 MCP Server for any client automatically via <a href="https://smithery.ai/server/@upstash/context7-mcp" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Smithery</a>:</p>
<pre class="language-bash"><code class="language-bash code-highlight hljs" data-highlighted="yes">npx -y @smithery/cli@latest install @upstash/context7-mcp --client &lt;CLIENT_NAME&gt; --key &lt;YOUR_SMITHERY_KEY&gt;
</code><div class="copied" data-code="npx -y @smithery/cli@latest install @upstash/context7-mcp --client <CLIENT_NAME> --key <YOUR_SMITHERY_KEY>
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>You can find your Smithery key in the <a href="https://smithery.ai/server/@upstash/context7-mcp" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Smithery.ai webpage</a>.</p>
</details>
<details>
<summary><b>Install in Cursor</b></summary>
<p>Go to: <code>Settings</code> -&gt; <code>Cursor Settings</code> -&gt; <code>MCP</code> -&gt; <code>Add new global MCP server</code></p>
<p>Pasting the following configuration into your Cursor <code>~/.cursor/mcp.json</code> file is the recommended approach. You may also install in a specific project by creating <code>.cursor/mcp.json</code> in your project folder. See <a href="https://docs.cursor.com/context/model-context-protocol" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Cursor MCP docs</a> for more info.</p>
<h4 id="cursor-remote-server-connection"><a aria-hidden="true" tabindex="-1" href="#cursor-remote-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Cursor Remote Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"url"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"https://mcp.context7.com/mcp"</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;url&quot;: &quot;https://mcp.context7.com/mcp&quot;
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<h4 id="cursor-local-server-connection"><a aria-hidden="true" tabindex="-1" href="#cursor-local-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Cursor Local Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<details>
<summary>Alternative: Use Bun</summary>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"bunx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;bunx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary>Alternative: Use Deno</summary>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"deno"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"run"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"--allow-env"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"--allow-net"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"npm:@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;deno&quot;,
      &quot;args&quot;: [&quot;run&quot;, &quot;--allow-env&quot;, &quot;--allow-net&quot;, &quot;npm:@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
</details>
<details>
<summary><b>Install in Windsurf</b></summary>
<p>Add this to your Windsurf MCP config file. See <a href="https://docs.windsurf.com/windsurf/mcp" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Windsurf MCP docs</a> for more info.</p>
<h4 id="windsurf-remote-server-connection"><a aria-hidden="true" tabindex="-1" href="#windsurf-remote-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Windsurf Remote Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"serverUrl"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"https://mcp.context7.com/sse"</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;serverUrl&quot;: &quot;https://mcp.context7.com/sse&quot;
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<h4 id="windsurf-local-server-connection"><a aria-hidden="true" tabindex="-1" href="#windsurf-local-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Windsurf Local Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in VS Code</b></summary>
<p><a href="https://insiders.vscode.dev/redirect?url=vscode%3Amcp%2Finstall%3F%7B%22name%22%3A%22context7%22%2C%22command%22%3A%22npx%22%2C%22args%22%3A%5B%22-y%22%2C%22%40upstash%2Fcontext7-mcp%40latest%22%5D%7D" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Install in VS Code (npx)" src="https://img.shields.io/badge/VS_Code-VS_Code?style=flat-square&amp;label=Install%20Context7%20MCP&amp;color=0098FF"></a>
<a href="https://insiders.vscode.dev/redirect?url=vscode-insiders%3Amcp%2Finstall%3F%7B%22name%22%3A%22context7%22%2C%22command%22%3A%22npx%22%2C%22args%22%3A%5B%22-y%22%2C%22%40upstash%2Fcontext7-mcp%40latest%22%5D%7D" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Install in VS Code Insiders (npx)" src="https://img.shields.io/badge/VS_Code_Insiders-VS_Code_Insiders?style=flat-square&amp;label=Install%20Context7%20MCP&amp;color=24bfa5"></a></p>
<p>Add this to your VS Code MCP config file. See <a href="https://code.visualstudio.com/docs/copilot/chat/mcp-servers" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">VS Code MCP docs</a> for more info.</p>
<h4 id="vs-code-remote-server-connection"><a aria-hidden="true" tabindex="-1" href="#vs-code-remote-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>VS Code Remote Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-attr">"mcp"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"servers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"type"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"http"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"url"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"https://mcp.context7.com/mcp"</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="&quot;mcp&quot;: {
  &quot;servers&quot;: {
    &quot;context7&quot;: {
      &quot;type&quot;: &quot;http&quot;,
      &quot;url&quot;: &quot;https://mcp.context7.com/mcp&quot;
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<h4 id="vs-code-local-server-connection"><a aria-hidden="true" tabindex="-1" href="#vs-code-local-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>VS Code Local Server Connection</h4>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-attr">"mcp"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"servers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"type"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"stdio"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="&quot;mcp&quot;: {
  &quot;servers&quot;: {
    &quot;context7&quot;: {
      &quot;type&quot;: &quot;stdio&quot;,
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in Zed</b></summary>
<p>It can be installed via <a href="https://zed.dev/extensions?query=Context7" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Zed Extensions</a> or you can add this to your Zed <code>settings.json</code>. See <a href="https://zed.dev/docs/assistant/context-servers" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Zed Context Server docs</a> for more info.</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"context_servers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"Context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
        <span class="hljs-attr">"path"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
        <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
      <span class="hljs-punctuation">}</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"settings"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span><span class="hljs-punctuation">}</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;context_servers&quot;: {
    &quot;Context7&quot;: {
      &quot;command&quot;: {
        &quot;path&quot;: &quot;npx&quot;,
        &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
      },
      &quot;settings&quot;: {}
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in Claude Code</b></summary>
<p>Run this command. See <a href="https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/tutorials#set-up-model-context-protocol-mcp" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Claude Code MCP docs</a> for more info.</p>
<h4 id="claude-code-remote-server-connection"><a aria-hidden="true" tabindex="-1" href="#claude-code-remote-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Claude Code Remote Server Connection</h4>
<pre class="language-sh"><code class="language-sh code-highlight hljs language-bash" data-highlighted="yes">claude mcp add --transport sse context7 https://mcp.context7.com/sse
</code><div class="copied" data-code="claude mcp add --transport sse context7 https://mcp.context7.com/sse
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<h4 id="claude-code-local-server-connection"><a aria-hidden="true" tabindex="-1" href="#claude-code-local-server-connection" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><span class="icon icon-link"></span></a>Claude Code Local Server Connection</h4>
<pre class="language-sh"><code class="language-sh code-highlight hljs language-bash" data-highlighted="yes">claude mcp add context7 -- npx -y @upstash/context7-mcp
</code><div class="copied" data-code="claude mcp add context7 -- npx -y @upstash/context7-mcp
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in Claude Desktop</b></summary>
<p>Add this to your Claude Desktop <code>claude_desktop_config.json</code> file. See <a href="https://modelcontextprotocol.io/quickstart/user" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Claude Desktop MCP docs</a> for more info.</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"Context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;Context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in BoltAI</b></summary>
<p>Open the "Settings" page of the app, navigate to "Plugins," and enter the following JSON:</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>Once saved, enter in the chat <code>get-library-docs</code> followed by your Context7 documentation ID (e.g., <code>get-library-docs /nuxt/ui</code>). More information is available on <a href="https://docs.boltai.com/docs/plugins/mcp-servers" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">BoltAI's Documentation site</a>. For BoltAI on iOS, <a href="https://docs.boltai.com/docs/boltai-mobile/mcp-servers" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">see this guide</a>.</p>
</details>
<details>
<summary><b>Using Docker</b></summary>
<p>If you prefer to run the MCP server in a Docker container:</p>
<ol>
<li>
<p><strong>Build the Docker Image:</strong></p>
<p>First, create a <code>Dockerfile</code> in the project root (or anywhere you prefer):</p>
<details>
<summary>Click to see Dockerfile content</summary>
<pre class="language-dockerfile"><code class="language-Dockerfile code-highlight hljs" data-highlighted="yes"><span class="hljs-keyword">FROM</span> node:<span class="hljs-number">18</span>-alpine

<span class="hljs-keyword">WORKDIR</span><span class="language-bash"> /app</span>

<span class="hljs-comment"># Install the latest version globally</span>
<span class="hljs-keyword">RUN</span><span class="language-bash"> npm install -g @upstash/context7-mcp</span>

<span class="hljs-comment"># Expose default port if needed (optional, depends on MCP client interaction)</span>
<span class="hljs-comment"># EXPOSE 3000</span>

<span class="hljs-comment"># Default command to run the server</span>
<span class="hljs-keyword">CMD</span><span class="language-bash"> [<span class="hljs-string">"context7-mcp"</span>]</span>
</code><div class="copied" data-code="FROM node:18-alpine

WORKDIR /app

# Install the latest version globally
RUN npm install -g @upstash/context7-mcp

# Expose default port if needed (optional, depends on MCP client interaction)
# EXPOSE 3000

# Default command to run the server
CMD [&quot;context7-mcp&quot;]
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<p>Then, build the image using a tag (e.g., <code>context7-mcp</code>). <strong>Make sure Docker Desktop (or the Docker daemon) is running.</strong> Run the following command in the same directory where you saved the <code>Dockerfile</code>:</p>
<pre class="language-bash"><code class="language-bash code-highlight hljs" data-highlighted="yes">docker build -t context7-mcp .
</code><div class="copied" data-code="docker build -t context7-mcp .
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</li>
<li>
<p><strong>Configure Your MCP Client:</strong></p>
<p>Update your MCP client's configuration to use the Docker command.</p>
<p><em>Example for a cline_mcp_settings.json:</em></p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"Сontext7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"autoApprove"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-punctuation">]</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"disabled"</span><span class="hljs-punctuation">:</span> <span class="hljs-literal"><span class="hljs-keyword">false</span></span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"timeout"</span><span class="hljs-punctuation">:</span> <span class="hljs-number">60</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"docker"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"run"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"-i"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"--rm"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"context7-mcp"</span><span class="hljs-punctuation">]</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"transportType"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"stdio"</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;Сontext7&quot;: {
      &quot;autoApprove&quot;: [],
      &quot;disabled&quot;: false,
      &quot;timeout&quot;: 60,
      &quot;command&quot;: &quot;docker&quot;,
      &quot;args&quot;: [&quot;run&quot;, &quot;-i&quot;, &quot;--rm&quot;, &quot;context7-mcp&quot;],
      &quot;transportType&quot;: &quot;stdio&quot;
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p><em>Note: This is an example configuration. Please refer to the specific examples for your MCP client (like Cursor, VS Code, etc.) earlier in this README to adapt the structure (e.g., <code>mcpServers</code> vs <code>servers</code>). Also, ensure the image name in <code>args</code> matches the tag used during the <code>docker build</code> command.</em></p>
</li>
</ol>
</details>
<details>
<summary><b>Install in Windows</b></summary>
<p>The configuration on Windows is slightly different compared to Linux or macOS (<em><code>Cline</code> is used in the example</em>). The same principle applies to other editors; refer to the configuration of <code>command</code> and <code>args</code>.</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"github.com/upstash/context7-mcp"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"cmd"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"/c"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp@latest"</span><span class="hljs-punctuation">]</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"disabled"</span><span class="hljs-punctuation">:</span> <span class="hljs-literal"><span class="hljs-keyword">false</span></span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"autoApprove"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;github.com/upstash/context7-mcp&quot;: {
      &quot;command&quot;: &quot;cmd&quot;,
      &quot;args&quot;: [&quot;/c&quot;, &quot;npx&quot;, &quot;-y&quot;, &quot;@upstash/context7-mcp@latest&quot;],
      &quot;disabled&quot;: false,
      &quot;autoApprove&quot;: []
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Install in Augment Code</b></summary>
<p>To configure Context7 MCP in Augment Code, follow these steps:</p>
<ol>
<li>Press Cmd/Ctrl Shift P or go to the hamburger menu in the Augment panel</li>
<li>Select Edit Settings</li>
<li>Under Advanced, click Edit in settings.json</li>
<li>Add the server configuration to the <code>mcpServers</code> array in the <code>augment.advanced</code> object</li>
</ol>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-attr">"augment.advanced"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span>
        <span class="hljs-punctuation">{</span>
            <span class="hljs-attr">"name"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"context7"</span><span class="hljs-punctuation">,</span>
            <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
            <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
        <span class="hljs-punctuation">}</span>
    <span class="hljs-punctuation">]</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="&quot;augment.advanced&quot;: {
    &quot;mcpServers&quot;: [
        {
            &quot;name&quot;: &quot;context7&quot;,
            &quot;command&quot;: &quot;npx&quot;,
            &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
        }
    ]
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>Once the MCP server is added, restart your editor. If you receive any errors, check the syntax to make sure closing brackets or commas are not missing.</p>
</details>
<h2 id="-environment-variables"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-environment-variables" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>🔧 Environment Variables</h2>
<p>The Context7 MCP server supports the following environment variables:</p>
<ul>
<li><code>DEFAULT_MINIMUM_TOKENS</code>: Set the minimum token count for documentation retrieval (default: 10000)</li>
</ul>
<p>Example configuration with environment variables:</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"env"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
        <span class="hljs-attr">"DEFAULT_MINIMUM_TOKENS"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"6000"</span>
      <span class="hljs-punctuation">}</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;],
      &quot;env&quot;: {
        &quot;DEFAULT_MINIMUM_TOKENS&quot;: &quot;6000&quot;
      }
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<h2 id="-available-tools"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-available-tools" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>🔨 Available Tools</h2>
<p>Context7 MCP provides the following tools that LLMs can use:</p>
<ul>
<li>
<p><code>resolve-library-id</code>: Resolves a general library name into a Context7-compatible library ID.</p>
<ul>
<li><code>libraryName</code> (required): The name of the library to search for</li>
</ul>
</li>
<li>
<p><code>get-library-docs</code>: Fetches documentation for a library using a Context7-compatible library ID.</p>
<ul>
<li><code>context7CompatibleLibraryID</code> (required): Exact Context7-compatible library ID (e.g., <code>/mongodb/docs</code>, <code>/vercel/next.js</code>)</li>
<li><code>topic</code> (optional): Focus the docs on a specific topic (e.g., "routing", "hooks")</li>
<li><code>tokens</code> (optional, default 10000): Max number of tokens to return. Values less than the configured <code>DEFAULT_MINIMUM_TOKENS</code> value or the default value of 10000 are automatically increased to that value.</li>
</ul>
</li>
</ul>
<h2 id="-development"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-development" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>💻 Development</h2>
<p>Clone the project and install dependencies:</p>
<pre class="language-bash"><code class="language-bash code-highlight hljs" data-highlighted="yes">bun i
</code><div class="copied" data-code="bun i
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>Build:</p>
<pre class="language-bash"><code class="language-bash code-highlight hljs" data-highlighted="yes">bun run build
</code><div class="copied" data-code="bun run build
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<details>
<summary><b>Local Configuration Example</b></summary>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"tsx"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"/path/to/folder/context7-mcp/src/index.ts"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;tsx&quot;, &quot;/path/to/folder/context7-mcp/src/index.ts&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>Testing with MCP Inspector</b></summary>
<pre class="language-bash"><code class="language-bash code-highlight hljs" data-highlighted="yes">npx -y @modelcontextprotocol/inspector npx @upstash/context7-mcp
</code><div class="copied" data-code="npx -y @modelcontextprotocol/inspector npx @upstash/context7-mcp
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<h2 id="-troubleshooting"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-troubleshooting" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>🚨 Troubleshooting</h2>
<details>
<summary><b>Module Not Found Errors</b></summary>
<p>If you encounter <code>ERR_MODULE_NOT_FOUND</code>, try using <code>bunx</code> instead of <code>npx</code>:</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"bunx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;bunx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
<p>This often resolves module resolution issues in environments where <code>npx</code> doesn't properly install or resolve packages.</p>
</details>
<details>
<summary><b>ESM Resolution Issues</b></summary>
<p>For errors like <code>Error: Cannot find module 'uriTemplate.js'</code>, try the <code>--experimental-vm-modules</code> flag:</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"--node-options=--experimental-vm-modules"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp@1.0.6"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;--node-options=--experimental-vm-modules&quot;, &quot;@upstash/context7-mcp@1.0.6&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>TLS/Certificate Issues</b></summary>
<p>Use the <code>--experimental-fetch</code> flag to bypass TLS-related problems:</p>
<pre class="language-json"><code class="language-json code-highlight hljs" data-highlighted="yes"><span class="hljs-punctuation">{</span>
  <span class="hljs-attr">"mcpServers"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
    <span class="hljs-attr">"context7"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">{</span>
      <span class="hljs-attr">"command"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"npx"</span><span class="hljs-punctuation">,</span>
      <span class="hljs-attr">"args"</span><span class="hljs-punctuation">:</span> <span class="hljs-punctuation">[</span><span class="hljs-string">"-y"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"--node-options=--experimental-fetch"</span><span class="hljs-punctuation">,</span> <span class="hljs-string">"@upstash/context7-mcp"</span><span class="hljs-punctuation">]</span>
    <span class="hljs-punctuation">}</span>
  <span class="hljs-punctuation">}</span>
<span class="hljs-punctuation">}</span>
</code><div class="copied" data-code="{
  &quot;mcpServers&quot;: {
    &quot;context7&quot;: {
      &quot;command&quot;: &quot;npx&quot;,
      &quot;args&quot;: [&quot;-y&quot;, &quot;--node-options=--experimental-fetch&quot;, &quot;@upstash/context7-mcp&quot;]
    }
  }
}
"><svg class="octicon-copy" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg><svg class="octicon-check" aria-hidden="true" viewBox="0 0 16 16" fill="currentColor" height="12" width="12"><path fill-rule="evenodd" d="M13.78 4.22a.75.75 0 010 1.06l-7.25 7.25a.75.75 0 01-1.06 0L2.22 9.28a.75.75 0 011.06-1.06L6 10.94l6.72-6.72a.75.75 0 011.06 0z"></path></svg></div></pre>
</details>
<details>
<summary><b>General MCP Client Errors</b></summary>
<ol>
<li>Try adding <code>@latest</code> to the package name</li>
<li>Use <code>bunx</code> as an alternative to <code>npx</code></li>
<li>Consider using <code>deno</code> as another alternative</li>
<li>Ensure you're using Node.js v18 or higher for native fetch support</li>
</ol>
</details>
<h2 id="️-disclaimer"><a class="anchor" aria-hidden="true" tabindex="-1" href="#️-disclaimer" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>⚠️ Disclaimer</h2>
<p>Context7 projects are community-contributed and while we strive to maintain high quality, we cannot guarantee the accuracy, completeness, or security of all library documentation. Projects listed in Context7 are developed and maintained by their respective owners, not by Context7. If you encounter any suspicious, inappropriate, or potentially harmful content, please use the "Report" button on the project page to notify us immediately. We take all reports seriously and will review flagged content promptly to maintain the integrity and safety of our platform. By using Context7, you acknowledge that you do so at your own discretion and risk.</p>
<h2 id="-connect-with-us"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-connect-with-us" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>🤝 Connect with Us</h2>
<p>Stay updated and join our community:</p>
<ul>
<li>📢 Follow us on <a href="https://x.com/contextai" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">X</a> for the latest news and updates</li>
<li>🌐 Visit our <a href="https://context7.com" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Website</a></li>
<li>💬 Join our <a href="https://upstash.com/discord" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Discord Community</a></li>
</ul>
<h2 id="-context7-in-media"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-context7-in-media" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>📺 Context7 In Media</h2>
<ul>
<li><a href="https://youtu.be/52FC3qObp9E" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Better Stack: "Free Tool Makes Cursor 10x Smarter"</a></li>
<li><a href="https://www.youtube.com/watch?v=G7gK8H6u7Rs" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Cole Medin: "This is Hands Down the BEST MCP Server for AI Coding Assistants"</a></li>
<li><a href="https://www.youtube.com/watch?v=-ggvzyLpK6o" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Income Stream Surfers: "Context7 + SequentialThinking MCPs: Is This AGI?"</a></li>
<li><a href="https://www.youtube.com/watch?v=CTZm6fBYisc" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Julian Goldie SEO: "Context7: New MCP AI Agent Update"</a></li>
<li><a href="https://www.youtube.com/watch?v=-ls0D-rtET4" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">JeredBlu: "Context 7 MCP: Get Documentation Instantly + VS Code Setup"</a></li>
<li><a href="https://www.youtube.com/watch?v=PS-2Azb-C3M" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Income Stream Surfers: "Context7: The New MCP Server That Will CHANGE AI Coding"</a></li>
<li><a href="https://www.youtube.com/watch?v=qZfENAPMnyo" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">AICodeKing: "Context7 + Cline &amp; RooCode: This MCP Server Makes CLINE 100X MORE EFFECTIVE!"</a></li>
<li><a href="https://www.youtube.com/watch?v=LqTQi8qexJM" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow">Sean Kochel: "5 MCP Servers For Vibe Coding Glory (Just Plug-In &amp; Go)"</a></li>
</ul>
<h2 id="-star-history"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-star-history" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>⭐ Star History</h2>
<p><a href="https://www.star-history.com/#upstash/context7&amp;Date" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><img alt="Star History Chart" src="https://api.star-history.com/svg?repos=upstash/context7&amp;type=Date"></a></p>
<h2 id="-license"><a class="anchor" aria-hidden="true" tabindex="-1" href="#-license" node="[object Object]" target="_blank" rel="noopener noreferrer nofollow"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>📄 License</h2>
<p>MIT</p></div></div></div></div>'''
text_content = html2text.html2text(html_content, bodywidth=0)

print("HTML Content:")
# print(html_content)

print("\nMarkdown Content:")
print(text_content)