<script lang="ts">
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import Result from "./result.svelte";

  let { results, page = $bindable() } = $props();
  let perPage = $state(20);
</script>

{#each results.slice((page - 1) * perPage, page * perPage) as result}
  <Result {result} />
{/each}

<div class="my-10">
{#if results.length > perPage}
  <Pagination.Root count={results.length} {perPage} bind:page={page}>
    {#snippet children({ pages, currentPage })}
      <Pagination.Content>
        <Pagination.Item>
          <Pagination.PrevButton />
        </Pagination.Item>
        {#each pages as page (page.key)}
          {#if page.type === "ellipsis"}
            <Pagination.Item>
              <Pagination.Ellipsis />
            </Pagination.Item>
          {:else}
            <Pagination.Item>
              <Pagination.Link {page} isActive={currentPage === page.value}>
                {page.value}
              </Pagination.Link>
            </Pagination.Item>
          {/if}
        {/each}
        <Pagination.Item>
          <Pagination.NextButton />
        </Pagination.Item>
      </Pagination.Content>
    {/snippet}
  </Pagination.Root>
{/if}
</div>
