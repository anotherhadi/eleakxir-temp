<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
    import { Github } from "@lucide/svelte";
    import Button from "../ui/button/button.svelte";
</script>

<Card.Root>
  <Card.Header class="gap-5">
    <div
      class="rounded-full w-16 h-16 bg-gray-400/40 grid items-center justify-center"
    >
      <Github class="size-10" />
    </div>
    <Card.Title>Contributing</Card.Title>
  </Card.Header>
  <Card.Content>
    <p>
    This project is open source and welcomes contributions. Whether it's fixing bugs, improving performance, or adding new features, feel free to fork the repository and submit a pull request.
    </p>
    <Button class="mt-4" variant="outline" href="https://github.com/anotherhadi/eleakxir">See the repo</Button>
  </Card.Content>
</Card.Root>
