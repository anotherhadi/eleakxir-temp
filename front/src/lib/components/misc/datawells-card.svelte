<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import { DatabaseIcon } from "@lucide/svelte";
  import LastDatawells from "../last-datawells.svelte";
  import Button from "../ui/button/button.svelte";

  const { dataleaks } = $props();
</script>

<Card.Root>
  <Card.Header class="gap-5">
    <div
      class="rounded-full w-16 h-16 bg-orange-400/40 grid items-center justify-center"
    >
      <DatabaseIcon class="size-10" />
    </div>
    <Card.Title>Last Data Wells Added</Card.Title>
  </Card.Header>
  <Card.Content>
    {#if dataleaks && dataleaks.Dataleaks.length}
      <LastDatawells dataleaks={dataleaks.Dataleaks.slice(-5).reverse()} />
      <Button variant="link" size="sm" href="/dataleaks"
        >See all data wells</Button
      >
    {:else}
      <p class="text-muted-foreground">
        Cannot find any data wells.
        Refresh the page or check your connection.
      </p>
    {/if}
  </Card.Content>
</Card.Root>
