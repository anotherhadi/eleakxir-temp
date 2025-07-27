<script lang="ts">
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import { FormatSize } from "$src/lib/leaks";

  let { dataleaks, totalDocuments, totalSize, totalDataleaks, page = $bindable() } = $props();

  let perPage = $state(20);
</script>

<div class="grid grid-cols-1">
  <Table.Root>
    <Table.Header>
      <Table.Row>
        <Table.Head>Name</Table.Head>
        <Table.Head>Documents</Table.Head>
        <Table.Head>Size (Mb)</Table.Head>
        <Table.Head>Columns</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each dataleaks.slice((page - 1) * perPage, page * perPage) as dataleak}
        <Table.Row>
          <Table.Cell class="font-medium"
            >{dataleak.Name}</Table.Cell
          >
          <Table.Cell>{dataleak.Length.toLocaleString("fr")
          }</Table.Cell>
          <Table.Cell>{FormatSize(dataleak.Size)}</Table.Cell>
          <Table.Cell>{dataleak.Columns.join(", ")}</Table.Cell>
        </Table.Row>
      {/each}
    </Table.Body>
    <Table.Footer>
      <Table.Row>
        <Table.Cell>Total: {totalDataleaks} dataleaks</Table.Cell>
        <Table.Cell>{totalDocuments.toLocaleString("fr")}</Table.Cell>
        <Table.Cell>{FormatSize(totalSize)}</Table.Cell>
        <Table.Cell></Table.Cell>
      </Table.Row>
    </Table.Footer>
  </Table.Root>
</div>

<div class="my-10">
  {#if dataleaks.length > perPage}
    <Pagination.Root count={dataleaks.length} {perPage} bind:page>
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
