<script setup lang="ts">
import { computed } from 'vue'
import MarkdownInline from './MarkdownInline.vue'

const props = defineProps<{ content: string }>()

type MarkdownBlock =
  | { kind: 'heading'; level: number; text: string }
  | { kind: 'paragraph'; text: string }
  | { kind: 'list'; ordered: boolean; items: string[] }
  | { kind: 'quote'; text: string }
  | { kind: 'code'; language: string; content: string }
  | { kind: 'table'; headers: string[]; rows: string[][] }

const blocks = computed(() => parseMarkdown(props.content))

function parseMarkdown(markdown: string): MarkdownBlock[] {
  const lines = markdown.replace(/\r\n/g, '\n').split('\n')
  const result: MarkdownBlock[] = []
  let paragraph: string[] = []
  const flushParagraph = () => {
    if (paragraph.length) result.push({ kind: 'paragraph', text: paragraph.join('\n') })
    paragraph = []
  }
  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index]
    const codeFence = line.match(/^```\s*([^`]*)$/)
    const heading = line.match(/^(#{1,3})\s+(.+)$/)
    const unordered = line.match(/^[-*+]\s+(.+)$/)
    const ordered = line.match(/^\d+[.)]\s+(.+)$/)
    const quote = line.match(/^>\s?(.+)$/)
    const tableHeaders = splitTableRow(line)
    if (codeFence) {
      flushParagraph()
      const code: string[] = []
      index += 1
      while (index < lines.length && !lines[index].startsWith('```')) {
        code.push(lines[index])
        index += 1
      }
      result.push({ kind: 'code', language: codeFence[1].trim(), content: code.join('\n') })
    } else if (heading) {
      flushParagraph()
      result.push({ kind: 'heading', level: heading[1].length, text: heading[2] })
    } else if (unordered || ordered) {
      flushParagraph()
      const isOrdered = Boolean(ordered)
      const items: string[] = []
      while (index < lines.length) {
        const item = isOrdered ? lines[index].match(/^\d+[.)]\s+(.+)$/) : lines[index].match(/^[-*+]\s+(.+)$/)
        if (!item) break
        items.push(item[1])
        index += 1
      }
      index -= 1
      result.push({ kind: 'list', ordered: isOrdered, items })
    } else if (tableHeaders && isTableDelimiter(lines[index + 1] ?? '')) {
      flushParagraph()
      const rows: string[][] = []
      index += 2
      while (index < lines.length) {
        const row = splitTableRow(lines[index])
        if (!row || row.length !== tableHeaders.length) break
        rows.push(row)
        index += 1
      }
      index -= 1
      result.push({ kind: 'table', headers: tableHeaders, rows })
    } else if (quote) {
      flushParagraph()
      result.push({ kind: 'quote', text: quote[1] })
    } else if (!line.trim()) {
      flushParagraph()
    } else {
      paragraph.push(line)
    }
  }
  flushParagraph()
  return result
}

function splitTableRow(line: string): string[] | undefined {
  if (!line.includes('|')) return undefined
  const cells = line.trim().replace(/^\|/, '').replace(/\|$/, '').split('|').map((cell) => cell.trim())
  return cells.length > 1 ? cells : undefined
}

function isTableDelimiter(line: string): boolean {
  const cells = splitTableRow(line)
  return Boolean(cells?.length && cells.every((cell) => /^:?-{3,}:?$/.test(cell)))
}
</script>

<template>
  <article class="markdown-article">
    <template v-for="(block, index) in blocks" :key="index">
      <component :is="`h${block.level}`" v-if="block.kind === 'heading'"><MarkdownInline :content="block.text" /></component>
      <p v-else-if="block.kind === 'paragraph'"><MarkdownInline :content="block.text" /></p>
      <pre v-else-if="block.kind === 'code'" class="markdown-code"><code>{{ block.content }}</code></pre>
      <component :is="block.ordered ? 'ol' : 'ul'" v-else-if="block.kind === 'list'">
        <li v-for="(item, itemIndex) in block.items" :key="itemIndex"><MarkdownInline :content="item" /></li>
      </component>
      <div v-else-if="block.kind === 'table'" class="markdown-table-wrap"><table><thead><tr><th v-for="(header, headerIndex) in block.headers" :key="headerIndex"><MarkdownInline :content="header" /></th></tr></thead><tbody><tr v-for="(row, rowIndex) in block.rows" :key="rowIndex"><td v-for="(cell, cellIndex) in row" :key="cellIndex"><MarkdownInline :content="cell" /></td></tr></tbody></table></div>
      <blockquote v-else><MarkdownInline :content="block.text" /></blockquote>
    </template>
  </article>
</template>

<style scoped>
.markdown-code {
  margin: 18px 0;
  max-width: 100%;
  overflow: auto;
  padding: 14px;
  border-radius: 8px;
  background: #0f172a;
  color: #e2e8f0;
  font: 13px/1.6 ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  white-space: pre;
}

.markdown-code code {
  padding: 0;
  background: transparent;
  color: inherit;
  font: inherit;
}
</style>
