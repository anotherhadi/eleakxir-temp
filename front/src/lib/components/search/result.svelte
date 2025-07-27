<script lang="ts">
  import * as Table from "$lib/components/ui/table/index.js";
  import {
    ChevronsUpDown,
    ClipboardCopy,
    Database as DatabaseIcon,
  } from "@lucide/svelte";
  import * as Collapsible from "$lib/components/ui/collapsible/index.js";
  import Button from "$src/lib/components/ui/button/button.svelte";
  import {
    GetIndexOrFirstNonEmpty,
    GetIndexToPrioritise,
  } from "$src/lib/leaks";
  import { clsx } from "clsx";
  import { toast } from "svelte-sonner";
  import type { Result } from "$src/lib/types";

  const { result }: { result: Result } = $props();

  var open = $state<boolean>(false);

  function copy() {
    navigator.clipboard.writeText(JSON.stringify(result));
    toast.success("Document copied to clipboard!");
  }
</script>

<Collapsible.Root bind:open>
  <div
    class={clsx(
      "h-18 border-b w-full flex items-center hover:bg-card/40 px-8 gap-2 overflow-scroll whitespace-nowrap",
      open && "bg-card/40 border-b",
    )}
  >
    <Collapsible.Trigger
      class="flex h-full items-center justify-between w-full gap-2"
    >
      <div>
        {GetIndexOrFirstNonEmpty(result.Content, GetIndexToPrioritise(result))}
      </div>
      <div class="flex gap-2 justify-center items-center">
        <div class="flex gap-2 items-center">
          <DatabaseIcon size={14} />
          <p class="text-muted-foreground">
            {result.DataleakName}
          </p>
        </div>
        <Button variant="ghost">
          <ChevronsUpDown />
        </Button>
      </div>
    </Collapsible.Trigger>
    {#if open}
      <Button variant="ghost" onclick={copy}>
        <ClipboardCopy />
      </Button>
    {/if}
  </div>

  <Collapsible.Content class="">
    {#each result.Content as cell, index}
      {#if cell !== "" && result.Columns[index] !== ""}
        <Table.Row class="odd:bg-card/50 even:bg-card/30 grid grid-cols-6 py-2">
          <Table.Cell class="pl-12 col-span-2 text-muted-foreground overflow-x-scroll"
            >{result.Columns[index]}</Table.Cell
          >
          <Table.Cell class="col-span-4 w-full overflow-x-scroll">{cell}</Table.Cell>
        </Table.Row>
      {/if}
    {/each}
  </Collapsible.Content>
</Collapsible.Root>
