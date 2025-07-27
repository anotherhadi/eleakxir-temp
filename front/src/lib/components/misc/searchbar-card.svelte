<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import { GraduationCap, SquareEqual } from "@lucide/svelte";
</script>

<Card.Root>
  <Card.Header class="gap-5">
    <div
      class="rounded-full w-16 h-16 bg-green-400/40 grid items-center justify-center"
    >
      <GraduationCap class="size-10" />
    </div>
    <Card.Title>How to use the search bar</Card.Title>
  </Card.Header>
  <Card.Content>
    <ul class="pl-5 list-disc space-y-2 text-sm leading-relaxed">
      <li>
        <strong>Tabs:</strong> Choose the index to search in. <code>All</code>
        includes all indexes. <code>"Full Text"</code> searches the entire leaks directly
        (very slow).
      </li>
      <li>
        <strong>Terms:</strong> Write multiple keywords (e.g.
        <code>john doe</code>) to find results containing <em>all</em> of them (AND
        logic).
      </li>
      <li>
        <strong>Exact Match:</strong> Toggle the <SquareEqual
          class="inline size-4 mx-1"
        /> icon to search for an exact phrase (e.g. <code>"john doe"</code> as one
        field).
      </li>
    </ul>
  </Card.Content>
</Card.Root>
