<script lang="ts">
  import Button from "$src/lib/components/ui/button/button.svelte";

  const illustrations = [
    "/illustration1.png",
    "/illustration2.png",
    "/illustration3.png",
  ];

  const randomIllustration =
    illustrations[Math.floor(Math.random() * illustrations.length)];
</script>

<div class="grid grid-cols-1 xl:grid-cols-2 gap-5 h-full w-full">
  <div class="w-full h-full flex justify-center items-center ">
    <div class="flex flex-col gap-5 max-w-3xl">
      <div class="flex gap-2 items-end mb-4">
        <img class="size-28" src="/big.png" alt="Eleakxir Logo" />
        <h2 class="mb-2">Eleakxir.</h2>
      </div>
      <h2 class="leading-16 text-5xl mb-6 mt-0 font-medium">
        <span class="text-primary">Search data wells</span> from your computer. Not
        the fastest, but it's <span class="text-primary">local</span> and
        <span class="text-primary">privacy friendly</span> :)
      </h2>
      <p class="mb-4 text-muted-foreground">
        Eleakxir allows you to search through your local csv files. It's an
        open-source app available at <a href="https://github.com/anotherhadi/eleakxir">anotherhadi/eleakxir</a>.
      </p>
      <div class="flex gap-5">
        <Button size="lg" href="/search">Let's search</Button>
        <Button size="lg" variant="outline" href="/dataleaks">Data wells</Button
        >
      </div>
    </div>
  </div>

  <div class="w-full h-full pr-0 py-2 relative hidden xl:block">
    <div class="fixed top-8 right-8">
      <Button size="lg" href="/search">Search</Button>
    </div>
    <div
      class={"rounded-lg w-full h-full bg-cover bg-right"}
      style="background-image: url('{randomIllustration}');"
    ></div>
  </div>
</div>
