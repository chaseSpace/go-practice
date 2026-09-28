<script setup lang="ts">
const props = defineProps<{ content: string }>()

type InlineSegment = { kind: 'text' | 'bold' | 'code'; text: string }

function segments(value: string): InlineSegment[] {
  const result: InlineSegment[] = []
  const matcher = /\*\*(.+?)\*\*|`([^`]+)`/g
  let cursor = 0
  for (const match of value.matchAll(matcher)) {
    if (match.index! > cursor) result.push({ kind: 'text', text: value.slice(cursor, match.index) })
    result.push(match[1] ? { kind: 'bold', text: match[1] } : { kind: 'code', text: match[2] })
    cursor = match.index! + match[0].length
  }
  if (cursor < value.length) result.push({ kind: 'text', text: value.slice(cursor) })
  return result
}
</script>

<template>
  <template v-for="(segment, index) in segments(props.content)" :key="index">
    <strong v-if="segment.kind === 'bold'">{{ segment.text }}</strong>
    <code v-else-if="segment.kind === 'code'">{{ segment.text }}</code>
    <template v-else>{{ segment.text }}</template>
  </template>
</template>
