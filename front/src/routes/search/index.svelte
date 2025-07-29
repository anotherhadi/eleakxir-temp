<script lang="ts">
  import Input from "$src/lib/components/ui/input/input.svelte";
  import { Loader, Search, SquareEqual } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import Button from "$src/lib/components/ui/button/button.svelte";
  import type {
    Result as ResultType,
    Dataleaks,
    ResearchStatus,
  } from "$src/lib/types";
  import Toggle from "$src/lib/components/ui/toggle/toggle.svelte";
  import * as Tooltip from "$lib/components/ui/tooltip/index.js";
  import Stats from "$src/lib/components/search/stats.svelte";
  import ProgressBar from "$src/lib/components/search/progress-bar.svelte";
  import NoResult from "$src/lib/components/search/no-result.svelte";
  import Searching from "$src/lib/components/search/searching.svelte";
  import SearchbarCard from "$src/lib/components/misc/searchbar-card.svelte";
  import DatawellsCard from "$src/lib/components/misc/datawells-card.svelte";
  import ResultPaginated from "$src/lib/components/search/resultPaginated.svelte";
  import ContributionCard from "$src/lib/components/misc/contribution-card.svelte";
  import { onDestroy, onMount } from "svelte";
  import axios from "axios";

  var dataleaks = $state<Dataleaks | null>(null);

  onMount(async () => {
    axios
      .get("/api/dataleaks")
      .then((r) => {
        dataleaks = r.data.dataleaks;
      })
      .catch((error) => {
        console.error("Error fetching dataleaks:", error);
        toast.error("Failed to fetch dataleaks. Please try again later.");
      });
  });

  // // Gérer la fermeture de la connexion lorsque le composant est détruit
  onDestroy(() => {
    if (eventSource) {
      eventSource.close();
    }
  });

  var query = $state<string>("");
  var exactMatch = $state<boolean>(false);
  var page = $state<number>(1);

  var choosenColumn = $state<string>("all");
  var availableColumns = $state<string[]>([
    "all",
    "username",
    "email",
    "name",
    "password",
    "full text",
  ]);

  let progressPercentage = $state(0);
  let results = $state<ResultType[]>([]);
  let searchStatus = $state<ResearchStatus>("idle");
  var researchTime = $state<number | null>(null);

  let eventSource: any = null;
  let startTime: number | null = null;

  function startSearch() {
    const q = JSON.stringify({
      Terms: query.trim().split(" "),
      ExactMatch: exactMatch,
    });

    console.log("Chosen column:", choosenColumn);
    let cols = choosenColumn;
    if (cols === "all") {
      cols = "username,email,name,password";
    }
    cols = cols.replace(" ", "_");
    cols = cols.split(",").map((col) => {
      return col.trim() === "name" ? "full_name" : col.trim();
    }).join(",");
    if (cols === "name"){
      cols = "full_name";
    }
    console.log("Columns to search:", cols);

    // TODO: Doesn't work
    if (!q) {
      alert("Please enter a search query.");
      return;
    } else if (q.length < 6) {
      alert("Query is too short. Please enter at least 6 characters.");
      return;
    } else if (!cols) {
      alert("Please specify columns to search.");
      return;
    } else if (searchStatus === "searching") {
      alert("Search already in progress. Please wait for it to complete.");
      return;
    }

    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }

    results = [];
    page = 1;

    const url = `/api/search?q=${encodeURIComponent(q)}&columns=${encodeURIComponent(cols)}`;

    try {
      eventSource = new EventSource(url);

      eventSource.addEventListener("start", () => {
        searchStatus = "searching";
        startTime = performance.now();
        researchTime = null;
        progressPercentage = 0;
      });

      eventSource.addEventListener("progress", (event: any) => {
        const data = JSON.parse(event.data);
        progressPercentage = data.percentage;
      });

      eventSource.addEventListener("new_results", (event: any) => {
        const data = JSON.parse(event.data);
        if (data.results && data.results.length !== 0) {
          results = [...results, ...data.results];
        }
      });

      eventSource.addEventListener("file_error", (event: any) => {
        const data = JSON.parse(event.data);
        console.error("File error:", data);
      });

      eventSource.addEventListener("error", (event: any) => {
        const data = event.data
          ? JSON.parse(event.data)
          : { message: "Unknown error" };
        searchStatus = "error";
        toast.error(`Error: ${data.message || "An error occurred"}`, {
          description: "Please try again later.",
        });
        if (startTime !== null) {
          researchTime = performance.now() - startTime; // Calcule le temps même en cas d'erreur
        }
        console.error("General EventSource error:", data, event);
        eventSource.close();
        eventSource = null;
      });

      eventSource.addEventListener("complete", () => {
        progressPercentage = 100;
        searchStatus = "complete";
        if (startTime !== null) {
          researchTime = performance.now() - startTime; // Calcule le temps à la fin
        }
        eventSource.close();
        eventSource = null;
      });
    } catch (e: any) {
      console.error("Failed to create EventSource:", e);
      if (startTime !== null) {
        researchTime = performance.now() - startTime; // Calcule le temps même si la création échoue
      }
      searchStatus = "error";
    }
  }
