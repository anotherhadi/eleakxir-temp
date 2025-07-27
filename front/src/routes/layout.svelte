<script lang="ts">
  import { type Snippet } from "svelte";
  let { children }: { children: Snippet } = $props();
  import { ModeWatcher } from "mode-watcher";
  import { Toaster } from "$lib/components/ui/sonner/index.js";
  import AppSidebar from "$lib/components/navigation/app-sidebar.svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { cn } from "$src/lib/utils";
  import { buttonVariants } from "$src/lib/components/ui/button";
</script>

<ModeWatcher />
<Toaster />

<Sidebar.Provider open={false} style="--sidebar-width: 16rem;">
  <AppSidebar />
  <Sidebar.Inset class="m-0">
    <Sidebar.Trigger
      class={cn(
        "absolute top-3 left-3 z-[100]",
        buttonVariants({ variant: "outline", size: "icon" }),
      )}
    />
    {@render children()}
  </Sidebar.Inset>
</Sidebar.Provider>
