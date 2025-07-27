<script lang="ts">
  import {
    Timer,
    File,
    ChartSpline,
    ClipboardCopy,
    Database as DatabaseIcon,
  } from "@lucide/svelte";
  import Button from "$src/lib/components/ui/button/button.svelte";
  import { toast } from "svelte-sonner";
  import type { Dataleaks, Result } from "$src/lib/types";

  const {
    dataleaks,
    results,
    researchTime,
    query,
    exactMatch,
    searchStatus,
  }: {
    dataleaks: Dataleaks|null;
    results: Result[];
    researchTime: number|null;
    query: string;
    exactMatch: boolean;
    searchStatus: any;
  } = $props();

  function copyResult() {
    const copyData = {
      Query: {
        Terms: query.split(" "),
        ExactMatch: exactMatch,
      },
      Results: results,
      Time: researchTime,
    };
    navigator.clipboard.writeText(JSON.stringify(copyData, null, 2));
    toast.success("Document copied to clipboard!");
  }
</script>

<div class="flex flex-wrap gap-4 justify-between text-sm mt-5">
  <div class="flex items-center gap-2">
    <Timer class="text-green-400" />
    <span
      >{researchTime === null
        ? "--"
        : Math.round(researchTime / 1000) + "s"}</span
    >
    <span class="text-muted-foreground hidden lg:inline">Elapsed</span>
  </div>
  <div class="flex items-center gap-2">
    <File class="text-blue-400" />
    <span>{results.length}</span>
    <span class="text-muted-foreground hidden lg:inline">Results</span>
  </div>
  <div class="flex items-center gap-2">
    <ChartSpline class="text-red-400" />
    <span>
      {dataleaks?.TotalRows.toLocaleString("fr-FR") || "-- --- --- ---"}
    </span>
    <span class="text-muted-foreground hidden lg:inline">Documents</span>
  </div>
  <div class="flex items-center gap-2">
    <DatabaseIcon class="text-orange-400" />
    <span>{dataleaks?.TotalDataleaks || "---"}</span>
    <span class="text-muted-foreground hidden lg:inline">Data Wells</span>
  </div>
  <Button
    variant="ghost"
    onclick={copyResult}
    disabled={searchStatus === "searching"}
  >
    <ClipboardCopy class="mr-1" />
    <span class="text-muted-foreground hidden lg:inline">Copy result</span>
  </Button>
</div>