</script>

<main>
  <header class="border-b pb-8 mb-10">
    <h1 class="mb-2">Search</h1>
    <p class="text-muted-foreground mb-8">
      Search across
      {dataleaks?.TotalRows.toLocaleString("fr") || "-- --- --- ---"}
      documents in your data wells.
    </p>

    <div class="flex flex-wrap gap-2 mb-6">
      {#each availableColumns as column}
        {#if column === choosenColumn}
          <Button class="capitalize" onclick={() => (choosenColumn = column)}>
            {column}
          </Button>
        {:else}
          <Button
            class="capitalize text-muted-foreground"
            variant="ghost"
            disabled={searchStatus === "searching"}
            onclick={() => (choosenColumn = column)}
          >
            {column}
          </Button>
        {/if}
      {/each}
    </div>

    <form
      class="flex flex-col sm:flex-row gap-2 sm:gap-0 items-stretch border rounded-md h-16"
      onsubmit={(e) => {
        e.preventDefault();
        startSearch();
      }}
    >
      <div class="relative w-full">
        <div class="absolute left-3 top-1/2 -translate-y-1/2">
          <Search class="text-muted-foreground" />
        </div>
        <Input
          bind:value={query}
          type="text"
          placeholder="Search"
          class="pl-10 pr-12 text-xl border-none h-full w-full rounded-none sm:rounded-l-md"
          disabled={searchStatus === "searching"}
        />
        <div class="absolute right-2 top-1/2 -translate-y-1/2">
          <Tooltip.Provider>
            <Tooltip.Root>
              <Tooltip.Trigger>
                <Toggle
                  aria-label="Exact match"
                  bind:pressed={exactMatch}
                  disabled={searchStatus === "searching"}
                >
                  <SquareEqual class="size-4" />
                </Toggle>
              </Tooltip.Trigger>
              <Tooltip.Content>
                <p>Search for exact match.</p>
              </Tooltip.Content>
            </Tooltip.Root>
          </Tooltip.Provider>
        </div>
      </div>
      <Button
        variant="secondary"
        size="lg"
        type="submit"
        class="sm:rounded-l-none h-full"
        disabled={searchStatus === "searching"}
      >
        {#if searchStatus === "searching"}
          <Loader class="animate-spin size-5" />
        {:else}
          <Search class="size-5" />
        {/if}
      </Button>
    </form>

    {#if searchStatus === "searching" || searchStatus === "complete"}
      <Stats
        {dataleaks}
        {searchStatus}
        {results}
        {researchTime}
        {query}
        {exactMatch}
      />
    {/if}
  </header>

  {#if searchStatus === "searching"}
    <ProgressBar progress={progressPercentage} />
  {/if}

  <ResultPaginated bind:page {results} />

  {#if searchStatus === "complete" && results.length === 0}
    <NoResult />
  {:else if searchStatus === "searching" && results.length === 0}
    <Searching />
  {/if}

  {#if searchStatus === "idle"}
    <div
      class="grid gap-5 grid-cols-1 2xl:grid-cols-2 [&>*:last-child]:col-span-full"
    >
      <SearchbarCard />
      <DatawellsCard {dataleaks} />
      <ContributionCard />
    </div>
  {/if}
</main>
