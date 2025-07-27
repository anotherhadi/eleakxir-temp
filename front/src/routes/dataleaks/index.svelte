<script lang="ts">
  import { type Dataleak, type Dataleaks } from "$src/lib/types";
  import axios from "axios";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import DatawellsTable from "$src/lib/components/misc/datawells-table.svelte";
  import Input from "$src/lib/components/ui/input/input.svelte";
  import { Search } from "@lucide/svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";

  var dataleaks = $state<Dataleaks | null>(null);

  let page = $state(1);
  let search = $state<string>("");
  let dataleaksFilter = $state<Dataleak[] | null>(null);

  onMount(async () => {
    axios
      .get("/api/dataleaks")
      .then((r) => {
        dataleaks = r.data.dataleaks;
        if(dataleaks){
          dataleaksFilter = dataleaks.Dataleaks;
        }
      })
      .catch((error) => {
        console.error("Error fetching dataleaks:", error);
        toast.error("Failed to fetch dataleaks. Please try again later.");
      });
  });

  function updateFilter() {
    if (dataleaks && dataleaksFilter) {
      page = 1;
      dataleaksFilter = dataleaks.Dataleaks.filter((d) => {
        return (
          d.Name.toLowerCase().includes(search.toLowerCase()) ||
          d.Columns.some((c) => c.toLowerCase().includes(search.toLowerCase()))
        );
      });
    }
  }
</script>

<main>
  <h1>Data wells</h1>

  <p class="text-muted-foreground">
    This page lets you explore all the indexed dataleaks stored on the server.
    You can filter by filename or by available columns. For each leak, you'll
    see its size, number of entries, available fields, and get access to more
    detailed information.
    To add new data wells, you just have to move them to the datawells folder of the backend.
  </p>

  <hr />

  <form
    class="flex flex-col sm:flex-row gap-2 sm:gap-0 items-stretch border rounded-md h-10 mb-6"
    onsubmit={(e) => {
      e.preventDefault();
      updateFilter();
    }}
  >
    <div class="relative w-full">
      <Input
        bind:value={search}
        type="text"
        placeholder="Filter"
        class="text-xl border-none h-full w-full rounded-none sm:rounded-l-md"
      />
    </div>
    <Button
      variant="secondary"
      size="lg"
      type="submit"
      class="sm:rounded-l-none h-full"
    >
      <Search class="size-5" />
    </Button>
  </form>

  {#if dataleaks}
    {#if dataleaksFilter === null}
      <p class="text-muted-foreground">Loading dataleaks...</p>
    {:else if dataleaksFilter.length === 0}
      <p class="text-muted-foreground">No dataleaks found.</p>
    {:else}
      <DatawellsTable
        dataleaks={dataleaksFilter}
        totalDocuments={dataleaks.TotalRows}
        totalSize={dataleaks.TotalSize}
        totalDataleaks={dataleaks.TotalDataleaks}
        bind:page
      />
    {/if}
  {/if}
</main>
