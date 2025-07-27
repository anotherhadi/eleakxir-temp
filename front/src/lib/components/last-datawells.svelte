<script lang="ts">
  import { DatabaseIcon } from "@lucide/svelte";

  const { dataleaks } = $props();
</script>

{#if dataleaks.length === 0}
  <p class="text-muted-foreground text-center">No data wells available.</p>
{:else}
  {#each dataleaks as d}
    <div class="flex items-center gap-2 my-2">
      <DatabaseIcon class="text-muted-foreground" />
      <span>{d.Name}</span>
      <span class="text-muted-foreground"
        >({d.Length.toLocaleString("fr-FR")} documents)</span
      >
    </div>
  {/each}
{/if}
